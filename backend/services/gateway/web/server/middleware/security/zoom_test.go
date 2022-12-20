package security

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	_ZOOM_CONTEXT = "DO6MJD0Iw0TZVTR2LgAA5QAAAG2Pz7Ok-BIokLdFx3rE6rFN4U9HiZoM5gBuF5rIqiDQGtyy985lf4wOdy8aPrKUscsjJOqOT-_X6ekGMOy3r2CTe0pFYZ9zXoN30yAfQXr1UrbzgjNJYrY-HTW1FzW3TOiI5GbUr-4zhakph1tGV8Bdpo0jfYtzIjKsB6ko2I8vHUwy0Guz04golMTOnQnSLcF-jhrZPrr1rNHMSsJbhXH-5cMw3ui_RNYRHX7O6SU1_lApDWB2fgKBv53xg_jOHNwrX-JB3ad2u1zXTcTvc_yYO2fC-trHhadxMRD-VuLzxJxSnodQEl0TsqF4nfHDjLunHkO6"
	_ZOOM_SECRET  = "bllTHHZB1GHnBt5tevD76YyYphHBI4kK"

	_EXPECTED_DECODED_CONTEXT = "0cee8c243d08c344d95534762e0000e50000006d8fcfb3a4f8122890b745c77ac4eab14de14f47899a0ce6006e179ac8aa20d01adcb2f7ce657f8c0e772f1a3eb294b1cb2324ea8e4fefd7e9e90630ecb7af60937b4a45619f735e8377d3201f417af552b6f382334962b63e1d35b51735b74ce888e466d4afee3385a929875b4657c05da68d237d8b732232ac07a928d88f2f1d4c32d06bb3d3882894c4ce9d09d22dc17e8e1ad93ebaf5acd1cc4ac25b8571fee5c330dee8bf44d6111d7ecee92535fe50290d60767e0281bf9df183f8ce1cdc2b5fe241dda776bb5cd74dc4ef73fc983b67c2fadac785a7713110fe56e2f3c49c529e8750125d13b2a1789df1c38cbba71e43ba"
	_EXPECTED_IV              = "ee8c243d08c344d95534762e"
	_EXPECTED_AAD             = ""
	_EXPECTED_CIPHER_TEXT     = "6d8fcfb3a4f8122890b745c77ac4eab14de14f47899a0ce6006e179ac8aa20d01adcb2f7ce657f8c0e772f1a3eb294b1cb2324ea8e4fefd7e9e90630ecb7af60937b4a45619f735e8377d3201f417af552b6f382334962b63e1d35b51735b74ce888e466d4afee3385a929875b4657c05da68d237d8b732232ac07a928d88f2f1d4c32d06bb3d3882894c4ce9d09d22dc17e8e1ad93ebaf5acd1cc4ac25b8571fee5c330dee8bf44d6111d7ecee92535fe50290d60767e0281bf9df183f8ce1cdc2b5fe241dda776bb5cd74dc4ef73fc983b67c2fadac785a7713110fe56e2f3c49c529e87"
	_EXPECTED_TAG             = "50125d13b2a1789df1c38cbba71e43ba"
)

func TestDecodeContext(t *testing.T) {
	t.Parallel()

	var buffer []byte
	var key []byte
	var ok bool

	ctx, err := buildDecodedBuffer(_ZOOM_CONTEXT, _ZOOM_SECRET)

	assert.Nil(t, err)

	if buffer, ok = ctx.Value(contextBuffer{}).([]byte); !ok {
		assert.Fail(t, "could not extract buffer from the context")
	}

	if key, ok = ctx.Value(secretKey{}).([]byte); !ok {
		assert.Fail(t, "could not extract key from the context")
	}

	assert.Equal(t, 32, len(key))
	assert.Equal(t, _EXPECTED_DECODED_CONTEXT, fmt.Sprintf("%x", buffer))
}

func TestExtractIV(t *testing.T) {
	t.Parallel()

	var iv []byte
	var ok bool

	ctx, _ := buildDecodedBuffer(_ZOOM_CONTEXT, _ZOOM_SECRET)
	ctx, err := extractIV(ctx)

	assert.Nil(t, err)

	if iv, ok = ctx.Value(contextIV{}).([]byte); !ok {
		assert.Fail(t, "could not extract iv from the context")
	}

	assert.Equal(t, _EXPECTED_IV, fmt.Sprintf("%x", iv))
}

func TestExtractAll(t *testing.T) {
	ctx, _ := buildDecodedBuffer(_ZOOM_CONTEXT, _ZOOM_SECRET)
	ctx, err := pipe(extractIV, extractAAD, extractCipherText)(ctx)
	assert.Nil(t, err)

	if iv, ok := ctx.Value(contextIV{}).([]byte); !ok {
		assert.Fail(t, "could not extract iv from the context")
	} else {
		assert.Equal(t, _EXPECTED_IV, fmt.Sprintf("%x", iv))
	}

	if aad, ok := ctx.Value(contextAad{}).([]byte); !ok {
		assert.Fail(t, "could not extract aad from the context")
	} else {
		assert.Equal(t, _EXPECTED_AAD, fmt.Sprintf("%x", aad))
	}

	if cipherText, ok := ctx.Value(contextText{}).([]byte); !ok {
		assert.Fail(t, "could not extract cipherText from the context")
	} else {
		assert.Equal(t, _EXPECTED_CIPHER_TEXT, fmt.Sprintf("%x", cipherText))
	}

	if tag, ok := ctx.Value(contextTag{}).([]byte); !ok {
		assert.Fail(t, "could not extract tag from the context")
	} else {
		assert.Equal(t, _EXPECTED_TAG, fmt.Sprintf("%x", tag))
	}
}

func TestDecryptZoomContext(t *testing.T) {
	ctx, _ := buildDecodedBuffer(_ZOOM_CONTEXT, _ZOOM_SECRET)
	ctx, err := pipe(extractIV, extractAAD, extractCipherText)(ctx)
	assert.Nil(t, err)
	zcontext, err := decrypt(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, zcontext)
	assert.Less(t, zcontext.Exp, int(time.Now().UnixMilli()))
}

func TestExtractZoomContext(t *testing.T) {
	zcontext, err := ExtractZoomContext(_ZOOM_CONTEXT, _ZOOM_SECRET)
	assert.Nil(t, err)
	assert.NotNil(t, zcontext)
}
