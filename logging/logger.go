// Package logging defines initialization logic for loggers
package logging

import (
	"io"
	"log"
	"os"

	"tobtoby/trackr/config"
)

var (
	GlobalLogger  *log.Logger
	PollingLogger *log.Logger
)

func InitializeLogger() {
	initializeApplicationLogger()
	initializePollLogger()
}

func initializeApplicationLogger() {
	applicationLogFile := config.SafeFetchVar("LOG_FILE")

	appLogWriter, err := os.OpenFile(applicationLogFile, os.O_RDWR, 0o666)
	if err != nil {
		GlobalLogger.Fatalf("Failed to open file: %s\n", err.Error())
	}

	// Clear log file if previously had contents
	if err = appLogWriter.Truncate(0); err != nil {
		GlobalLogger.Fatalf("Failed to truncate existing file: %s\n", err.Error())
	}

	GlobalLogger = log.New(io.MultiWriter(appLogWriter, os.Stdout), "", log.LstdFlags)

	GlobalLogger.Println("Successfully initialized application logger")
}

func initializePollLogger() {
	pollingLogFile := config.SafeFetchVar("POLL_LOGS")

	pollLogWriter, err := os.OpenFile(pollingLogFile, os.O_RDWR, 0o666)
	if err != nil {
		PollingLogger.Fatalf("Failed to open file: %s\n", err.Error())
	}

	if err = pollLogWriter.Truncate(0); err != nil {
		PollingLogger.Fatalf("Failed to truncate existing file: %s\n", err.Error())
	}

	PollingLogger = log.New(pollLogWriter, "", log.LstdFlags)

	GlobalLogger.Println("Successfully initialized poll logger")
	PollingLogger.Println("Successfully initialized poll logger")
}
