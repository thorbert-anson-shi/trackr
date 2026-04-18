package firebase

import (
	"context"
	"fmt"
	"tobtoby/trackr/config"
	"tobtoby/trackr/logging"

	fb "firebase.google.com/go"

	"google.golang.org/api/option"
)

func ConnectFirebase() {
	logging.GlobalLogger.Println("Reading service account configuration")

	var err error
	opt := option.WithAuthCredentialsFile(option.ServiceAccount, config.SafeFetchVar("SERVICE_ACCOUNT_PATH"))

	app, err := fb.NewApp(context.Background(), nil, opt)
	if err != nil {
		logging.GlobalLogger.Fatalln(fmt.Errorf("Error initializing app: %v", err))
	}

	MsgClient, err = app.Messaging(context.Background())
	if err != nil {
		logging.GlobalLogger.Fatalln("Failed to connect to Firebase")
	}

	logging.GlobalLogger.Println("Connection to Firebase established")
}
