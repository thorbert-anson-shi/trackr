package logging

import (
	"io"
	"log"
	"os"
	"tobtoby/trackr/config"
)

var GlobalLogger *log.Logger
var PollingLogger *log.Logger

func InitializeLogger() {
	initializeApplicationLogger()
	initializePollLogger()
}

func initializeApplicationLogger() {
	applicationLogFile := config.SafeFetchVar("LOG_FILE")

	appLogWriter, err := os.OpenFile(applicationLogFile, os.O_RDWR, 0666)
	if err != nil {
		GlobalLogger.Fatalf("Failed to open file: %s\n", err.Error())
	}

	// Clear log file if previously had contents
	appLogWriter.Truncate(0)

	GlobalLogger = log.New(io.MultiWriter(appLogWriter, os.Stdout), "", log.LstdFlags)

	GlobalLogger.Println("Successfully initialized application logger")

}

func initializePollLogger() {
	pollingLogFile := config.SafeFetchVar("POLL_LOGS")

	pollLogWriter, err := os.OpenFile(pollingLogFile, os.O_RDWR, 0666)
	if err != nil {
		PollingLogger.Fatalf("Failed to open file: %s\n", err.Error())
	}

	pollLogWriter.Truncate(0)

	PollingLogger = log.New(pollLogWriter, "", log.LstdFlags)

	PollingLogger.Println("Successfully initialized poll logger")

}
