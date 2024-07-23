package automata

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/vertexai/genai"
	vision "cloud.google.com/go/vision/apiv1"
)

// Client is a client for the automata service.
// It contains the genai and vision clients.
// The genai client is used to generate content.
// The vision client is used to detect labels and colors in images.
type Client struct {
	genai      *genai.Client
	vision     *vision.ImageAnnotatorClient
	embeddings Embeddings
	ctx        context.Context
	modelName  string
	projectID  string
}

// Embeddings is a struct that contains the client, URL.
// The client is used to make HTTP requests.
// The URL is the URL of the embeddings model.
// Its used to make embeddings requests.
type Embeddings struct {
	client *http.Client
	URL    string
}

// NewClient creates a new client for the automata service.
// It returns a client and an error if the client could not be created.
// The function takes the projectID and region as parameters.
func NewClient(projectID, region string) (*Client, error) {
	ctx := context.TODO()

	genaiClient, err := genai.NewClient(ctx, projectID, region)
	if err != nil {
		return nil, fmt.Errorf("genai.NewClient: %v", err)
	}

	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("vision.NewImageAnnotatorClient: %v", err)
	}

	return &Client{
		genai:  genaiClient,
		vision: visionClient,
		embeddings: Embeddings{
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
			URL: fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:predict", region, projectID, region, "multimodalembedding@001"),
		},
		ctx:       ctx,
		modelName: "gemini-1.5-flash-001",
		projectID: projectID,
	}, nil
}
