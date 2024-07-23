package main

import (
	"fmt"

	"github.com/thoughtgears/project-search/internal/automata"
	"github.com/thoughtgears/project-search/internal/db"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ProjectID    string `envconfig:"GCP_PROJECT_ID"`
	Region       string `envconfig:"GCP_REGION"`
	Collection   string `envconfig:"FIRESTORE_COLLECTION" default:"image-data"`
	VertexBucket string `envconfig:"VERTEX_BUCKET"`
}

var config Config

func init() {
	envconfig.MustProcess("", &config)

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
}

func main() {
	database, err := db.NewDB(config.ProjectID, config.Collection)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database")
	}

	documents, err := database.GetDocuments(2000)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get documents")
	}

	for _, document := range documents {
		imageURI := fmt.Sprintf("gs://%s/%s", document.Bucket, document.Path)
		client, err := automata.NewClient(config.ProjectID, config.Region, config.VertexBucket)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create client")
		}

		exists, err := client.CheckEmbeddings(document.ID)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to check if document exists")
		}

		if exists {
			log.Info().Interface("ID", document.ID).Msg("Document already exists in GCS")
			continue
		}

		labels, err := client.GetLabels(imageURI)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get labels")
		}

		colors, err := client.GetColors(imageURI)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get colors")
		}

		description, err := client.GetDescription(imageURI, labels)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get description")
		}

		embeddings, err := client.GetEmbeddings(description, labels, imageURI)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get embeddings")
		}

		document.Description = description
		document.TextEmbeddings = embeddings.Predictions[0].TextEmbedding
		document.ImageEmbeddings = embeddings.Predictions[0].ImageEmbedding
		document.Metadata.Labels = labels
		document.Metadata.Colors = colors

		if err := client.PutEmbeddings(document); err != nil {
			log.Fatal().Err(err).Msg("Failed to put embeddings")
		}
	}
}
