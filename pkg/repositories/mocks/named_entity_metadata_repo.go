// Mock implementation of a workflow repo to be used for tests.
package mocks

import (
	"context"

	"github.com/lyft/flyteadmin/pkg/repositories/interfaces"
	"github.com/lyft/flyteadmin/pkg/repositories/models"
)

type UpdateNamedEntityMetadataFunc func(input models.NamedEntityMetadata) error
type GetNamedEntityMetadataFunc func(input interfaces.GetNamedEntityMetadataInput) (models.NamedEntityMetadata, error)

type MockNamedEntityMetadataRepo struct {
	updateFunction UpdateNamedEntityMetadataFunc
	getFunction    GetNamedEntityMetadataFunc
}

func (r *MockNamedEntityMetadataRepo) Update(ctx context.Context, NamedEntityMetadata models.NamedEntityMetadata) error {
	if r.updateFunction != nil {
		return r.updateFunction(NamedEntityMetadata)
	}
	return nil
}

func (r *MockNamedEntityMetadataRepo) Get(
	ctx context.Context, input interfaces.GetNamedEntityMetadataInput) (models.NamedEntityMetadata, error) {
	if r.getFunction != nil {
		return r.getFunction(input)
	}
	return models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: input.ResourceType,
			Project:      input.Project,
			Domain:       input.Domain,
			Name:         input.Name,
		},
		NamedEntityMetadataFields: models.NamedEntityMetadataFields{
			Description: "",
		},
	}, nil
}

func (r *MockNamedEntityMetadataRepo) SetGetCallback(getFunction GetNamedEntityMetadataFunc) {
	r.getFunction = getFunction
}

func (r *MockNamedEntityMetadataRepo) SetUpdateCallback(updateFunction UpdateNamedEntityMetadataFunc) {
	r.updateFunction = updateFunction
}

func NewMockNamedEntityMetadataRepo() interfaces.NamedEntityMetadataRepoInterface {
	return &MockNamedEntityMetadataRepo{}
}
