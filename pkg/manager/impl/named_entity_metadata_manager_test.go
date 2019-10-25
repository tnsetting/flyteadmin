package impl

import (
	"context"
	"testing"

	"github.com/lyft/flyteadmin/pkg/manager/impl/testutils"
	"github.com/lyft/flyteadmin/pkg/repositories"
	"github.com/lyft/flyteadmin/pkg/repositories/interfaces"
	repositoryMocks "github.com/lyft/flyteadmin/pkg/repositories/mocks"
	"github.com/lyft/flyteadmin/pkg/repositories/models"
	runtimeInterfaces "github.com/lyft/flyteadmin/pkg/runtime/interfaces"
	runtimeMocks "github.com/lyft/flyteadmin/pkg/runtime/mocks"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	mockScope "github.com/lyft/flytestdlib/promutils"
	"github.com/stretchr/testify/assert"
)

var metadataIdentifier = admin.NamedEntityIdentifier{
	Project: project,
	Domain:  domain,
	Name:    name,
}

var badIdentifier = admin.NamedEntityIdentifier{
	Project: project,
	Domain:  domain,
	Name:    "",
}

func getMockRepositoryForNEMTest() repositories.RepositoryInterface {
	return repositoryMocks.NewMockRepository()
}

func getMockConfigForNEMTest() runtimeInterfaces.Configuration {
	mockConfig := runtimeMocks.NewMockConfigurationProvider(
		testutils.GetApplicationConfigWithDefaultProjects(), nil, nil, nil, nil)
	return mockConfig
}

func TestNamedEntityMetadataManager_GetMetadata(t *testing.T) {
	repository := getMockRepositoryForNEMTest()
	manager := NewNamedEntityMetadataManager(repository, getMockConfigForNEMTest(), mockScope.NewTestScope())

	metadataGetFunction := func(input interfaces.GetNamedEntityMetadataInput) (models.NamedEntityMetadata, error) {
		return models.NamedEntityMetadata{
			NamedEntityMetadataKey: models.NamedEntityMetadataKey{
				ResourceType: input.ResourceType,
				Project:      input.Project,
				Domain:       input.Domain,
				Name:         input.Name,
			},
			NamedEntityMetadataFields: models.NamedEntityMetadataFields{
				Description: description,
			},
		}, nil
	}
	repository.NamedEntityMetadataRepo().(*repositoryMocks.MockNamedEntityMetadataRepo).SetGetCallback(metadataGetFunction)
	response, err := manager.GetNamedEntityMetadata(context.Background(), admin.GetNamedEntityMetadataRequest{
		ResourceType: resourceType,
		Id:           &metadataIdentifier,
	})
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestNamedEntityMetadataManager_GetMetadata_BadRequest(t *testing.T) {
	repository := getMockRepositoryForNEMTest()
	manager := NewNamedEntityMetadataManager(repository, getMockConfigForNEMTest(), mockScope.NewTestScope())

	response, err := manager.GetNamedEntityMetadata(context.Background(), admin.GetNamedEntityMetadataRequest{
		ResourceType: core.ResourceType_UNSPECIFIED,
		Id:           &metadataIdentifier,
	})
	assert.Error(t, err)
	assert.Nil(t, response)

	response, err = manager.GetNamedEntityMetadata(context.Background(), admin.GetNamedEntityMetadataRequest{
		ResourceType: resourceType,
		Id:           &badIdentifier,
	})
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestNamedEntityMetadataManager_Update(t *testing.T) {
	repository := getMockRepositoryForNEMTest()
	manager := NewNamedEntityMetadataManager(repository, getMockConfigForNEMTest(), mockScope.NewTestScope())
	updatedDescription := "updated description"

	metadataUpdateFunction := func(input models.NamedEntityMetadata) error {
		assert.Equal(t, input.Description, updatedDescription)
		assert.Equal(t, input.ResourceType, resourceType)
		assert.Equal(t, input.Project, project)
		assert.Equal(t, input.Domain, domain)
		assert.Equal(t, input.Name, name)
		return nil
	}
	repository.NamedEntityMetadataRepo().(*repositoryMocks.MockNamedEntityMetadataRepo).SetUpdateCallback(metadataUpdateFunction)
	updatedMetadata := admin.NamedEntityMetadata{
		Description: updatedDescription,
	}
	response, err := manager.UpdateNamedEntityMetadata(context.Background(), admin.NamedEntityMetadataUpdateRequest{
		Metadata:     &updatedMetadata,
		ResourceType: resourceType,
		Id:           &metadataIdentifier,
	})
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestNamedEntityMetadataManager_Update_BadRequest(t *testing.T) {
	repository := getMockRepositoryForNEMTest()
	manager := NewNamedEntityMetadataManager(repository, getMockConfigForNEMTest(), mockScope.NewTestScope())
	updatedDescription := "updated description"

	updatedMetadata := admin.NamedEntityMetadata{
		Description: updatedDescription,
	}
	response, err := manager.UpdateNamedEntityMetadata(context.Background(), admin.NamedEntityMetadataUpdateRequest{
		Metadata:     &updatedMetadata,
		ResourceType: core.ResourceType_UNSPECIFIED,
		Id:           &metadataIdentifier,
	})
	assert.Error(t, err)
	assert.Nil(t, response)

	response, err = manager.UpdateNamedEntityMetadata(context.Background(), admin.NamedEntityMetadataUpdateRequest{
		Metadata:     &updatedMetadata,
		ResourceType: resourceType,
		Id:           &badIdentifier,
	})
	assert.Error(t, err)
	assert.Nil(t, response)
}
