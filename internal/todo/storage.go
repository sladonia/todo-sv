package todo

import (
	"context"
	"time"

	"github.com/sladonia/todo-sv/pkg/todopb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Storage interface {
	ByID(ctx context.Context, projectID string) (*todopb.Project, error)
	AllUserProjects(ctx context.Context, userID string) ([]*todopb.Project, error)
	Insert(ctx context.Context, project *todopb.Project) error
	Replace(ctx context.Context, prev, curr *todopb.Project) error
	Delete(ctx context.Context, projectID string) error
}

type ProjectBSON struct {
	ID           string              `bson:"_id"`
	Name         string              `bson:"name"`
	OwnerID      string              `bson:"owner_id"`
	Participants []string            `bson:"participants"`
	Tasks        map[string]TaskBSON `bson:"tasks"`
	CreatedAt    time.Time           `bson:"created_at"`
	UpdatedAt    time.Time           `bson:"updated_at"`
	Version      string              `bson:"version"`
}

type TaskBSON struct {
	ID          string    `bson:"id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Tags        []string  `bson:"tags"`
	IsImportant bool      `bson:"is_important"`
	IsFinished  bool      `bson:"is_finished"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	Version     string    `bson:"version"`
}

func NewProjectBSON(p *todopb.Project) ProjectBSON {
	tasks := make(map[string]TaskBSON)

	for id, task := range p.Tasks {
		taskBSON := NewTaskBSON(task)
		tasks[id] = taskBSON
	}

	return ProjectBSON{
		ID:           p.Id,
		Name:         p.Name,
		OwnerID:      p.OwnerId,
		Participants: p.Participants,
		Tasks:        tasks,
		CreatedAt:    p.CreatedAt.AsTime(),
		UpdatedAt:    p.UpdatedAt.AsTime(),
		Version:      p.Version,
	}
}

func NewTaskBSON(t *todopb.Task) TaskBSON {
	return TaskBSON{
		ID:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		Tags:        t.Tags,
		IsImportant: t.IsImportant,
		IsFinished:  t.IsFinished,
		CreatedAt:   t.CreatedAt.AsTime(),
		UpdatedAt:   t.UpdatedAt.AsTime(),
		Version:     t.Version,
	}
}

func (p *ProjectBSON) Project() *todopb.Project {
	tasks := make(map[string]*todopb.Task)

	for id, taskBSON := range p.Tasks {
		task := taskBSON.Task()
		tasks[id] = task
	}

	return &todopb.Project{
		Id:           p.ID,
		Name:         p.Name,
		OwnerId:      p.OwnerID,
		Participants: p.Participants,
		Tasks:        tasks,
		CreatedAt:    timestamppb.New(p.CreatedAt),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
		Version:      p.Version,
	}
}

func (t *TaskBSON) Task() *todopb.Task {
	return &todopb.Task{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Tags:        t.Tags,
		IsImportant: t.IsImportant,
		IsFinished:  t.IsFinished,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
		Version:     t.Version,
	}
}
