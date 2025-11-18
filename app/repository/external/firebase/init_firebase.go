package firebase

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var firebase_app *firebase.App
var auth_client *auth.Client

// initialize_firebase は Firebase Admin SDK を初期化します
func initialize_firebase() error {
	var ctx context.Context
	var credentials_path string
	var opt option.ClientOption
	var app *firebase.App
	var client *auth.Client
	var err error

	ctx = context.Background()
	credentials_path = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentials_path == "" {
		return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
	}
	opt = option.WithCredentialsFile(credentials_path)
	app, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize firebase app: %w", err)
	}
	firebase_app = app
	client, err = app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth client: %w", err)
	}
	auth_client = client
	return nil
}

// get_auth_client は Firebase Auth クライアントを取得します
func get_auth_client() (*auth.Client, error) {
	if firebase_app == nil || auth_client == nil {
		return nil, fmt.Errorf("firebase is not initialized")
	}
	return auth_client, nil
}
