package adminservice

import (
	"context"

	"github.com/lyft/flyteadmin/pkg/rpc/adminservice/util"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *AdminService) GetNamedEntityMetadata(ctx context.Context, request *admin.GetNamedEntityMetadataRequest) (*admin.NamedEntityMetadata, error) {
	defer m.interceptPanic(ctx, request)
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Incorrect request, nil requests not allowed")
	}

	// We require a valid resource type to retrieve metadata
	if request.ResourceType == core.ResourceType_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "resource_type cannot be undefined")
	}

	var response *admin.NamedEntityMetadata
	var err error
	m.Metrics.namedEntityMetadataEndpointMetrics.get.Time(func() {
		response, err = m.NamedEntityMetadataManager.GetNamedEntityMetadata(ctx, *request)
	})
	if err != nil {
		return nil, util.TransformAndRecordError(err, &m.Metrics.namedEntityMetadataEndpointMetrics.get)
	}
	m.Metrics.namedEntityMetadataEndpointMetrics.get.Success()
	return response, nil

}

func (m *AdminService) UpdateNamedEntityMetadata(ctx context.Context, request *admin.NamedEntityMetadataUpdateRequest) (
	*admin.NamedEntityMetadataUpdateResponse, error) {
	defer m.interceptPanic(ctx, request)
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Incorrect request, nil requests not allowed")
	}

	// We require a valid resource type to look up / set metadata
	if request.ResourceType == core.ResourceType_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "resource_type cannot be undefined")
	}

	var response *admin.NamedEntityMetadataUpdateResponse
	var err error
	m.Metrics.namedEntityMetadataEndpointMetrics.update.Time(func() {
		response, err = m.NamedEntityMetadataManager.UpdateNamedEntityMetadata(ctx, *request)
	})
	if err != nil {
		return nil, util.TransformAndRecordError(err, &m.Metrics.namedEntityMetadataEndpointMetrics.update)
	}
	m.Metrics.namedEntityMetadataEndpointMetrics.update.Success()
	return response, nil
}
