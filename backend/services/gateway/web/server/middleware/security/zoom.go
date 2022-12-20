package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var errInvalidIVBuffer = errors.New("could not extract IV from an invalid buffer")
var errInvalidAADBuffer = errors.New("could not extract AAD from an invalid buffer")
var errInvalidCipherTextBuffer = errors.New("could not extract cipher text from an invalid buffer")

type ZoomContextExtractionError struct {
	Part string
}

func (e *ZoomContextExtractionError) Error() string {
	return fmt.Sprintf("could not extract %s from the context", e.Part)
}

type contextBuffer struct{}
type contextIV struct{}
type contextAad struct{}
type contextText struct{}
type contextTag struct{}
type secretKey struct{}

type ZoomContext struct {
	Type string `json:"typ"`
	Uid  string `json:"uid"`
	Mid  string `json:"mid"`
	Aud  string `json:"aud"`
	Iss  string `json:"iss"`
	Ts   int    `json:"ts"`
	Exp  int    `json:"exp"`
}

func GenerateState(secret string) (string, error) {
	ts, err := randomHex(64)
	if err != nil {
		return "", err
	}

	hmac, err := hmacBase64(ts, secret)
	if err != nil {
		return "", err
	}

	return url.QueryEscape(strings.ReplaceAll(strings.Join([]string{hmac, ts}, "."), "+", "")), nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hmacBase64(message string, secret string) (string, error) {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)

	if _, err := h.Write([]byte(message)); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func ExtractZoomContext(context string, secret string) (ZoomContext, error) {
	ctx, err := buildDecodedBuffer(context, secret)
	if err != nil {
		return ZoomContext{}, err
	}

	ctx, err = pipe(extractIV, extractAAD, extractCipherText)(ctx)
	if err != nil {
		return ZoomContext{}, err
	}

	return decrypt(ctx)
}

func pipe(fs ...func(context.Context) (context.Context, error)) func(context.Context) (context.Context, error) {
	return func(c context.Context) (context.Context, error) {
		var err error
		for _, f := range fs {
			if c, err = f(c); err != nil {
				return nil, err
			}
		}
		return c, nil
	}
}

func buildDecodedBuffer(zctx string, secret string) (context.Context, error) {
	udecoded, err := base64.RawURLEncoding.DecodeString(zctx)
	if err != nil {
		return nil, err
	}

	key := sha256.Sum256([]byte(secret))

	ctx := context.WithValue(context.Background(), secretKey{}, key[:])
	return context.WithValue(ctx, contextBuffer{}, udecoded), nil
}

func extractIV(ctx context.Context) (context.Context, error) {
	var buffer []byte
	var ok bool
	if buffer, ok = ctx.Value(contextBuffer{}).([]byte); !ok {
		return nil, errInvalidIVBuffer
	}

	if len(buffer) < 2 {
		return nil, errInvalidIVBuffer
	}

	ivLen := buffer[0]
	buffer = buffer[1:]

	if len(buffer) <= int(ivLen) {
		return nil, errInvalidIVBuffer
	}

	iv := buffer[:ivLen]
	buf := buffer[ivLen:]

	ctx = context.WithValue(ctx, contextBuffer{}, buf)
	return context.WithValue(ctx, contextIV{}, iv), nil
}

func extractAAD(ctx context.Context) (context.Context, error) {
	var buffer []byte
	var ok bool
	if buffer, ok = ctx.Value(contextBuffer{}).([]byte); !ok {
		return nil, errInvalidIVBuffer
	}

	if len(buffer) < 2 {
		return nil, errInvalidAADBuffer
	}

	aadLen := binary.LittleEndian.Uint16(buffer[:2])
	buffer = buffer[2:]

	aad := buffer[:aadLen]
	buf := buffer[aadLen:]

	ctx = context.WithValue(ctx, contextBuffer{}, buf)
	return context.WithValue(ctx, contextAad{}, aad), nil
}

func extractCipherText(ctx context.Context) (context.Context, error) {
	var buffer []byte
	var ok bool
	if buffer, ok = ctx.Value(contextBuffer{}).([]byte); !ok {
		return nil, errInvalidIVBuffer
	}

	if len(buffer) < 4 {
		return nil, errInvalidCipherTextBuffer
	}

	cipherLength := binary.LittleEndian.Uint32(buffer[:4])
	if int(cipherLength) >= len(buffer) {
		return nil, errInvalidCipherTextBuffer
	}

	buffer = buffer[4:]
	cipherText := buffer[:cipherLength]
	tag := buffer[cipherLength:]

	ctx = context.WithValue(ctx, contextText{}, cipherText)
	return context.WithValue(ctx, contextTag{}, tag), nil
}

func decrypt(ctx context.Context) (ZoomContext, error) {
	var iv, cipherText, tag, key []byte
	var ok bool

	if iv, ok = ctx.Value(contextIV{}).([]byte); !ok {
		return ZoomContext{}, &ZoomContextExtractionError{
			Part: "iv",
		}
	}

	if cipherText, ok = ctx.Value(contextText{}).([]byte); !ok {
		return ZoomContext{}, &ZoomContextExtractionError{
			Part: "cipherText",
		}
	}

	if tag, ok = ctx.Value(contextTag{}).([]byte); !ok {
		return ZoomContext{}, &ZoomContextExtractionError{
			Part: "tag",
		}
	}

	if key, ok = ctx.Value(secretKey{}).([]byte); !ok {
		return ZoomContext{}, &ZoomContextExtractionError{
			Part: "secret key",
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return ZoomContext{}, err
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return ZoomContext{}, err
	}

	ciphertext := append(cipherText, tag...)

	plaintext, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return ZoomContext{}, err
	}

	var zoom ZoomContext
	if err := json.Unmarshal(plaintext, &zoom); err != nil {
		return ZoomContext{}, err
	}

	return zoom, nil
}
