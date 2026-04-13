package logging

import (
	"io"
	"log"
	"os"
	"tobtoby/trackr/config"
)

func InitializeLogger() {
	GlobalLogger = log.Default()

	logOutput, err := os.OpenFile(config.Config["LOG_FILE"], os.O_RDWR, 0666)
	if err != nil {
		GlobalLogger.Fatalf("Failed to open file: %s\n", err.Error())
	}

	outputChannels := io.MultiWriter(logOutput, os.Stdout)

	GlobalLogger.SetOutput(outputChannels)

	GlobalLogger.Println("Successfully initialized logger")
}
