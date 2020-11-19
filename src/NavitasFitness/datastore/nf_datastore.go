package datastore

import (
	"context"
	log "logger"

	"cloud.google.com/go/datastore"
)

const PROJECT_ID = "navitas-fitness-aarhus"

var dsClient *datastore.Client = nil

func createClient() error {
	var err error
	ctx := context.Background()
	log.Infof(ctx, "Creating DataStore Client")
	// Create a datastore client. In a typical application, you would create
	// a single client which is reused for every datastore operation.
	dsClient, err = datastore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return err
	}

	return nil
}

func GetDsClient() (*datastore.Client, error) {
	if dsClient != nil {
		return dsClient, nil
	}

	if err := createClient(); err != nil {
		return nil, err
	}

	return dsClient, nil
}
