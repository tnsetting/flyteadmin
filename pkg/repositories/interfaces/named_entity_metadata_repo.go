package interfaces

import (
	"context"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/lyft/flyteadmin/pkg/repositories/models"
)

type GetNamedEntityMetadataInput struct {
	ResourceType core.ResourceType
	Project      string
	Domain       string
	Name         string
}

// Defines the interface for interacting with NamedEntityMetadata models
type NamedEntityMetadataRepoInterface interface {
	// Updates metadata associated with a NamedEntity
	Update(ctx context.Context, input models.NamedEntityMetadata) error
	// Gets metadata (if available) associated with a NamedEntity
	Get(ctx context.Context, input GetNamedEntityMetadataInput) (models.NamedEntityMetadata, error)
}