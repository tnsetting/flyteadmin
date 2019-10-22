package gormimpl

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/lyft/flyteadmin/pkg/repositories/errors"
	"github.com/lyft/flyteadmin/pkg/repositories/interfaces"
	"github.com/lyft/flyteadmin/pkg/repositories/models"
	"github.com/lyft/flytestdlib/promutils"
)

// Implementation of NamedEntityMetadataRepoInterface.
type NamedEntityMetadataRepo struct {
	db               *gorm.DB
	errorTransformer errors.ErrorTransformer
	metrics          gormMetrics
}

func (r *NamedEntityMetadataRepo) Update(ctx context.Context, input models.NamedEntityMetadata) error {
	timer := r.metrics.UpdateDuration.Start()
	// TODO: Check for existence, create if non-existent. Otherwise update.
	tx := r.db.Update(&input)
	timer.Stop()
	if tx.Error != nil {
		return r.errorTransformer.ToFlyteAdminError(tx.Error)
	}
	return nil
}

func (r *NamedEntityMetadataRepo) Get(ctx context.Context, input interfaces.GetNamedEntityMetadataInput) (models.NamedEntityMetadata, error) {
	var NamedEntityMetadata models.NamedEntityMetadata
	timer := r.metrics.GetDuration.Start()
	tx := r.db.Where(&models.NamedEntityMetadata{
		NamedEntityMetadataKey: models.NamedEntityMetadataKey{
			ResourceType: input.ResourceType,
			Project:      input.Project,
			Domain:       input.Domain,
			Name:         input.Name,
		},
	}).First(&NamedEntityMetadata)
	timer.Stop()
	if tx.Error != nil {
		return models.NamedEntityMetadata{}, r.errorTransformer.ToFlyteAdminError(tx.Error)
	}
	// If a record is not found, we will return empty metadata
	if tx.RecordNotFound() {
		return models.NamedEntityMetadata{}, nil
	}
	return NamedEntityMetadata, nil
}

// Returns an instance of NamedEntityMetadataRepoInterface
func NewNamedEntityMetadataRepo(
	db *gorm.DB, errorTransformer errors.ErrorTransformer, scope promutils.Scope) interfaces.NamedEntityMetadataRepoInterface {
	metrics := newMetrics(scope)
	// TODO: specialized metrics?

	return &NamedEntityMetadataRepo{
		db:               db,
		errorTransformer: errorTransformer,
		metrics:          metrics,
	}
}
