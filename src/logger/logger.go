package logger

import "fmt"

// Logger is the logger.
type Logger int

// Log is the basic logger.
var Log Logger = 0

// Logging levels.
const (
	LEVEL_INFO    = 0
	LEVEL_DEBUG   = 10
	LEVEL_WARNING = 20
	LEVEL_FATAL   = 30
	LEVEL_PANIC   = 999
)

// Logln is the base function to use to log.
func (lg *Logger) Logln(level int, args ...interface{}) {
	if int(*lg) < level {
		fmt.Println("[", int(*lg), "]", args)
	}
}
