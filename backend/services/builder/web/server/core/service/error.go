package service

import "fmt"

type InvalidServiceParameterError struct {
	Name   string
	Reason string
}

func (e *InvalidServiceParameterError) Error() string {
	return fmt.Sprintf("invald service parameter %s. Reason: %s", e.Name, e.Reason)
}
