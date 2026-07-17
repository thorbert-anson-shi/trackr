// Package polling contains logic for polling devices at certain intervals
package polling

import (
	"context"
	"strconv"
	"strings"
	"time"

	"tobtoby/trackr/config"
	"tobtoby/trackr/database"
	"tobtoby/trackr/firebase"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"firebase.google.com/go/v4/messaging"
)

func InitializePoller(appCtx context.Context) {
	isDevEnvironment := strings.EqualFold(config.SafeFetchVar("DEVELOPMENT"), "true")

	devTickInterval, err := strconv.Atoi(config.SafeFetchVar("TICK_INTERVAL_DEV"))
	if err != nil {
		logging.GlobalLogger.Fatalln("TICK_INTERVAL_DEV needs to be an int")
	}
	prodTickInterval, err := strconv.Atoi(config.SafeFetchVar("TICK_INTERVAL_PROD"))
	if err != nil {
		logging.GlobalLogger.Fatalln("TICK_INTERVAL_PROD needs to be an int")
	}

	var ticker *time.Ticker
	if isDevEnvironment {
		ticker = time.NewTicker(time.Duration(devTickInterval) * time.Second)
	} else {
		ticker = time.NewTicker(time.Duration(prodTickInterval) * time.Second)
	}

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Catch if application is cancelled mid-poll
				select {
				case <-appCtx.Done():
					logging.PollingLogger.Println("Poller stopped by application")
					return
				default:
					jobCtx, cancel := context.WithTimeout(appCtx, 5*time.Second)
					if err = requestLocationUpdates(jobCtx, generated.New(database.DB)); err != nil {
						logging.PollingLogger.Printf("An error occurred when requesting locations: %s\n", err.Error())
					}
					cancel()
				}
			case <-appCtx.Done():
				logging.PollingLogger.Println("Poller stopped by application")
				return
			}
		}
	}()
}

func requestLocationUpdates(c context.Context, queries *generated.Queries) error {
	firebaseIds, err := queries.ListFirebaseIDs(c)
	if err != nil {
		logging.PollingLogger.Println("An error occurred when fetching Firebase IDs")
		return err
	}

	// If no valid Firebase IDs found, don't poll
	if len(firebaseIds) == 0 {
		return nil
	}

	logging.PollingLogger.Printf("Installation IDs: %s\n", firebaseIds)

	batchResponse, err := firebase.MsgClient.SendEachForMulticast(c, &messaging.MulticastMessage{
		Fids: firebaseIds,
		Data: map[string]string{"event": "sendLocation"},
	})
	if err != nil {
		logging.PollingLogger.Printf("Failed to send all messages: %s\n", err.Error())
		return err
	}

	if batchResponse.FailureCount > 0 {
		var failedIDs []string
		for idx, resp := range batchResponse.Responses {
			if !resp.Success {
				failedIDs = append(failedIDs, firebaseIds[idx])
			}
		}

		logging.PollingLogger.Printf("Installation IDs causing failure: %v\n", failedIDs)
	}

	return nil
}
