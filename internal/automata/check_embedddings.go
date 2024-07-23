package automata

import (
	"errors"
	"fmt"

	"cloud.google.com/go/storage"
)

func (c *Client) CheckEmbeddings(id string) (bool, error) {
	object := fmt.Sprintf("%s.json", id)
	if _, err := c.storage.Bucket(c.embeddings.Bucket).Object(object).Attrs(c.ctx); err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object: %w", err)
	}
	return true, nil
}
