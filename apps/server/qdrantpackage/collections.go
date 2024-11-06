package qdrantpackage

import (
	"context"
	"fmt"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

type CollectionName string

const DEFAULT_VECTOR_SIZE = 1536

const (
	Messages CollectionName = "messages"
)

func CreateQdrantCollections(client *qdrant.Client, collectionNames []CollectionName) error {
	ctx := context.Background()

	for _, collectionName := range collectionNames {
		exists, err := client.CollectionExists(ctx, string(collectionName))
		if err != nil {
			return fmt.Errorf("error checking collection existence: %w", err)
		}

		if exists {
			log.Printf("Collection %s already exists. skipping...", collectionName)
			continue
		}

		// Collection doesn't exist, create it
		err = client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: string(collectionName),
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     DEFAULT_VECTOR_SIZE,
				Distance: qdrant.Distance_Cosine,
			}),
		})

		if err != nil {
			return fmt.Errorf("error creating collection: %w", err)
		}
		log.Printf("Successfully created collection %s", collectionName)
	}

	return nil
}
