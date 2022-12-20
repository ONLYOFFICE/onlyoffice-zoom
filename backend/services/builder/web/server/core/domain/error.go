package domain

import "fmt"

type InvalidModelFieldError struct {
	Model  string
	Field  string
	Reason string
}

func (e *InvalidModelFieldError) Error() string {
	return fmt.Sprintf("invald %s field %s. Reason: %s", e.Model, e.Field, e.Reason)
}
