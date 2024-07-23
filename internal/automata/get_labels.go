package automata

import (
	"fmt"

	vision "cloud.google.com/go/vision/apiv1"
)

// GetLabels returns the labels of the image.
// It uses the Vision API to detect the labels.
// The function returns a slice of strings containing the labels.
// If an error occurs, it returns an error.
// The function takes the imageURI as a parameter, where the uri has to start with "gs://".
func (c *Client) GetLabels(imageURI string) ([]string, error) {
	image := vision.NewImageFromURI(imageURI)
	detectedLabels, err := c.vision.DetectLabels(c.ctx, image, nil, 10)
	if err != nil {
		return nil, fmt.Errorf("DetectLabels: %v", err)
	}

	var labels []string

	for _, l := range detectedLabels {
		labels = append(labels, l.Description)
	}

	return labels, nil
}
