package hookfs

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// LogLevelMin is the minimum log level
const LogLevelMin = 0

// LogLevelMax is the maximum log level
const LogLevelMax = 2

var logLevel int

func initLog() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
}

// LogLevel gets the log level.
func LogLevel() int {
	return logLevel
}

// SetLogLevel sets the log level. newLevel must be >= LogLevelMin, and <= LogLevelMax.
func SetLogLevel(newLevel int) {
	if newLevel < LogLevelMin || newLevel > LogLevelMax {
		log.Fatalf("Bad log level: %d (must be %d..%d)", newLevel, LogLevelMin, LogLevelMax)
	}

	logLevel = newLevel
	if logLevel > 0 {
		log.SetLevel(log.DebugLevel)
	}
}

func init() {
	initLog()
	SetLogLevel(0)
}
