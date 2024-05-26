package config

import (
	"context"
	"fmt"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db" // Import Realtime Database package
	"google.golang.org/api/option"
)

type FireDB struct {
	*db.Client
}

var fireDB FireDB

func FirebaseDB() *FireDB {
	return &fireDB
}

func (db *FireDB) Connect() error {
	ctx := context.Background()
	serviceAccountKeyFilePath, err := filepath.Abs("../../../serviceAccountKey.json")
	if err != nil {
		return fmt.Errorf("error loading service account key: %v", err)
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	config := &firebase.Config{DatabaseURL: "https://macformula-ui-testing-default-rtdb.firebaseio.com/"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}
	db.Client = client
	return nil
}
