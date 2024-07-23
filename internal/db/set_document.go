package db

import "fmt"

func (d *DB) SetDocument(data Document) error {
	_, err := d.client.Collection(d.collection).Doc(data.ID).Set(d.ctx, data)
	if err != nil {
		return fmt.Errorf("error saving document %w", err)
	}
	return nil
}
