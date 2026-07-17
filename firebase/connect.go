package firebase

import (
	"context"
	"fmt"

	"tobtoby/trackr/config"
	"tobtoby/trackr/logging"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func ConnectFirebase() {
	logging.GlobalLogger.Println("Reading service account configuration")

	var err error
	opt := option.WithAuthCredentialsFile(option.ServiceAccount, config.SafeFetchVar("SERVICE_ACCOUNT_PATH"))

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logging.GlobalLogger.Fatalln(fmt.Errorf("error initializing app: %v", err))
	}

	MsgClient, err = app.Messaging(context.Background())
	if err != nil {
		logging.GlobalLogger.Fatalf("Failed to connect to Firebase: %s\n", err.Error())
	}

	logging.GlobalLogger.Println("Connection to Firebase established")
}
