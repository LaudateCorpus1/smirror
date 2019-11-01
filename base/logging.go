package base

import (
	"github.com/viant/toolbox"
	"os"
)

//LoggingEnvKey logging key
const LoggingEnvKey = "LOGGING"

//IsLoggingEnabled returns true if logging enabled
func IsLoggingEnabled() bool {
	return IsFnLoggingEnabled(LoggingEnvKey)
}

//IsFnLoggingEnabled returns true if logging is enabled
func IsFnLoggingEnabled(key string) bool {
	return toolbox.AsBoolean(os.Getenv(key))
}
