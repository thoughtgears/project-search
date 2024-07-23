package db

import "time"

// Document represents a document for images.
// The ID is the ID of the image.
// The Bucket is the bucket of the image.
// The Name is the name of the image.
// The Path is the path of the image, including name.
// The URL is the URL of the image, its a public URL.
// The Description is the description of the image.
// The Published is the status of the image, if it is published or not.
// The Valid is the status of the image, if it is valid or not.
// The TimeCreated is the time the image was created.
// The TimeUpdated is the time the image was updated.
type Document struct {
	ID              string    `firestore:"imageId" json:"id"`
	Bucket          string    `firestore:"bucket" json:"-"`
	Name            string    `firestore:"imageName" json:"-"`
	Path            string    `firestore:"imagePath" json:"-"`
	URL             string    `firestore:"imageUrl" json:"url"`
	Description     string    `firestore:"imageDescription,omitempty" json:"description"`
	Published       bool      `firestore:"published,omitempty" json:"-"`
	Valid           bool      `firestore:"valid,omitempty" json:"-"`
	TimeCreated     time.Time `firestore:"timeCreated" json:"-"`
	TimeUpdated     time.Time `firestore:"timeUpdated" json:"-"`
	Metadata        Metadata  `firestore:"metadata,omitempty" json:"metadata"`
	TextEmbeddings  []float64 `firestore:"textEmbeddings,omitempty" json:"text_embeddings"`
	ImageEmbeddings []float64 `firestore:"imageEmbeddings,omitempty" json:"image_embeddings"`
}

// Metadata represents the metadata of the image.
// The labels are the labels of the image. ie. cat, dog, etc.
// The colors are the colors of the image. Based on the Color struct.
// The width is the width of the image. ie. 100, 200, etc.
// The height is the height of the image. ie. 100, 200, etc.
type Metadata struct {
	Labels []string `firestore:"labels,omitempty" json:"labels"`
	Colors []Color  `firestore:"colors,omitempty" json:"colors"`
	Width  int      `firestore:"width" json:"-"`
	Height int      `firestore:"height" json:"-"`
}

// Color represents a color with its name, shade, and weight.
// The name is the name of the color, ie. red, blue, green, etc.
// The shade is the shade of the color, ie. light, dark, etc.
// The weight is the weight of the color in float numbers.
type Color struct {
	Name   string  `firestore:"name" json:"name"`
	Shade  string  `firestore:"shade" json:"shade"`
	Weight float32 `firestore:"weight" json:"weight"`
}
