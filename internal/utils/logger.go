package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger initializes the global logger
func InitLogger(debug bool) {
	Logger = logrus.New()

	// Set output to stderr (standard for CLI tools)
	Logger.SetOutput(os.Stderr)

	// Set log level based on debug flag
	if debug {
		Logger.SetLevel(logrus.DebugLevel)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}

	// Use a simple formatter for CLI output
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
	})
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if Logger == nil {
		InitLogger(false)
	}
	return Logger
}
