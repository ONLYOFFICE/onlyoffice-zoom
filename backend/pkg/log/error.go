package log

import "fmt"

// LogFileInitializationError fires if Lumberjack throws an error.
type LogFileInitializationError struct {
	Path  string
	Cause error
}

func (e *LogFileInitializationError) Error() string {
	return fmt.Sprintf("could not open/create a log file with path: %s. Cause: %s", e.Path, e.Cause.Error())
}

// LogElasticInitializationError fires when an elastic client throws an error.
type LogElasticInitializationError struct {
	Address string
	Cause   error
}

func (e *LogElasticInitializationError) Error() string {
	return fmt.Sprintf("could not initialize an elastic client/hook with address: %s. Cause: %s", e.Address, e.Cause.Error())
}
