// Package polling contains logic for polling devices at certain intervals
package polling

import (
	"context"
	"strconv"
	"time"

	"tobtoby/trackr/config"
	"tobtoby/trackr/database"
	"tobtoby/trackr/firebase"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"firebase.google.com/go/messaging"
	"github.com/jackc/pgx/v5/pgtype"
)

func InitializePoller(appCtx context.Context) {
	isDevEnvironment := config.SafeFetchVar("DEVELOPMENT") == "true"
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
				}
				logging.PollingLogger.Println("Polling FCM")
				jobCtx, cancel := context.WithTimeout(appCtx, 5*time.Second)
				if err = requestLocationUpdates(jobCtx, generated.New(database.DB)); err != nil {
					logging.PollingLogger.Printf("An error occurred when requesting locations: %s\n", err.Error())
				}
				cancel()
			case <-appCtx.Done():
				logging.PollingLogger.Println("Poller stopped by application")
				return
			}
		}
	}()
}

func requestLocationUpdates(c context.Context, queries *generated.Queries) error {
	registrationTokens, err := queries.ListRegistrationTokens(c)
	if err != nil {
		logging.PollingLogger.Println("An error occurred when fetching registration tokens")
		return err
	}

	nativeRegistrationTokens := registrationTokenPGToNative(registrationTokens)

	// If no valid registration tokens found, don't poll
	// NOTE: This check is done after parsing because admin returns a NULL registration token,
	// causing the len(registrationTokens) == 0 check to fail
	if len(nativeRegistrationTokens) == 0 {
		return nil
	}

	logging.PollingLogger.Printf("Registration tokens: %s\n", nativeRegistrationTokens)

	batchResponse, err := firebase.MsgClient.SendMulticast(c, &messaging.MulticastMessage{
		Tokens: nativeRegistrationTokens,
		Data:   map[string]string{"event": "sendLocation"},
	})
	if err != nil {
		logging.PollingLogger.Printf("Failed to send all messages: %s\n", err.Error())
		return err
	}

	if batchResponse.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range batchResponse.Responses {
			if !resp.Success {
				failedTokens = append(failedTokens, nativeRegistrationTokens[idx])
			}
		}

		logging.PollingLogger.Printf("Tokens causing failure: %v\n", failedTokens)
	}

	return nil
}

func registrationTokenPGToNative(arr []pgtype.Text) []string {
	var res []string
	for _, token := range arr {
		// check if token string in DB is not NULL
		if token.Valid {
			res = append(res, token.String)
		}
	}

	return res
}
