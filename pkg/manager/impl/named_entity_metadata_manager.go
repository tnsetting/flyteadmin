package impl

import (
	"context"

	"github.com/lyft/flyteadmin/pkg/manager/impl/util"
	"github.com/lyft/flyteadmin/pkg/manager/impl/validation"
	"github.com/lyft/flyteadmin/pkg/manager/interfaces"
	"github.com/lyft/flyteadmin/pkg/repositories"
	"github.com/lyft/flyteadmin/pkg/repositories/transformers"
	runtimeInterfaces "github.com/lyft/flyteadmin/pkg/runtime/interfaces"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"
	"github.com/lyft/flytestdlib/promutils"
)

type NamedEntityMetadataMetrics struct {
	Scope promutils.Scope
}

type NamedEntityMetadataManager struct {
	db      repositories.RepositoryInterface
	config  runtimeInterfaces.Configuration
	metrics NamedEntityMetadataMetrics
}

func (m *NamedEntityMetadataManager) UpdateNamedEntityMetadata(ctx context.Context, request admin.NamedEntityMetadataUpdateRequest) (
	*admin.NamedEntityMetadataUpdateResponse, error) {
	if err := validation.ValidateResourceType(request.ResourceType); err != nil {
		return nil, err
	}
	if err := validation.ValidateNamedEntityIdentifier(request.Id); err != nil {
		return nil, err
	}

	metadataModel := transformers.CreateNamedEntityMetadataModel(&request)
	err := m.db.NamedEntityMetadataRepo().Update(ctx, metadataModel)
	if err != nil {
		logger.Debugf(ctx, "Failed to update named_entity_metadata for [%+v] with err %v", request.Id, err)
		return nil, err
	}
	return &admin.NamedEntityMetadataUpdateResponse{}, nil
}

func (m *NamedEntityMetadataManager) GetNamedEntityMetadata(ctx context.Context, request admin.GetNamedEntityMetadataRequest) (
	*admin.NamedEntityMetadata, error) {
	if err := validation.ValidateResourceType(request.ResourceType); err != nil {
		return nil, err
	}
	if err := validation.ValidateNamedEntityIdentifier(request.Id); err != nil {
		return nil, err
	}
	return util.GetNamedEntityMetadata(ctx, m.db, request.ResourceType, *request.Id)
}

func NewNamedEntityMetadataManager(
	db repositories.RepositoryInterface,
	config runtimeInterfaces.Configuration,
	scope promutils.Scope) interfaces.NamedEntityMetadataInterface {

	metrics := NamedEntityMetadataMetrics{
		Scope: scope,
	}
	return &NamedEntityMetadataManager{
		db:      db,
		config:  config,
		metrics: metrics,
	}
}
