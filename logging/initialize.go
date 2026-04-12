package logging

import "log"

func InitializeLogger() {
	GlobalLogger = log.Default()

	GlobalLogger.Println("Successfully initialized logger")
}
