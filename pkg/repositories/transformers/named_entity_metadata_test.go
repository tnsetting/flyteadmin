package transformers

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/lyft/flyteadmin/pkg/repositories/models"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/stretchr/testify/assert"
)

func TestCreateNamedEntityMetadataModel(t *testing.T) {

	model := CreateNamedEntityMetadataModel(&admin.NamedEntityMetadataUpdateRequest{
		ResourceType: core.ResourceType_WORKFLOW,
		Id: &admin.NamedEntityIdentifier{
			Project: "project",
			Domain:  "domain",
			Name:    "name",
		},
		Metadata: &admin.NamedEntityMetadata{
			Description: "description",
		},
	})

	assert.Equal(t, models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: core.ResourceType_WORKFLOW,
			Project:      "project",
			Domain:       "domain",
			Name:         "name",
		},
		NamedEntityMetadataFields: models.NamedEntityMetadataFields{
			Description: "description",
		},
	}, model)
}

func TestFromNamedEntityMetadataModel(t *testing.T) {
	model := models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: core.ResourceType_WORKFLOW,
			Project:      "project",
			Domain:       "domain",
			Name:         "name",
		},
		NamedEntityMetadataFields: models.NamedEntityMetadataFields{
			Description: "description",
		},
	}

	metadata := FromNamedEntityMetadataModel(model)
	assert.True(t, proto.Equal(&admin.NamedEntityMetadata{
		Description: "description",
	}, &metadata))
}

func TestFromNamedEntityMetadataFields(t *testing.T) {
	model := models.NamedEntityMetadataFields{
		Description: "description",
	}

	metadata := FromNamedEntityMetadataFields(model)
	assert.True(t, proto.Equal(&admin.NamedEntityMetadata{
		Description: "description",
	}, &metadata))
}
