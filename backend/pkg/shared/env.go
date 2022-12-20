package shared

const (
	DEV  = "development"
	TEST = "testing"
	PROD = "production"
)

var SUPPORTED_ENVIRONMENTS = map[string]string{
	"development": DEV,
	"testing":     TEST,
	"production":  PROD,
	"dev":         DEV,
	"test":        TEST,
	"prod":        PROD,
	"d":           DEV,
	"t":           TEST,
	"p":           PROD,
}
