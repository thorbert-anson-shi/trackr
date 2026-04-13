package polling

import (
	"context"
	"time"
	"tobtoby/trackr/database"
	"tobtoby/trackr/firebase"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"firebase.google.com/go/messaging"
	"github.com/jackc/pgx/v5/pgtype"
)

func InitializePoller(appCtx context.Context) {
	queries := generated.New(database.DB)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Catch if application is cancelled mid-poll
				select {
				case <-appCtx.Done():
					logging.GlobalLogger.Println("Poller stopped by application")
					return
				default:
				}
				logging.GlobalLogger.Println("Polling FCM")
				jobCtx, cancel := context.WithTimeout(appCtx, 5*time.Second)
				requestLocationUpdates(jobCtx, queries)
				cancel()
			case <-appCtx.Done():
				logging.GlobalLogger.Println("Poller stopped by application")
				return
			}
		}
	}()
}

func requestLocationUpdates(c context.Context, queries *generated.Queries) error {
	registrationTokens, err := queries.ListRegistrationTokens(c)
	if err != nil {
		logging.GlobalLogger.Println("An error occurred when fetching registration tokens")
		return err
	}

	nativeRegistrationTokens := registrationTokenPGToNative(registrationTokens)

	logging.GlobalLogger.Printf("Registration tokens: %s\n", nativeRegistrationTokens)

	batchResponse, err := firebase.MsgClient.SendMulticast(c, &messaging.MulticastMessage{
		Tokens: nativeRegistrationTokens,
		Data:   map[string]string{"event": "sendLocation"},
	})
	if err != nil {
		logging.GlobalLogger.Printf("Failed to send all messages: %s\n", err.Error())
		return err
	}

	if batchResponse.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range batchResponse.Responses {
			if !resp.Success {
				failedTokens = append(failedTokens, nativeRegistrationTokens[idx])
			}
		}

		logging.GlobalLogger.Printf("Tokens causing failure: %v\n", failedTokens)
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
