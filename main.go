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
	ProjectID  string `envconfig:"GCP_PROJECT_ID"`
	Region     string `envconfig:"GCP_REGION"`
	Collection string `envconfig:"FIRESTORE_COLLECTION" default:"image-data"`
}

var config Config
var database *db.DB
var client *automata.Client

func init() {
	var err error
	envconfig.MustProcess("", &config)

	database, err = db.NewDB(config.ProjectID, config.Collection)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database")
	}

	client, err = automata.NewClient(config.ProjectID, config.Region)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create client")
	}

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
}

func main() {
	documents, err := database.GetDocuments(20000)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get documents")
	}

	for _, document := range documents {
		imageURI := fmt.Sprintf("gs://%s/%s", document.Bucket, document.Path)
		labels, err := client.GetLabels(imageURI)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get labels")
		}

		description, err := client.GetDescription(imageURI, labels)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get description")
		}

		colors, err := client.GetColors(imageURI)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get colors")
		}

		embeddings, err := client.GetEmbeddings(imageURI, description, labels)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get embeddings")
		}

		document.Description = description
		document.ImageEmbeddings = embeddings.Predictions[0].ImageEmbedding
		document.TextEmbeddings = embeddings.Predictions[0].TextEmbedding
		document.Metadata.Colors = colors
		document.Metadata.Labels = labels

		if err := database.SetDocument(document); err != nil {
			log.Error().Err(err).Msg("Failed to set document")
		}
	}

}
