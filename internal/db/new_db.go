package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// DB is a client for the firestore database.
// It contains the firestore client and the context.
// The client is used to interact with the firestore database.
// The context is used to manage the lifecycle of the client.
// The collection is the name of the collection in the firestore database.
type DB struct {
	client     *firestore.Client
	ctx        context.Context
	collection string
}

// NewDB creates a new client for the firestore database.
func NewDB(projectID, collection string) (*DB, error) {
	ctx := context.TODO()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}

	return &DB{
		client:     client,
		ctx:        ctx,
		collection: collection,
	}, nil
}
