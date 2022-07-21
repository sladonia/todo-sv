package test

import (
	"context"
	"time"

	"github.com/sladonia/todo-sv/internal/todo"
	"github.com/sladonia/todo-sv/pkg/todopb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (s *Suite) TestGetProject() {
	ctx := context.Background()

	cases := []struct {
		name      string
		request   *todopb.GetProjectRequest
		expected  *todopb.Project
		errorCode codes.Code
	}{
		{
			name: "success",
			request: &todopb.GetProjectRequest{
				UserId:    "1",
				ProjectId: "2",
			},
			expected:  projectFixtureInserted1,
			errorCode: codes.OK,
		},
		{
			name: "not_found",
			request: &todopb.GetProjectRequest{
				UserId:    "1",
				ProjectId: "unexisting",
			},
			expected:  nil,
			errorCode: codes.NotFound,
		},
		{
			name: "unexisting_user",
			request: &todopb.GetProjectRequest{
				UserId:    "unexisting",
				ProjectId: "2",
			},
			expected:  nil,
			errorCode: codes.PermissionDenied,
		},
	}

	for _, testCase := range cases {
		s.Run(testCase.name, func() {
			res, err := s.service.GetProject(ctx, testCase.request)

			st, ok := status.FromError(err)
			if !ok {
				s.FailNow("not a grpc error")
			}

			s.Equal(testCase.expected, res)
			s.Equal(testCase.errorCode, st.Code())
		})
	}
}

func (s *Suite) TestCreateProject() {
	ownerID := "4"
	ctx := context.Background()

	cases := []struct {
		name      string
		request   *todopb.CreateProjectRequest
		expected  *todopb.Project
		errorCode codes.Code
	}{
		{
			name: "success",
			request: &todopb.CreateProjectRequest{
				Name:         "home stuff",
				OwnerId:      ownerID,
				Participants: []string{"5", "6"},
			},
			expected: &todopb.Project{
				Name:         "home stuff",
				OwnerId:      ownerID,
				Participants: []string{"5", "6"},
				Tasks:        nil,
			},
			errorCode: 0,
		},
		{
			name: "invalid_request",
			request: &todopb.CreateProjectRequest{
				Name:    "home stuff",
				OwnerId: "",
			},
			expected:  nil,
			errorCode: codes.InvalidArgument,
		},
	}

	for _, testCase := range cases {
		s.Run(testCase.name, func() {
			res, err := s.service.CreateProject(ctx, testCase.request)

			st, ok := status.FromError(err)
			if !ok {
				s.FailNow("not a grpc error")
			}

			if testCase.expected != nil {
				s.Equal(testCase.expected.Name, res.Name)
				s.Equal(testCase.expected.OwnerId, res.OwnerId)
				s.Equal(len(testCase.expected.Participants), len(res.Participants))
			}

			s.Equal(testCase.errorCode, st.Code())
		})
	}
}

func (s *Suite) TestUpdateProject() {
	ctx := context.Background()

	s.Run("success", func() {
		request := &todopb.UpdateProjectRequest{
			ProjectId: "2",
			UserId:    "1",
			Name:      "new_name",
			OwnerId:   "2",
			FieldMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "owner_id"},
			},
		}

		_, err := s.service.UpdateProject(ctx, request)
		s.NoError(err)

		updated, err := s.storage.ByID(ctx, request.ProjectId)

		s.NoError(err)
		s.Equal(request.Name, updated.Name)
		s.Equal(request.OwnerId, updated.OwnerId)
		s.True(updated.Version > projectFixtureInserted1.Version)
	})

	s.Run("invalid_request", func() {
		request := &todopb.UpdateProjectRequest{
			ProjectId: "2",
			UserId:    "",
			Name:      "new_name",
			OwnerId:   "2",
			FieldMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "owner_id"},
			},
		}

		_, err := s.service.UpdateProject(ctx, request)

		st, ok := status.FromError(err)
		if !ok {
			s.FailNow("not a grpc error")
		}

		s.Equal(codes.InvalidArgument, st.Code())
	})
}

func (s *Suite) TestAllProjects() {
	ctx := context.Background()

	testCases := []struct {
		name     string
		request  *todopb.AllProjectsRequest
		expected *todopb.AllProjectsResponse
		errCode  codes.Code
	}{
		{
			name: "success",
			request: &todopb.AllProjectsRequest{
				UserId: "2",
			},
			expected: &todopb.AllProjectsResponse{Projects: []*todopb.Project{
				projectFixtureInserted1,
				projectFixtureInserted2,
			}},
			errCode: 0,
		},
		{
			name: "no_projects_found",
			request: &todopb.AllProjectsRequest{
				UserId: "unexisting",
			},
			expected: &todopb.AllProjectsResponse{},
			errCode:  0,
		},
		{
			name: "invalid request",
			request: &todopb.AllProjectsRequest{
				UserId: "",
			},
			expected: nil,
			errCode:  codes.InvalidArgument,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			resp, err := s.service.AllProjects(ctx, testCase.request)

			st, ok := status.FromError(err)
			if !ok {
				s.FailNow("not a grpc error")
			}

			s.Equal(testCase.errCode, st.Code())

			if resp != nil {
				s.True(len(testCase.expected.Projects) == len(resp.Projects))
			}
		})
	}
}

func (s *Suite) TestAddTask() {
	ctx := context.Background()

	s.Run("success", func() {
		request := &todopb.AddTaskRequest{
			Title:       "remember the milk",
			ProjectId:   "2",
			UserId:      "2",
			Description: "remember to buy the milk",
			Tags:        []string{"groceries"},
			IsImportant: true,
		}

		_, err := s.service.AddTask(ctx, request)
		s.NoError(err)

		updated, err := s.storage.ByID(ctx, request.ProjectId)

		s.NoError(err)
		if !s.Len(updated.Tasks, 1) {
			s.FailNow("task not added")
		}
		addedTask := updated.TaskList()[0]
		s.Equal(request.Title, addedTask.Title)
		s.Equal(request.Description, addedTask.Description)
		s.Equal(request.Tags, addedTask.Tags)
		s.Equal(request.IsImportant, addedTask.IsImportant)
	})

	s.Run("permission_denied", func() {
		request := &todopb.AddTaskRequest{
			Title:       "remember the milk",
			ProjectId:   "2",
			UserId:      "5",
			Description: "remember to buy the milk",
			Tags:        []string{"groceries"},
			IsImportant: true,
		}

		_, err := s.service.AddTask(ctx, request)

		st, ok := status.FromError(err)
		if !ok {
			s.FailNow("not a grpc error")
		}

		s.Equal(codes.PermissionDenied, st.Code())
	})

	s.Run("invalid_request", func() {
		request := &todopb.AddTaskRequest{
			Title:       "remember the milk",
			ProjectId:   "",
			UserId:      "5",
			Description: "remember to buy the milk",
			Tags:        []string{"groceries"},
			IsImportant: true,
		}

		_, err := s.service.AddTask(ctx, request)

		st, ok := status.FromError(err)
		if !ok {
			s.FailNow("not a grpc error")
		}

		s.Equal(codes.InvalidArgument, st.Code())
	})
}

func (s *Suite) TestUpdateTask() {
	ctx := context.Background()

	testCases := []struct {
		name     string
		request  *todopb.UpdateTaskRequest
		expected *todopb.Task
		errCode  codes.Code
	}{
		{
			name: "success",
			request: &todopb.UpdateTaskRequest{
				TaskId:      "1",
				ProjectId:   "3",
				UserId:      "3",
				Title:       "pay",
				Description: "pay the god damn bill already!",
				Tags:        []string{"home", "bills"},
				IsImportant: true,
				IsFinished:  false,
				FieldMask:   nil, // all fields
			},
			expected: &todopb.Task{
				Title:       "pay",
				Description: "pay the god damn bill already!",
				Tags:        []string{"home", "bills"},
				IsImportant: true,
				IsFinished:  false,
			},
			errCode: 0,
		},
		{
			name: "task_not_found",
			request: &todopb.UpdateTaskRequest{
				TaskId:    "unexisting",
				ProjectId: "3",
				UserId:    "3",
				Title:     "buy",
				FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"title"}},
			},
			expected: nil,
			errCode:  codes.NotFound,
		},
		{
			name: "project_not_found",
			request: &todopb.UpdateTaskRequest{
				TaskId:    "1",
				ProjectId: "unexisting",
				UserId:    "3",
				Title:     "buy",
				FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"title"}},
			},
			expected: nil,
			errCode:  codes.NotFound,
		},
		{
			name: "permission_denied",
			request: &todopb.UpdateTaskRequest{
				TaskId:    "1",
				ProjectId: "3",
				UserId:    "unexisting",
				Title:     "buy",
				FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"title"}},
			},
			expected: nil,
			errCode:  codes.PermissionDenied,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			_, err := s.service.UpdateTask(ctx, testCase.request)

			st, ok := status.FromError(err)
			if !ok {
				s.FailNow("not a grpc error")
			}

			if !s.Equal(testCase.errCode, st.Code()) {
				s.FailNow("wrong status code")
			}

			if testCase.expected != nil {
				updated, err := s.storage.ByID(ctx, testCase.request.ProjectId)

				s.NoError(err)
				if !s.Len(updated.Tasks, 1) {
					s.FailNow("wrong tasks len")
				}

				task := updated.TaskList()[0]

				s.Equal(testCase.expected.Title, task.Title)
				s.Equal(testCase.expected.Description, task.Description)
				s.True(len(testCase.expected.Tags) == len(task.Tags))
				s.Equal(testCase.expected.IsImportant, task.IsImportant)
				s.Equal(testCase.expected.IsFinished, task.IsFinished)
			}
		})
	}
}

func (s *Suite) TestDeleteTask() {
	ctx := context.Background()

	s.Run("success", func() {
		request := &todopb.DeleteTaskRequest{
			ProjectId: "3",
			TaskId:    "1",
			UserId:    "3",
		}

		_, err := s.service.DeleteTask(ctx, request)
		s.NoError(err)

		updatedProject, err := s.storage.ByID(ctx, "3")
		s.NoError(err)
		s.Len(updatedProject.Tasks, 0)
	})

	s.Run("no_task_found", func() {
		request := &todopb.DeleteTaskRequest{
			ProjectId: "2",
			TaskId:    "1",
			UserId:    "1",
		}

		_, err := s.service.DeleteTask(ctx, request)
		s.NoError(err)
	})

	s.Run("permission_denied", func() {
		request := &todopb.DeleteTaskRequest{
			ProjectId: "3",
			TaskId:    "1",
			UserId:    "unexisting",
		}

		_, err := s.service.DeleteTask(ctx, request)

		st, ok := status.FromError(err)
		if !ok {
			s.FailNow("not a grpc error")
		}

		s.Equal(codes.PermissionDenied, st.Code())
	})
}

func (s *Suite) TestDeleteProject() {
	ctx := context.Background()

	testCases := []struct {
		name    string
		request *todopb.DeleteProjectRequest
		errCode codes.Code
	}{
		{
			name: "success",
			request: &todopb.DeleteProjectRequest{
				ProjectId: "2",
				UserId:    "1",
			},
			errCode: 0,
		},
		{
			name: "not_found",
			request: &todopb.DeleteProjectRequest{
				ProjectId: "unexisting",
				UserId:    "1",
			},
			errCode: 0,
		},
		{
			name: "permission_denied",
			request: &todopb.DeleteProjectRequest{
				ProjectId: "3",
				UserId:    "3",
			},
			errCode: codes.PermissionDenied,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			_, err := s.service.DeleteProject(ctx, testCase.request)

			st, ok := status.FromError(err)
			if !ok {
				s.FailNow("not a grpc error")
			}

			s.Equal(testCase.errCode, st.Code())

			if err == nil {
				_, err := s.storage.ByID(ctx, testCase.request.ProjectId)
				s.ErrorIs(err, todo.ErrProjectNotFound)
			}
		})
	}
}

func (s *Suite) TestSubscribeToProjectsUpdates() {
	ctx := context.Background()

	r := &todopb.ProjectsUpdatesRequest{
		UserId:   projectFixtureInserted1.OwnerId,
		DeviceId: "1",
	}

	subscribeServer := newMockSubscribeServer()

	go func() {
		err := s.service.SubscribeToProjectsUpdates(r, subscribeServer)
		s.Require().NoError(err)
	}()

	time.Sleep(50 * time.Millisecond)

	addTask := &todopb.AddTaskRequest{
		Title:       "to buy something",
		ProjectId:   "2",
		UserId:      "1",
		Description: "just buy",
		Tags:        []string{"groceries"},
		IsImportant: true,
	}

	_, err := s.service.AddTask(ctx, addTask)
	s.NoError(err)

	ev := <-subscribeServer.eventCh
	s.Equal(todopb.EventType_PROJECT_UPDATED, ev.Type)
	s.Len(ev.Project.Tasks, 1)

	_, err = s.service.DeleteProject(ctx, &todopb.DeleteProjectRequest{
		ProjectId: "2",
		UserId:    "1",
	})

	ev = <-subscribeServer.eventCh
	s.Equal(todopb.EventType_PROJECT_DELETED, ev.Type)
}
