package automata

import (
	"encoding/json"
	"fmt"

	"github.com/thoughtgears/project-search/internal/db"

	"cloud.google.com/go/vertexai/genai"
	vision "cloud.google.com/go/vision/apiv1"
)

// GetColors returns the colors of the image.
// It uses the Vision API to detect the colors.
// The function returns a slice of Color containing the colors.
// If an error occurs, it returns an error.
// The function takes the imageURI as a parameter, where the uri has to start with "gs://".
func (c *Client) GetColors(imageURI string) ([]db.Color, error) {
	image := vision.NewImageFromURI(imageURI)
	detectedColors, err := c.vision.DetectImageProperties(c.ctx, image, nil)
	if err != nil {
		return nil, fmt.Errorf("DetectImageProperties: %v", err)
	}

	var colors string
	for _, color := range detectedColors.DominantColors.Colors {
		colors += fmt.Sprintf("RGB: %s, weight: %v -", color.Color.String(), color.Score)
	}

	return c.parseRGB(colors)
}

// parseRGB parses the colors from the string.
// It uses the GenAI API to parse the colors.
// The function returns a slice of Color containing the colors.
// If an error occurs, it returns an error.
func (c *Client) parseRGB(colors string) ([]db.Color, error) {
	model := c.genai.GenerativeModel(c.modelName)
	model.SetMaxOutputTokens(1024)
	model.SetTemperature(0.1)
	model.ResponseMIMEType = "application/json"
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text("You are a computer who should calculate the colors and the color weight.")},
	}

	parts := []genai.Part{
		genai.Text("Based on the input colors, give me the colors that the RGB values make."),
		genai.Text("A color MUST be a single word, i.e. 'red', 'blue', 'green', etc."),
		genai.Text("The response must be a list of objects with the following keys:"),
		genai.Text("name: name of color, shade: light, dark etc, weight: weight of color in float numbers"),
		genai.Text("The color MUST only appear once in the list with shade combinations."),
		genai.Text("The colors are:"),
		genai.Text(colors),
	}

	resp, err := model.GenerateContent(c.ctx, parts...)
	if err != nil {
		return nil, fmt.Errorf("GenerateContent: %v", err)
	}

	var parsedColors []db.Color

	if len(resp.Candidates) != 0 {
		part := resp.Candidates[0].Content.Parts[0]
		switch p := part.(type) {
		case genai.Text:
			if err := json.Unmarshal([]byte(p), &parsedColors); err != nil {
				return nil, fmt.Errorf("json.Unmarshal: %v", err)
			}
		default:
			return nil, nil
		}
	}

	return parsedColors, nil
}
