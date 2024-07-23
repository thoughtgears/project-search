package automata

import (
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

// GetDescription returns the description of the image.
// It uses the generative model to generate the description.
// The function returns the description as a string.
// If an error occurs, it returns an error.
// The function takes the imageURI and labels as parameters.
// The imageURI is the URI of the image.
// The labels are the labels of the image.
func (c *Client) GetDescription(imageURI string, labels []string) (string, error) {
	model := c.genai.GenerativeModel(c.modelName)
	model.SetTemperature(0.9)
	model.SetTopP(0.95)
	model.SetTopK(20)
	model.SetMaxOutputTokens(512)

	parts := []genai.Part{
		genai.Text("Describe the image with only 150 tokens."),
		genai.Text("The images are related to Checkatrade.com and usually contain home improvement projects."),
		genai.Text("You must assess the images properly to ensure the description matches the image presented."),
		genai.Text("The description MUST be something that could be shown in an ALT IMG tag."),
		genai.Text("You can use the following labels to help you describe the image: " + strings.Join(labels, ", ")),
		genai.FileData{
			MIMEType: "image/jpeg",
			FileURI:  imageURI,
		},
	}

	resp, err := model.GenerateContent(c.ctx, parts...)
	if err != nil {
		return "", fmt.Errorf("GenerateContent: %v", err)
	}

	if len(resp.Candidates) == 0 ||
		len(resp.Candidates[0].Content.Parts) == 0 {
		return "Empty Description", nil
	}

	var description string

	if resp.Candidates[0].FinishReason == genai.FinishReasonSafety {
		return "Finish with SafetyReason", nil
	}

	if len(resp.Candidates) != 0 {
		part := resp.Candidates[0].Content.Parts[0]
		switch p := part.(type) {
		case genai.Text:
			description = string(p)
		default:
			return "Empty Description", nil
		}
	}

	return description, nil
}
