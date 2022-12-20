package crypto

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
)

var ErrJwtManagerSigning = errors.New("could not generate a new jwt")
var ErrJwtManagerEmptyToken = errors.New("could not verify an empty jwt")
var ErrJwtManagerEmptyDecodingBody = errors.New("could not decode a jwt. Got empty interface")
var ErrJwtManagerInvalidSigningMethod = errors.New("unexpected jwt signing method")
var ErrJwtManagerCastOrInvalidToken = errors.New("could not cast claims or invalid jwt")

type JwtManager interface {
	Sign(payload interface {
		Valid() error
	}) (string, error)
	Verify(jwtToken string, body interface{}) error
}

type onlyofficeJwtManager struct {
	key []byte
}

func NewOnlyofficeJwtManager(key string) (JwtManager, error) {
	key = strings.TrimSpace(key)

	if key == "" {
		return onlyofficeJwtManager{}, errors.New("onlyoffice jwt manager's constructor expected a valid jwt secret")
	}

	return onlyofficeJwtManager{
		key: []byte(key),
	}, nil
}

func (j onlyofficeJwtManager) Sign(payload interface {
	Valid() error
}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	ss, err := token.SignedString(j.key)

	if err != nil {
		return "", ErrJwtManagerSigning
	}

	return ss, nil
}

func (j onlyofficeJwtManager) Verify(jwtToken string, body interface{}) error {
	if jwtToken == "" {
		return ErrJwtManagerEmptyToken
	}

	if body == nil {
		return ErrJwtManagerEmptyDecodingBody
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrJwtManagerInvalidSigningMethod
		}

		return j.key, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return ErrJwtManagerCastOrInvalidToken
	} else {
		return mapstructure.Decode(claims, body)
	}
}
