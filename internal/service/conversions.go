package service

import (
	"time"

	"github.com/rs/xid"
	"github.com/sladonia/todo-sv/internal/todo"
	"github.com/sladonia/todo-sv/pkg/todopb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TaskToProtobuf(t todo.Task) *todopb.Task {
	return &todopb.Task{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Tags:        t.Tags.Values(),
		IsImportant: t.IsImportant,
		IsFinished:  t.IsFinished,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdateAt),
		Version:     t.Version,
	}
}

func ProjectToProtobuf(p todo.Project) *todopb.Project {
	var tasks map[string]*todopb.Task

	for taskID, task := range p.Tasks {
		pbTask := TaskToProtobuf(task)

		tasks[taskID] = pbTask
	}

	return &todopb.Project{
		Id:           p.ID,
		Name:         p.Name,
		OwnerId:      p.OwnerID,
		Participants: p.Participants.Values(),
		Tasks:        tasks,
		CreatedAt:    timestamppb.New(p.CreatedAt),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
		Version:      p.Version,
	}
}

func UpdateProject(p todo.Project, r *todopb.UpdateProjectRequest) todo.Project {
	fm := r.FieldMask
	fieldsSet := todo.NewSet(fm.Paths...)

	if fieldsSet.Contains(todo.UpdateProjectNameField) || len(fm.Paths) == 0 {
		p = p.WithName(r.Name)
	}
	if fieldsSet.Contains(todo.UpdateProjectOwnerIDField) || len(fm.Paths) == 0 {
		p = p.WithOwnerID(r.OwnerId)
	}
	if fieldsSet.Contains(todo.UpdateProjectParticipantsField) || len(fm.Paths) == 0 {
		p = p.WithParticipants(r.Participants)
	}

	return p
}

func UpdateTask(t todo.Task, r *todopb.UpdateTaskRequest) todo.Task {
	fm := r.FieldMask
	fieldSet := todo.NewSet(fm.Paths...)

	if fieldSet.Contains(todo.UpdateTaskTitleField) || fieldSet.IsEmpty() {
		t = t.WithTitle(r.Title)
	}
	if fieldSet.Contains(todo.UpdateTaskDescriptionField) || fieldSet.IsEmpty() {
		t = t.WithDescription(r.Description)
	}
	if fieldSet.Contains(todo.UpdateTaskIsImportantField) || fieldSet.IsEmpty() {
		t = t.WithImportant(r.IsImportant)
	}
	if fieldSet.Contains(todo.UpdateTaskIsFinishedField) || fieldSet.IsEmpty() {
		t = t.WithFinished(r.IsFinished)
	}
	if fieldSet.Contains(todo.UpdateTaskTagsField) || fieldSet.IsEmpty() {
		t = t.WithTags(r.Tags)
	}

	return t
}

func ProjectsToProtobuf(projects []todo.Project) *todopb.AllProjectsResponse {
	var resp *todopb.AllProjectsResponse

	for _, project := range projects {
		protobufProject := ProjectToProtobuf(project)
		resp.Projects = append(resp.Projects, protobufProject)
	}

	return resp
}

func NewTask(r *todopb.AddTaskRequest) todo.Task {
	now := time.Now()

	return todo.Task{
		ID:          xid.New().String(),
		Title:       r.Title,
		Description: r.Description,
		Tags:        todo.NewSet(r.Tags...),
		IsImportant: r.IsImportant,
		CreatedAt:   now,
		UpdateAt:    now,
		Version:     xid.New().String(),
	}
}
