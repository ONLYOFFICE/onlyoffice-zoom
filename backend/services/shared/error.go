package shared

import (
	"fmt"
)

type InvalidConfigurationParameterError struct {
	Parameter string
	Reason    string
}

func (e *InvalidConfigurationParameterError) Error() string {
	return fmt.Sprintf("invald configuration [%s] parameter. Reason: %s", e.Parameter, e.Reason)
}
