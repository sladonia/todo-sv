package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sladonia/todo-sv/internal/todo"
	"github.com/sladonia/todo-sv/pkg/todopb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	todopb.UnimplementedToDoServiceServer
	storage todo.Storage
	log     *zap.Logger
}

func NewService(log *zap.Logger, storage todo.Storage) todopb.ToDoServiceServer {
	return &service{
		storage: storage,
		log:     log,
	}
}

func (s *service) CreateProject(ctx context.Context, r *todopb.CreateProjectRequest) (*todopb.Project, error) {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("create project invalid request", zap.String("error", err.Error()))
		return nil, s.wrapError(err)
	}

	project := todo.NewProject(r.Name, r.OwnerId, r.Participants)

	err = s.storage.Insert(ctx, project)
	if err != nil {
		s.log.Debug("failed to insert project to db", zap.String("error", err.Error()))
		return nil, s.wrapError(err)
	}

	return ProjectToProtobuf(project), nil
}

func (s *service) GetProject(ctx context.Context, r *todopb.GetProjectRequest) (*todopb.Project, error) {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("get project invalid request", zap.String("error", err.Error()))
		return nil, s.wrapError(err)
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, todo.ErrProjectNotFound) {
			s.log.Info("failed to retrieve project", zap.Error(err))
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

	// TODO: add field masks

	return ProjectToProtobuf(p), nil
}
func (s *service) UpdateProject(ctx context.Context, r *todopb.UpdateProjectRequest) (*emptypb.Empty, error) {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("get project invalid request", zap.String("error", err.Error()))
		return empty(), s.wrapError(err)
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, todo.ErrProjectNotFound) {
			s.log.Info("failed to retrieve project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	if p.OwnerID != r.UserId {
		s.log.Debug(
			"update project. permission denied",
			zap.String("user_id", r.UserId),
			zap.String("project_id", r.ProjectId),
		)
		return empty(), status.Error(codes.PermissionDenied, "user has not modify access wrights to the project")
	}

	updated := UpdateProject(p, r)

	err = s.storage.Replace(ctx, p, updated)
	if err != nil {
		s.log.Debug("failed to replace project", zap.String("error", err.Error()))
		return nil, s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) AllProjects(ctx context.Context, r *todopb.AllProjectsRequest) (*todopb.AllProjectsResponse, error) {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("get all projects invalid request", zap.String("error", err.Error()))
		return nil, s.wrapError(err)
	}

	projects, err := s.storage.AllUserProjects(ctx, r.UserId)
	if err != nil {
		s.log.Info("failed to retrieve projects from storage", zap.Error(err))
		return nil, s.wrapError(err)
	}

	return ProjectsToProtobuf(projects), nil
}

func (s *service) AddTask(ctx context.Context, r *todopb.AddTaskRequest) (*emptypb.Empty, error) {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("add task invalid request", zap.String("error", err.Error()))
		return empty(), s.wrapError(err)
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, todo.ErrProjectNotFound) {
			s.log.Info("failed to retrieve project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	if !p.CanEdit(r.UserId) {
		return nil, status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("user %s has no access to %s project", r.UserId, r.ProjectId),
		)
	}

	task := NewTask(r)
	updatedProject := p.WithTask(task)

	err = s.storage.Replace(ctx, p, updatedProject)
	if err != nil {
		s.log.Info("failed to update project", zap.Error(err))
		return empty(), s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) UpdateTask(ctx context.Context, r *todopb.UpdateTaskRequest) (*emptypb.Empty, error) {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("update task invalid request", zap.String("error", err.Error()))
		return empty(), s.wrapError(err)
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, todo.ErrProjectNotFound) {
			s.log.Info("failed to retrieve project", zap.Error(err))
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

	updatedTask := UpdateTask(task, r)
	updatedProject := p.WithTask(updatedTask)

	err = s.storage.Replace(ctx, p, updatedProject)
	if err != nil {
		s.log.Info("failed to update project", zap.Error(err))
		return empty(), s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) DeleteTask(ctx context.Context, r *todopb.DeleteTaskRequest) (*emptypb.Empty, error) {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("delete task invalid request", zap.String("error", err.Error()))
		return empty(), s.wrapError(err)
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, todo.ErrProjectNotFound) {
			s.log.Info("failed to retrieve project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	if !p.CanEdit(r.UserId) {
		return nil, status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("user %s has no access to %s project", r.UserId, r.ProjectId),
		)
	}

	updatedProject := p.WithoutTask(r.TaskId)

	err = s.storage.Replace(ctx, p, updatedProject)
	if err != nil {
		s.log.Info("failed to update project", zap.Error(err))
		return empty(), s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) DeleteProject(ctx context.Context, r *todopb.DeleteProjectRequest) (*emptypb.Empty, error) {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("delete project invalid request", zap.String("error", err.Error()))
		return empty(), s.wrapError(err)
	}

	p, err := s.storage.ByID(ctx, r.ProjectId)
	if err != nil {
		if !errors.Is(err, todo.ErrProjectNotFound) {
			s.log.Info("failed to retrieve project", zap.Error(err))
		}

		return empty(), s.wrapError(err)
	}

	if p.OwnerID != r.UserId {
		return empty(), status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("user %s has no access to delete %s project", r.UserId, r.ProjectId),
		)
	}

	err = s.storage.Delete(ctx, r.ProjectId)
	if err != nil {
		s.log.Info("failed to delete project", zap.Error(err))
		return empty(), s.wrapError(err)
	}

	return empty(), nil
}

func (s *service) SubscribeToProjectsUpdates(
	r *todopb.ProjectsUpdatesRequest,
	updateServer todopb.ToDoService_SubscribeToProjectsUpdatesServer,
) error {
	err := r.ValidateAll()
	if err != nil {
		s.log.Debug("subscribe to project updates invalid request", zap.String("error", err.Error()))
		return s.wrapError(err)
	}

	_, err = s.storage.AllUserProjects(updateServer.Context(), r.UserId)
	if err != nil {
		s.log.Info("failed retrieve user projects", zap.Error(err))
		return s.wrapError(err)
	}

	// TODO: implement subscription

	for {
		select {
		case <-updateServer.Context().Done():
			s.log.Debug("subscribe to project context done", zap.String("error", err.Error()))
			return nil
		}
	}

	return nil
}

func (s *service) wrapError(err error) error {
	if err == nil {
		return nil
	}

	var validationError todopb.CreateProjectRequestMultiError

	switch {
	case errors.As(err, &validationError):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, todo.ErrProjectNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, todo.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, todo.ErrIDsMismatch):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, todo.ErrVersionMismatch):
		return status.Error(codes.Aborted, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}

func empty() *emptypb.Empty {
	return &emptypb.Empty{}
}
