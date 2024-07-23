package automata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thoughtgears/project-search/internal/auth"
)

// EmbeddingsRequestBody represents the structure of the JSON request body
type EmbeddingsRequestBody struct {
	Instances []Instance `json:"instances"`
}

// Instance represents each instance in the request body
// It contains the text, image and parameters
// The text is the description of the image
// The image contains the image data
// The parameters contain the dimension of the embeddings
type Instance struct {
	Text       string     `json:"text"`
	Image      Image      `json:"image"`
	Parameters Parameters `json:"parameters"`
}

// Image represents the image data in the request body
// It contains the base64 encoded image, the GCS URI and the MIME type
// The MIME type should be the type of the image, i.e. image/jpeg
// The GCS URI should start with "gs://"
// The base64 encoded image should be the image encoded in base64
// GCS URI and base64 encoded image are mutually exclusive
type Image struct {
	BytesBase64Encoded string `json:"bytesBase64Encoded,omitempty"`
	GcsUri             string `json:"gcsUri,omitempty"`
	MimeType           string `json:"mimeType,omitempty"`
}

// Parameters represents the parameters in the request body
// It contains the dimension of the embeddings, 128, 256, 512 or 1408
type Parameters struct {
	Dimension int `json:"dimension"`
}

// EmbeddingsResponse represents the structure of the JSON response body
type EmbeddingsResponse struct {
	Predictions []Prediction `json:"predictions"`
}

// Prediction struct represents each prediction in the response
// It contains the image and text embeddings
type Prediction struct {
	ImageEmbedding []float64 `json:"imageEmbedding"`
	TextEmbedding  []float64 `json:"textEmbedding"`
}

// GetEmbeddings returns the embeddings of the image and text
// It uses the GenAI API to get the embeddings
// The function returns an EmbeddingsResponse containing the embeddings
// If an error occurs, it returns an error
// The function takes the description, labels, and imageURI as parameters
// The description is the description of the image
// The labels are the labels of the image
// The imageURI is the URI of the image, where the URI has to start with "gs://"
func (c *Client) GetEmbeddings(description string, labels []string, imageURI string) (EmbeddingsResponse, error) {
	requestBody := EmbeddingsRequestBody{
		Instances: []Instance{
			{
				Text: fmt.Sprintf("%s %s", description, labels),
				Image: Image{
					GcsUri: imageURI,
				},
				Parameters: Parameters{
					Dimension: 512,
				},
			},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return EmbeddingsResponse{}, fmt.Errorf("marshal: %w", err)
	}

	token, err := auth.IDTokenTokenSource(c.ctx, fmt.Sprintf("https://%s-aiplatform.googleapis.com/", c.projectID))
	if err != nil {
		return EmbeddingsResponse{}, fmt.Errorf("IDTokenTokenSource: %w", err)
	}

	req, _ := http.NewRequest("POST", c.embeddings.URL, bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := c.embeddings.client.Do(req)
	if err != nil {
		return EmbeddingsResponse{}, fmt.Errorf("do: %w", err)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return EmbeddingsResponse{}, fmt.Errorf("ReadAll: %w", err)
	}
	var embeddingsResponse EmbeddingsResponse
	if err := json.Unmarshal(responseBody, &embeddingsResponse); err != nil {
		return EmbeddingsResponse{}, fmt.Errorf("unmarshal: %w", err)
	}

	return embeddingsResponse, nil
}
