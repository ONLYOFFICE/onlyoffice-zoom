package handler

import "errors"

var ErrInvalidContextValue = errors.New("could not extract context value")
var ErrEmptyIdValue = errors.New("could not perform current action with an empty id")
var ErrUnauthorizedAccess = errors.New("unauthorized file access")
