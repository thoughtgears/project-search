package db

import (
	"errors"
	"fmt"

	"google.golang.org/api/iterator"
)

// GetDocuments retrieves documents from the database.
// It returns a slice of Document structs.
// If an error occurs, it returns an error.
func (d *DB) GetDocuments(limit int) ([]Document, error) {
	var docs []Document

	itr := d.client.Collection(d.collection).Limit(limit).Documents(d.ctx)
	for {
		doc, err := itr.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("itr.Next: %w", err)
		}
		var document Document
		if err := doc.DataTo(&document); err != nil {
			return nil, fmt.Errorf("doc.DataTo: %w", err)
		}
		document.ID = doc.Ref.ID
		docs = append(docs, document)
	}

	return docs, nil
}
