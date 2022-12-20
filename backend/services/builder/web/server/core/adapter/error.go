package adapter

import "errors"

var ErrSessionAlreadyExists = errors.New("session with this meetingID already exists")
