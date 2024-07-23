package automata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/thoughtgears/project-search/internal/db"
)

func (c *Client) PutEmbeddings(data db.Document) error {
	var buffer bytes.Buffer

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	buffer.Write(jsonData)

	objectName := fmt.Sprintf("%s.json", data.ID)

	writer := c.storage.Bucket(c.embeddings.Bucket).Object(objectName).NewWriter(c.ctx)
	writer.ContentType = "application/json"

	if _, err := io.Copy(writer, &buffer); err != nil {
		return fmt.Errorf("failed to write to bucket: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	log.Println("Successfully wrote JSON lines to", objectName)
	return nil
}
