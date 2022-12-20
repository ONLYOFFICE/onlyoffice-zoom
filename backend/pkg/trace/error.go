package trace

import (
	"errors"
)

var ErrTracerInvalidNameInitialization = errors.New("could not initialize a tracer with an empty name")
var ErrTracerInvalidAddressInitialization = errors.New("could not initialize a tracer without a URL")
