package mocks

import "github.com/lyft/flyteadmin/pkg/executioncluster"

type GetTargetFunc func(*executioncluster.ExecutionTargetSpec) (*executioncluster.ExecutionTarget, error)
type GetAllValidTargetsFunc func() []executioncluster.ExecutionTarget

type MockCluster struct {
	getTargetFunc          GetTargetFunc
	getAllValidTargetsFunc GetAllValidTargetsFunc
}

func (m *MockCluster) SetGetTargetCallback(getTargetFunc GetTargetFunc) {
	m.getTargetFunc = getTargetFunc
}

func (m *MockCluster) SetGetAllValidTargetsCallback(getAllValidTargetsFunc GetAllValidTargetsFunc) {
	m.getAllValidTargetsFunc = getAllValidTargetsFunc
}

func (m *MockCluster) GetTarget(execCluster *executioncluster.ExecutionTargetSpec) (*executioncluster.ExecutionTarget, error) {
	if m.getTargetFunc != nil {
		return m.getTargetFunc(execCluster)
	}
	return nil, nil
}

func (m *MockCluster) GetAllValidTargets() []executioncluster.ExecutionTarget {
	if m.getAllValidTargetsFunc != nil {
		return m.getAllValidTargetsFunc()
	}
	return nil
}
