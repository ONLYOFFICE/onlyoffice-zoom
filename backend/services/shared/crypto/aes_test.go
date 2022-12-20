package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_AES_SECRET    = "sec"
	_AES_TEXT      = "mock"
	_AES_ENCRYPTED = "79xG0i1SV2Kxz7KB8FKwMQqozRm1ADf45zLDq2MHBd8="
)

func TestEncryptText(t *testing.T) {
	type test struct {
		text  string
		isErr bool
	}

	t.Parallel()
	encryptor := NewAesEncryptor([]byte(_AES_SECRET))

	tests := []test{
		{text: _AES_TEXT, isErr: false},
		{text: "1235423523623", isErr: false},
		{text: "477bbd54-3475-4036-bb64-cafd07275632", isErr: false},
		{text: "b29b5cca7ea66fa4aaeda02238b652b5dad0f31ab52a6f81a785ca4b73c577e97dac14dbf0bc24f8a0371e891de6bd304bddda26bef10f921d7079df7e0a7ccca52b9ab4a47e170f3a2a2d3c3dffeae9", isErr: false},
	}

	for _, test := range tests {
		cipher, err := encryptor.Encrypt(test.text)
		if test.isErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotEmpty(t, cipher)
		}
	}
}

func TestDecryptText(t *testing.T) {
	t.Parallel()
	encryptor := NewAesEncryptor([]byte(_AES_SECRET))

	text, err := encryptor.Decrypt(_AES_ENCRYPTED)
	assert.NoError(t, err)
	assert.Equal(t, _AES_TEXT, text)
}
