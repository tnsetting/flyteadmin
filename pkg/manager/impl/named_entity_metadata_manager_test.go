package impl

import (
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
)

const project, domain, name, description = "project", "domain", "name", "description",

var metadataIdentifier = admin.NamedEntityIdentifier{
	ResourceType: core.ResourceType_WORKFLOW,
	Project:      project,
	Domain:       domain,
	Name:         name,
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
	// repository := getMockRepositoryForLpTest()
	// lpManager := NewLaunchPlanManager(repository, getMockConfigForLpTest(), mockScheduler, mockScope.NewTestScope())
	// state := int32(0)
	// lpRequest := testutils.GetLaunchPlanRequest()
	// workflowRequest := testutils.GetWorkflowRequest()

	// closure := admin.LaunchPlanClosure{
	// 	ExpectedInputs:  lpRequest.Spec.DefaultInputs,
	// 	ExpectedOutputs: workflowRequest.Spec.Template.Interface.Outputs,
	// }
	// specBytes, _ := proto.Marshal(lpRequest.Spec)
	// closureBytes, _ := proto.Marshal(&closure)

	// launchPlanGetFunc := func(input interfaces.GetResourceInput) (models.LaunchPlan, error) {
	// 	return models.LaunchPlan{
	// 		LaunchPlanKey: models.LaunchPlanKey{
	// 			Project: input.Project,
	// 			Domain:  input.Domain,
	// 			Name:    input.Name,
	// 			Version: input.Version,
	// 		},
	// 		Spec:       specBytes,
	// 		Closure:    closureBytes,
	// 		WorkflowID: 1,
	// 		State:      &state,
	// 	}, nil
	// }
	// repository.LaunchPlanRepo().(*repositoryMocks.MockLaunchPlanRepo).SetGetCallback(launchPlanGetFunc)
	// response, err := lpManager.GetLaunchPlan(context.Background(), admin.ObjectGetRequest{
	// 	Id: &launchPlanIdentifier,
	// })
	// assert.NoError(t, err)
	// assert.NotNil(t, response)
}


func TestNamedEntityMetadataManager_Update(t *testing.T) {
	// oldScheduleExpression := admin.Schedule{
	// 	ScheduleExpression: &admin.Schedule_Rate{
	// 		Rate: &admin.FixedRate{
	// 			Value: 2,
	// 			Unit:  admin.FixedRateUnit_HOUR,
	// 		},
	// 	},
	// }
	// oldLaunchPlanSpec := admin.LaunchPlanSpec{
	// 	EntityMetadata: &admin.LaunchPlanMetadata{
	// 		Schedule: &oldScheduleExpression,
	// 	},
	// }
	// oldLaunchPlanSpecBytes, _ := proto.Marshal(&oldLaunchPlanSpec)
	// newScheduleExpression := admin.Schedule{
	// 	ScheduleExpression: &admin.Schedule_CronExpression{
	// 		CronExpression: "cron",
	// 	},
	// }
	// newLaunchPlanSpec := admin.LaunchPlanSpec{
	// 	EntityMetadata: &admin.LaunchPlanMetadata{
	// 		Schedule: &newScheduleExpression,
	// 	},
	// }
	// newLaunchPlanSpecBytes, _ := proto.Marshal(&newLaunchPlanSpec)
	// mockScheduler := mocks.NewMockEventScheduler()
	// var removeCalled bool
	// mockScheduler.(*mocks.MockEventScheduler).SetRemoveScheduleFunc(
	// 	func(ctx context.Context, identifier admin.NamedEntityIdentifier) error {
	// 		assert.True(t, proto.Equal(&launchPlanNamedIdentifier, &identifier))
	// 		removeCalled = true
	// 		return nil
	// 	})
	// var addCalled bool
	// mockScheduler.(*mocks.MockEventScheduler).SetAddScheduleFunc(
	// 	func(ctx context.Context, input scheduleInterfaces.AddScheduleInput) error {
	// 		assert.True(t, proto.Equal(&launchPlanNamedIdentifier, &input.Identifier))
	// 		assert.True(t, proto.Equal(&newScheduleExpression, &input.ScheduleExpression))
	// 		addCalled = true
	// 		return nil
	// 	})
	// repository := getMockRepositoryForLpTest()
	// lpManager := NewLaunchPlanManager(repository, getMockConfigForLpTest(), mockScheduler, mockScope.NewTestScope())
	// err := lpManager.(*LaunchPlanManager).updateSchedules(
	// 	context.Background(),
	// 	models.LaunchPlan{
	// 		LaunchPlanKey: models.LaunchPlanKey{
	// 			Project: project,
	// 			Domain:  domain,
	// 			Name:    name,
	// 		},
	// 		Spec: newLaunchPlanSpecBytes,
	// 	},
	// 	&models.LaunchPlan{
	// 		LaunchPlanKey: models.LaunchPlanKey{
	// 			Project: project,
	// 			Domain:  domain,
	// 			Name:    name,
	// 		},
	// 		Spec: oldLaunchPlanSpecBytes,
	// 	})
	// assert.Nil(t, err)
	// assert.True(t, removeCalled)
	// assert.True(t, addCalled)
}

func TestNamedEntityMetadataManager_Update_BadRequest(t *testing.T) {}
