package client

import (
	"errors"
	"fmt"
)

var ErrInvalidAccessToken error = errors.New("could not perform zoom action due to invalid access token")
var ErrInvalidAuthorizationCode error = errors.New("could not perform zoom action due to invalid authorization code")
var ErrInvalidUrlFormat error = errors.New("url is not valid")
var ErrEmptyDeeplinkResponse error = errors.New("could not get deeplink")
var ErrFileDoesNotExist error = errors.New("could not file with this id")

type UnexpectedStatusCodeError struct {
	Action string
	Code   int
}

func (e *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("could not perform zoom %s action. Status code: %d", e.Action, e.Code)
}

type InvalidQueryParameterError struct {
	Parameter string
}

func (e *InvalidQueryParameterError) Error() string {
	return fmt.Sprintf("could not perform send a zoom api request. Invalid query parameter: %s", e.Parameter)
}
