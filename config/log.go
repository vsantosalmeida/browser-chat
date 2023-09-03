package config

import (
	"os"
	"strings"

	"github.com/apex/log"
	logJson "github.com/apex/log/handlers/json"
)

// InitLogging initializes logging with json format
func InitLogging() {
	log.SetHandler(logJson.New(os.Stderr))

	setLogLevel()
}

func setLogLevel() {
	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	}
}
