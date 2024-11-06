package initializers

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

type CollectionName string

const FILINGS_VECTOR_SIZE = 1536

const (
	Filings CollectionName = "filings"
)

func CreateQdrantCollections(client *qdrant.Client) error {
	ctx := context.Background()

	exists, err := client.CollectionExists(ctx, string(Filings))
	if err != nil {
		return fmt.Errorf("error checking collection existence: %w", err)
	}

	if exists {
		return nil
	}

	// Collection doesn't exist, create it
	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: string(Filings),
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     FILINGS_VECTOR_SIZE,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		return fmt.Errorf("error creating collection: %w", err)
	}

	return nil
}
