package interfaces

import (
	"context"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
)

// Interface for managing metadata associated with NamedEntityIdentifiers
type NamedEntityMetadataInterface interface {
	GetNamedEntityMetadata(ctx context.Context, request admin.GetNamedEntityMetadataRequest) (*admin.NamedEntityMetadata, error)
	UpdateNamedEntityMetadata(ctx context.Context, request admin.NamedEntityMetadataUpdateRequest) (*admin.NamedEntityMetadataUpdateResponse, error)
}
