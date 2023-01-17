package message

import "errors"

var ErrEmptyIdValue = errors.New("could not perform current action with an empty id")
var _ErrInvalidHandlerPayload = errors.New("invalid handler payload")
