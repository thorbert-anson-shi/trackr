package handlers

import (
	"context"
	"tobtoby/trackr/database"
	"tobtoby/trackr/firebase"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"firebase.google.com/go/messaging"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

func ListLocationsHandler(c fiber.Ctx) error {
	queries := generated.New(database.DB)

	locations, err := queries.ListLocations(c.Context())
	if err != nil {
		logging.GlobalLogger.Println("An error occurred when fetching from DB")
		return err
	}

	// When a user requests the other devices' locations,
	// the server publishes a "requestLocation" event to the queue.
	// It expects the subscribers to then send a request to the POST /locations endpoint
	return c.JSON(locations)
}

func requestLocationUpdates(c context.Context, queries *generated.Queries) error {
	registrationTokens, err := queries.ListRegistrationTokens(c)
	if err != nil {
		logging.GlobalLogger.Println("An error occurred when fetching registration tokens")
		return err
	}

	nativeRegistrationTokens := registrationTokenPGToNative(registrationTokens)

	batchResponse, err := firebase.MsgClient.SendMulticast(c, &messaging.MulticastMessage{
		Tokens: nativeRegistrationTokens,
		Data:   map[string]string{"event": "sendLocation"},
	})
	if err != nil {
		logging.GlobalLogger.Println("Failed to send all messages")
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

// TODO: Figure out what pgtype.Text.Valid means
func registrationTokenPGToNative(arr []pgtype.Text) []string {
	var res []string
	for _, token := range arr {
		res = append(res, token.String)
	}

	return res
}
