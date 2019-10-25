package gormimpl

import (
	"context"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/lyft/flyteadmin/pkg/repositories/errors"
	"github.com/lyft/flyteadmin/pkg/repositories/interfaces"
	"github.com/lyft/flyteadmin/pkg/repositories/models"
	mockScope "github.com/lyft/flytestdlib/promutils"
	"github.com/stretchr/testify/assert"
)

func getMockNamedEntityMetadataResponseFromDb(expected models.NamedEntityMetadata) map[string]interface{} {
	metadata := make(map[string]interface{})
	metadata["resource_type"] = expected.ResourceType
	metadata["project"] = expected.Project
	metadata["domain"] = expected.Domain
	metadata["name"] = expected.Name
	metadata["description"] = expected.Description
	return metadata
}

func TestGetNamedEntityMetadata(t *testing.T) {
	metadataRepo := NewNamedEntityMetadataRepo(GetDbForTest(t), errors.NewTestErrorTransformer(), mockScope.NewTestScope())

	results := make([]map[string]interface{}, 0)
	metadata := getMockNamedEntityMetadataResponseFromDb(models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: resourceType,
			Project:      project,
			Domain:       domain,
			Name:         name,
		},
		NamedEntityMetadataFields: models.NamedEntityMetadataFields{
			Description: description,
		},
	})
	results = append(results, metadata)

	GlobalMock := mocket.Catcher.Reset()
	GlobalMock.Logging = true
	GlobalMock.NewMock().WithQuery(
		`SELECT * FROM "named_entity_metadata"  WHERE "named_entity_metadata"."deleted_at" IS NULL AND (("named_entity_metadata"."resource_type" = 2) AND ("named_entity_metadata"."project" = project) AND ("named_entity_metadata"."domain" = domain) AND ("named_entity_metadata"."name" = name)) ORDER BY "named_entity_metadata"."id" ASC LIMIT 1`).WithReply(results)
	output, err := metadataRepo.Get(context.Background(), interfaces.GetNamedEntityMetadataInput{
		ResourceType: resourceType,
		Project:      project,
		Domain:       domain,
		Name:         name,
	})
	assert.NoError(t, err)
	assert.Equal(t, project, output.Project)
	assert.Equal(t, domain, output.Domain)
	assert.Equal(t, name, output.Name)
	assert.Equal(t, resourceType, output.ResourceType)
	assert.Equal(t, description, output.Description)
}

func TestUpdateNamedEntityMetadata_WithExisting(t *testing.T) {
	metadataRepo := NewNamedEntityMetadataRepo(GetDbForTest(t), errors.NewTestErrorTransformer(), mockScope.NewTestScope())
	const updatedDescription = "updated description"

	results := make([]map[string]interface{}, 0)
	metadata := getMockNamedEntityMetadataResponseFromDb(models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: resourceType,
			Project:      project,
			Domain:       domain,
			Name:         name,
		},
		NamedEntityMetadataFields: models.NamedEntityMetadataFields{
			Description: description,
		},
	})
	results = append(results, metadata)

	GlobalMock := mocket.Catcher.Reset()
	GlobalMock.Logging = true
	GlobalMock.NewMock().WithQuery(
		`SELECT * FROM "named_entity_metadata"  WHERE "named_entity_metadata"."deleted_at" IS NULL AND (("named_entity_metadata"."resource_type" = 2) AND ("named_entity_metadata"."project" = project) AND ("named_entity_metadata"."domain" = domain) AND ("named_entity_metadata"."name" = name)) ORDER BY "named_entity_metadata"."id" ASC LIMIT 1`).WithReply(results)

	mockQuery := GlobalMock.NewMock()
	mockQuery.WithQuery(
		`UPDATE "named_entity_metadata" SET "description" = ?, "updated_at" = ?  WHERE "named_entity_metadata"."deleted_at" IS NULL AND (("named_entity_metadata"."resource_type" = ?) AND ("named_entity_metadata"."project" = ?) AND ("named_entity_metadata"."domain" = ?) AND ("named_entity_metadata"."name" = ?))`)

	err := metadataRepo.Update(context.Background(), models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: resourceType,
			Project:      project,
			Domain:       domain,
			Name:         name,
		},
		NamedEntityMetadataFields: models.NamedEntityMetadataFields{
			Description: updatedDescription,
		},
	})
	assert.NoError(t, err)
	assert.True(t, mockQuery.Triggered)
}

func TestUpdateNamedEntityMetadata_CreateNew(t *testing.T) {
	metadataRepo := NewNamedEntityMetadataRepo(GetDbForTest(t), errors.NewTestErrorTransformer(), mockScope.NewTestScope())
	const updatedDescription = "updated description"

	GlobalMock := mocket.Catcher.Reset()
	GlobalMock.Logging = true

	mockQuery := GlobalMock.NewMock()
	mockQuery.WithQuery(
		`INSERT  INTO "named_entity_metadata" ("created_at","updated_at","deleted_at","resource_type","project","domain","name","description") VALUES (?,?,?,?,?,?,?,?)`)

	err := metadataRepo.Update(context.Background(), models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: resourceType,
			Project:      project,
			Domain:       domain,
			Name:         name,
		},
		NamedEntityMetadataFields: models.NamedEntityMetadataFields{
			Description: updatedDescription,
		},
	})
	assert.NoError(t, err)
	assert.True(t, mockQuery.Triggered)
}
