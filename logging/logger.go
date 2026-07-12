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
		log.Default().Fatalf("Failed to open file: %s\n", err.Error())
	}

	GlobalLogger = log.New(io.MultiWriter(appLogWriter, os.Stdout), "[trackr] ", log.LstdFlags)

	GlobalLogger.Println("Successfully initialized application logger")
}

func initializePollLogger() {
	pollingLogFile := config.SafeFetchVar("POLL_LOGS")

	pollLogWriter, err := os.OpenFile(pollingLogFile, os.O_RDWR, 0o666)
	if err != nil {
		log.Default().Fatalf("Failed to open file: %s\n", err.Error())
	}

	PollingLogger = log.New(pollLogWriter, "[poll] ", log.LstdFlags)

	GlobalLogger.Println("Successfully initialized poll logger")
	PollingLogger.Println("Successfully initialized poll logger")
}
