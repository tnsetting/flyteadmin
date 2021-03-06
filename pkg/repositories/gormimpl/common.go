package gormimpl

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"

	"github.com/lyft/flyteadmin/pkg/common"
	adminErrors "github.com/lyft/flyteadmin/pkg/errors"
	"github.com/lyft/flyteadmin/pkg/repositories/errors"
	"github.com/lyft/flyteadmin/pkg/repositories/interfaces"
	"github.com/lyft/flyteadmin/pkg/repositories/models"
)

const Project = "project"
const Domain = "domain"
const Name = "name"
const Version = "version"
const Closure = "closure"

const ProjectID = "project_id"
const ProjectName = "project_name"
const DomainID = "domain_id"
const DomainName = "domain_name"

const executionTableName = "executions"
const nodeExecutionTableName = "node_executions"
const nodeExecutionEventTableName = "node_event_executions"
const taskExecutionTableName = "task_executions"
const taskTableName = "tasks"

const limit = "limit"
const filters = "filters"

var identifierGroupBy = fmt.Sprintf("%s, %s, %s", Project, Domain, Name)

var innerJoinNodeExecToNodeEvents = fmt.Sprintf(
	"INNER JOIN %s ON %s.node_execution_id = %s.id",
	nodeExecutionTableName, nodeExecutionEventTableName, nodeExecutionTableName)

var innerJoinExecToNodeExec = fmt.Sprintf(
	"INNER JOIN %s ON %s.execution_project = %s.execution_project AND "+
		"%s.execution_domain = %s.execution_domain AND %s.execution_name = %s.execution_name",
	executionTableName, nodeExecutionTableName, executionTableName, nodeExecutionTableName, executionTableName,
	nodeExecutionTableName, executionTableName)

var innerJoinNodeExecToTaskExec = fmt.Sprintf(
	"INNER JOIN %s ON %s.node_id = %s.node_id AND %s.execution_project = %s.execution_project AND "+
		"%s.execution_domain = %s.execution_domain AND %s.execution_name = %s.execution_name",
	nodeExecutionTableName, taskExecutionTableName, nodeExecutionTableName, taskExecutionTableName,
	nodeExecutionTableName, taskExecutionTableName, nodeExecutionTableName, taskExecutionTableName,
	nodeExecutionTableName)

// Because dynamic tasks do NOT necessarily register static task definitions, we use a left join to not exclude
// dynamic tasks from list queries.
var leftJoinTaskToTaskExec = fmt.Sprintf(
	"LEFT JOIN %s ON %s.project = %s.project AND %s.domain = %s.domain AND %s.name = %s.name AND "+
		"%s.version = %s.version",
	taskTableName, taskExecutionTableName, taskTableName, taskExecutionTableName, taskTableName,
	taskExecutionTableName, taskTableName, taskExecutionTableName, taskTableName)

var entityToModel = map[common.Entity]interface{}{
	common.Execution:          models.Execution{},
	common.LaunchPlan:         models.LaunchPlan{},
	common.NodeExecution:      models.NodeExecution{},
	common.NodeExecutionEvent: models.NodeExecutionEvent{},
	common.Task:               models.Task{},
	common.TaskExecution:      models.TaskExecution{},
	common.Workflow:           models.Workflow{},
}

// Validates there are no missing but required parameters in ListResourceInput
func ValidateListInput(input interfaces.ListResourceInput) adminErrors.FlyteAdminError {
	if input.Limit == 0 {
		return errors.GetInvalidInputError(limit)
	}
	if len(input.InlineFilters) == 0 {
		return errors.GetInvalidInputError(filters)
	}
	return nil
}

func applyFilters(tx *gorm.DB, inlineFilters []common.InlineFilter, mapFilters []common.MapFilter) (*gorm.DB, error) {
	for _, filter := range inlineFilters {
		gormQueryExpr, err := filter.GetGormQueryExpr()
		if err != nil {
			return nil, errors.GetInvalidInputError(err.Error())
		}
		tx = tx.Where(gormQueryExpr.Query, gormQueryExpr.Args)
	}
	for _, mapFilter := range mapFilters {
		tx = tx.Where(mapFilter.GetFilter())
	}
	return tx, nil
}

func applyScopedFilters(tx *gorm.DB, inlineFilters []common.InlineFilter, mapFilters []common.MapFilter) (*gorm.DB, error) {
	for _, filter := range inlineFilters {
		entityModel, ok := entityToModel[filter.GetEntity()]
		if !ok {
			return nil, adminErrors.NewFlyteAdminErrorf(codes.InvalidArgument,
				"unrecognized entity in filter expression: %v", filter.GetEntity())
		}
		tableName := tx.NewScope(entityModel).TableName()
		gormQueryExpr, err := filter.GetGormJoinTableQueryExpr(tableName)
		if err != nil {
			return nil, err
		}
		tx = tx.Where(gormQueryExpr.Query, gormQueryExpr.Args)
	}
	for _, mapFilter := range mapFilters {
		tx = tx.Where(mapFilter.GetFilter())
	}
	return tx, nil
}
