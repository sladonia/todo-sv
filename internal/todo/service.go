package todo

import (
	"context"
	"errors"
	"fmt"

	"github.com/sladonia/todo-sv/pkg/todopb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	todopb.UnimplementedToDoServiceServer
	storage Storage
	pubSub  PubSub
	log     *zap.Logger
}

func NewService(log *zap.Logger, storage Storage, pubSub PubSub) todopb.ToDoServiceServer {
	return &service{
		storage: storage,
		log:     log,
		pubSub:  pubSub,
	}
}

func (s *service) CreateProject(ctx context.Context, r *todopb.CreateProjectRequest) (*todopb.Project, error) {
	s.log.Debug("create project request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("create project invalid request", zap.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	project := todopb.NewProject(r)

	err = s.storage.Insert(ctx, project)
	if err != nil {
		if !errors.Is(err, ErrAlreadyExists) {
			s.log.Error("failed to insert project to db", zap.String("error", err.Error()))
		}

		return nil, s.wrapError(err)
	}

	err = s.pubSub.Publish(ctx, project)
	if err != nil {
		s.log.Error("publish project", zap.Error(err))
		return nil, s.wrapError(err)
	}

	return project, nil
}

func (s *service) GetProject(ctx context.Context, r *todopb.GetProjectRequest) (*todopb.Project, error) {
	s.log.Debug("get project request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("get project invalid request", zap.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, ErrProjectNotFound) {
			s.log.Error("failed to retrieve project", zap.Error(err))
		}

		return nil, s.wrapError(err)
	}

	if !p.CanEdit(r.UserId) {
		s.log.Debug(
			"get project. permission denied",
			zap.String("user_id", r.UserId),
			zap.String("project_id", r.ProjectId),
		)
		return nil, status.Error(codes.PermissionDenied, "user has not get access wrights to the project")
	}

	return p, nil
}

func (s *service) UpdateProject(ctx context.Context, r *todopb.UpdateProjectRequest) (*emptypb.Empty, error) {
	s.log.Debug("update project request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("update project invalid request", zap.String("error", err.Error()))
		return empty(), status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, ErrProjectNotFound) {
			s.log.Error("failed to retrieve project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	if !p.IsOwner(r.UserId) {
		s.log.Debug(
			"update project. permission denied",
			zap.String("user_id", r.UserId),
			zap.String("project_id", r.ProjectId),
		)
		return empty(), status.Error(codes.PermissionDenied, "user has not modify access wrights to the project")
	}

	updated := p.Update(r)

	err = s.storage.Replace(ctx, p, updated)
	if err != nil {
		if !IsStorageError(err) {
			s.log.Error("failed to replace project", zap.String("error", err.Error()))
		}

		return nil, s.wrapError(err)
	}

	err = s.pubSub.Publish(ctx, updated)
	if err != nil {
		s.log.Error("publish project", zap.Error(err))
		return nil, s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) AllProjects(ctx context.Context, r *todopb.AllProjectsRequest) (*todopb.AllProjectsResponse, error) {
	s.log.Debug("all projects request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("get all projects invalid request", zap.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	projects, err := s.storage.AllUserProjects(ctx, r.UserId)
	if err != nil {
		s.log.Error("failed to retrieve projects from storage", zap.Error(err))
		return nil, s.wrapError(err)
	}

	return &todopb.AllProjectsResponse{Projects: projects}, nil
}

func (s *service) AddTask(ctx context.Context, r *todopb.AddTaskRequest) (*emptypb.Empty, error) {
	s.log.Debug("add task request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("add task invalid request", zap.String("error", err.Error()))
		return empty(), status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, ErrProjectNotFound) {
			s.log.Error("failed to retrieve project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	if !p.CanEdit(r.UserId) {
		return empty(), status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("user %s has no access to %s project", r.UserId, r.ProjectId),
		)
	}

	task := todopb.NewTask(r)
	updatedProject := p.WithTask(task)

	err = s.storage.Replace(ctx, p, updatedProject)
	if err != nil {
		if !IsStorageError(err) {
			s.log.Error("failed to update project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	err = s.pubSub.Publish(ctx, updatedProject)
	if err != nil {
		s.log.Error("publish project", zap.Error(err))
		return nil, s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) UpdateTask(ctx context.Context, r *todopb.UpdateTaskRequest) (*emptypb.Empty, error) {
	s.log.Debug("update task request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("update task invalid request", zap.String("error", err.Error()))
		return empty(), status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, ErrProjectNotFound) {
			s.log.Error("failed to retrieve project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	if !p.CanEdit(r.UserId) {
		return empty(), status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("user %s has no access to %s project", r.UserId, r.ProjectId),
		)
	}

	task, ok := p.Tasks[r.TaskId]
	if !ok {
		return empty(), status.Error(
			codes.NotFound,
			fmt.Sprintf("task_id=%s not found in project_id=%s", r.TaskId, r.ProjectId),
		)
	}

	updatedTask := task.UpdateTask(r)
	updatedProject := p.WithTask(updatedTask)

	err = s.storage.Replace(ctx, p, updatedProject)
	if err != nil {
		if !IsStorageError(err) {
			s.log.Error("failed to update project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	err = s.pubSub.Publish(ctx, updatedProject)
	if err != nil {
		s.log.Error("publish project", zap.Error(err))
		return nil, s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) DeleteTask(ctx context.Context, r *todopb.DeleteTaskRequest) (*emptypb.Empty, error) {
	s.log.Debug("delete task request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("delete task invalid request", zap.String("error", err.Error()))
		return empty(), status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, ErrProjectNotFound) {
			s.log.Error("failed to retrieve project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	if !p.CanEdit(r.UserId) {
		return empty(), status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("user %s has no access to %s project", r.UserId, r.ProjectId),
		)
	}

	updatedProject := p.WithoutTask(r.TaskId)

	err = s.storage.Replace(ctx, p, updatedProject)
	if err != nil {
		if !IsStorageError(err) {
			s.log.Error("failed to update project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	err = s.pubSub.Publish(ctx, updatedProject)
	if err != nil {
		s.log.Error("publish project", zap.Error(err))
		return nil, s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) DeleteProject(ctx context.Context, r *todopb.DeleteProjectRequest) (*emptypb.Empty, error) {
	s.log.Debug("delete project request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("delete project invalid request", zap.String("error", err.Error()))
		return empty(), status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, ErrProjectNotFound) {
			s.log.Error("failed to retrieve project", zap.Error(err))
		}

		return empty(), nil
	}

	if !p.IsOwner(r.UserId) {
		return empty(), status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("user %s has no access to delete %s project", r.UserId, r.ProjectId),
		)
	}

	err = s.storage.Delete(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, ErrProjectNotFound) {
			s.log.Error("failed to delete project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) SubscribeToProjectsUpdates(
	r *todopb.ProjectsUpdatesRequest,
	updateServer todopb.ToDoService_SubscribeToProjectsUpdatesServer,
) error {
	s.log.Debug("subscribe to project updates request", zap.Any("request_body", r))

	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("subscribe to project updates invalid request", zap.String("error", err.Error()))
		return status.Error(codes.InvalidArgument, err.Error())
	}

	projectsCh, err := s.pubSub.Subscribe(updateServer.Context(), fmt.Sprintf("%s:%s", r.UserId, r.DeviceId))
	if err != nil {
		s.log.Error("failed to subscribe", zap.Error(err))
		return s.wrapError(err)
	}

Loop:
	for {
		select {
		case <-updateServer.Context().Done():
			break Loop
		case p, ok := <-projectsCh:
			if !ok {
				break Loop
			}

			err := updateServer.Send(p)
			if err != nil {
				s.log.Error("failed to send project", zap.Error(err))
			}
		}
	}

	s.log.Debug("subscription canceled", zap.String("error", err.Error()))
	return nil
}

func (s *service) wrapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrProjectNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, ErrIDsMismatch):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrVersionMismatch):
		return status.Error(codes.Aborted, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}

func empty() *emptypb.Empty {
	return &emptypb.Empty{}
}
