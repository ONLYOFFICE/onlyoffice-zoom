package model

import "errors"

var ErrInvalidTokenFormat error = errors.New("could not perform zoom action due to unexpected token format")
