package transformers

import (
	"context"

	"github.com/lyft/flyteadmin/pkg/common"
	"github.com/lyft/flytestdlib/logger"

	"github.com/golang/protobuf/proto"

	"github.com/golang/protobuf/ptypes"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/lyft/flyteadmin/pkg/errors"
	"github.com/lyft/flyteadmin/pkg/repositories/models"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"google.golang.org/grpc/codes"
)

type ToNodeExecutionModelInput struct {
	Request               *admin.NodeExecutionEventRequest
	ParentTaskExecutionID uint
}

func addNodeRunningState(request *admin.NodeExecutionEventRequest, nodeExecutionModel *models.NodeExecution,
	closure *admin.NodeExecutionClosure) error {
	occurredAt, err := ptypes.Timestamp(request.Event.OccurredAt)
	if err != nil {
		return errors.NewFlyteAdminErrorf(codes.Internal, "failed to unmarshal occurredAt with error: %v", err)
	}

	nodeExecutionModel.StartedAt = &occurredAt
	startedAtProto, err := ptypes.TimestampProto(occurredAt)
	if err != nil {
		return errors.NewFlyteAdminErrorf(codes.Internal,
			"failed to marshal occurredAt into a timestamp proto with error: %v", err)
	}
	closure.StartedAt = startedAtProto
	return nil
}

func addTerminalState(
	request *admin.NodeExecutionEventRequest, nodeExecutionModel *models.NodeExecution,
	closure *admin.NodeExecutionClosure) error {
	if closure.StartedAt == nil {
		logger.Warning(context.Background(), "node execution is missing StartedAt")
	} else {
		endTime, err := ptypes.Timestamp(request.Event.OccurredAt)
		if err != nil {
			return errors.NewFlyteAdminErrorf(
				codes.Internal, "Failed to parse node execution occurred at timestamp: %v", err)
		}
		nodeExecutionModel.Duration = endTime.Sub(*nodeExecutionModel.StartedAt)
		closure.Duration = ptypes.DurationProto(nodeExecutionModel.Duration)
	}

	// Serialize output results (if they exist)
	if request.Event.GetOutputUri() != "" {
		closure.OutputResult = &admin.NodeExecutionClosure_OutputUri{
			OutputUri: request.Event.GetOutputUri(),
		}
	} else if request.Event.GetError() != nil {
		closure.OutputResult = &admin.NodeExecutionClosure_Error{
			Error: request.Event.GetError(),
		}
	}
	return nil
}

func CreateNodeExecutionModel(input ToNodeExecutionModelInput) (*models.NodeExecution, error) {
	nodeExecution := &models.NodeExecution{
		NodeExecutionKey: models.NodeExecutionKey{
			NodeID: input.Request.Event.Id.NodeId,
			ExecutionKey: models.ExecutionKey{
				Project: input.Request.Event.Id.ExecutionId.Project,
				Domain:  input.Request.Event.Id.ExecutionId.Domain,
				Name:    input.Request.Event.Id.ExecutionId.Name,
			},
		},
		Phase:    input.Request.Event.Phase.String(),
		InputURI: input.Request.Event.InputUri,
	}

	closure := admin.NodeExecutionClosure{
		Phase:     input.Request.Event.Phase,
		CreatedAt: input.Request.Event.OccurredAt,
		UpdatedAt: input.Request.Event.OccurredAt,
	}

	if input.Request.Event.Phase == core.NodeExecution_RUNNING {
		err := addNodeRunningState(input.Request, nodeExecution, &closure)
		if err != nil {
			return nil, err
		}
	}
	if common.IsNodeExecutionTerminal(input.Request.Event.Phase) {
		err := addTerminalState(input.Request, nodeExecution, &closure)
		if err != nil {
			return nil, err
		}
	}
	marshaledClosure, err := proto.Marshal(&closure)
	if err != nil {
		return nil, errors.NewFlyteAdminErrorf(
			codes.Internal, "failed to marshal node execution closure with error: %v", err)
	}
	nodeExecution.Closure = marshaledClosure
	nodeExecutionCreatedAt, err := ptypes.Timestamp(input.Request.Event.OccurredAt)
	if err != nil {
		return nil, errors.NewFlyteAdminErrorf(codes.Internal, "failed to read event timestamp")
	}
	nodeExecution.NodeExecutionCreatedAt = &nodeExecutionCreatedAt
	nodeExecution.NodeExecutionUpdatedAt = &nodeExecutionCreatedAt
	if input.Request.Event.ParentTaskMetadata != nil {
		nodeExecution.ParentTaskExecutionID = input.ParentTaskExecutionID
	}
	return nodeExecution, nil
}

func UpdateNodeExecutionModel(
	request *admin.NodeExecutionEventRequest, nodeExecutionModel *models.NodeExecution,
	targetExecution *core.WorkflowExecutionIdentifier) error {
	var nodeExecutionClosure admin.NodeExecutionClosure
	err := proto.Unmarshal(nodeExecutionModel.Closure, &nodeExecutionClosure)
	if err != nil {
		return errors.NewFlyteAdminErrorf(codes.Internal,
			"failed to unmarshal node execution closure with error: %+v", err)
	}
	nodeExecutionModel.Phase = request.Event.Phase.String()
	nodeExecutionClosure.Phase = request.Event.Phase
	nodeExecutionClosure.UpdatedAt = request.Event.OccurredAt
	if request.Event.Phase == core.NodeExecution_RUNNING {
		err := addNodeRunningState(request, nodeExecutionModel, &nodeExecutionClosure)
		if err != nil {
			return err
		}
	}
	if common.IsNodeExecutionTerminal(request.Event.Phase) {
		err := addTerminalState(request, nodeExecutionModel, &nodeExecutionClosure)
		if err != nil {
			return err
		}
	}

	// If the node execution kicked off a workflow execution update the closure if it wasn't set
	if targetExecution != nil && nodeExecutionClosure.GetWorkflowNodeMetadata() == nil {
		nodeExecutionClosure.TargetMetadata = &admin.NodeExecutionClosure_WorkflowNodeMetadata{
			WorkflowNodeMetadata: &admin.WorkflowNodeMetadata{
				ExecutionId: targetExecution,
			},
		}
	}

	marshaledClosure, err := proto.Marshal(&nodeExecutionClosure)
	if err != nil {
		return errors.NewFlyteAdminErrorf(
			codes.Internal, "failed to marshal node execution closure with error: %v", err)
	}
	nodeExecutionModel.Closure = marshaledClosure
	updatedAt, err := ptypes.Timestamp(request.Event.OccurredAt)
	if err != nil {
		return errors.NewFlyteAdminErrorf(codes.Internal, "failed to parse updated at timestamp")
	}
	nodeExecutionModel.NodeExecutionUpdatedAt = &updatedAt
	return nil
}

func FromNodeExecutionModel(nodeExecutionModel models.NodeExecution) (*admin.NodeExecution, error) {
	var closure admin.NodeExecutionClosure
	err := proto.Unmarshal(nodeExecutionModel.Closure, &closure)
	if err != nil {
		return nil, errors.NewFlyteAdminErrorf(codes.Internal, "failed to unmarshal closure")
	}

	return &admin.NodeExecution{
		Id: &core.NodeExecutionIdentifier{
			NodeId: nodeExecutionModel.NodeID,
			ExecutionId: &core.WorkflowExecutionIdentifier{
				Project: nodeExecutionModel.NodeExecutionKey.ExecutionKey.Project,
				Domain:  nodeExecutionModel.NodeExecutionKey.ExecutionKey.Domain,
				Name:    nodeExecutionModel.NodeExecutionKey.ExecutionKey.Name,
			},
		},
		InputUri: nodeExecutionModel.InputURI,
		Closure:  &closure,
	}, nil
}

func FromNodeExecutionModels(
	nodeExecutionModels []models.NodeExecution) ([]*admin.NodeExecution, error) {
	nodeExecutions := make([]*admin.NodeExecution, len(nodeExecutionModels))
	for idx, nodeExecutionModel := range nodeExecutionModels {
		nodeExecution, err := FromNodeExecutionModel(nodeExecutionModel)
		if err != nil {
			return nil, err
		}
		nodeExecutions[idx] = nodeExecution
	}
	return nodeExecutions, nil
}
