package transformers

import (
	"github.com/lyft/flyteadmin/pkg/repositories/models"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
)

func CreateNamedEntityMetadataModel(request *admin.NamedEntityMetadataUpdateRequest) models.NamedEntityMetadata {
	return models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: request.ResourceType,
			Project:      request.Id.Project,
			Domain:       request.Id.Domain,
			Name:         request.Id.Name,
		},
		NamedEntityMetadataFields: models.NamedEntityMetadataFields{
			Description: request.Metadata.Description,
		},
	}
}

func FromNamedEntityMetadataFields(metadata models.NamedEntityMetadataFields) admin.NamedEntityMetadata {
	return admin.NamedEntityMetadata{
		Description: metadata.Description,
	}
}

func FromNamedEntityMetadataModel(metadataModel models.NamedEntityMetadata) admin.NamedEntityMetadata {
	NamedEntityMetadata := admin.NamedEntityMetadata{
		Description: metadataModel.Description,
	}
	return NamedEntityMetadata
}
