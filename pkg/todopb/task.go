package todopb

import (
	"github.com/rs/xid"
	"github.com/sladonia/todo-sv/pkg/set"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	UpdateTaskTitleField       = "title"
	UpdateTaskDescriptionField = "description"
	UpdateTaskIsImportantField = "is_important"
	UpdateTaskIsFinishedField  = "is_finished"
	UpdateTaskTagsField        = "tags"
)

func NewTask(r *AddTaskRequest) *Task {
	now := timestamppb.Now()

	return &Task{
		Id:          xid.New().String(),
		Title:       r.Title,
		Description: r.Description,
		Tags:        unique(r.Tags),
		IsImportant: r.IsImportant,
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     xid.New().String(),
	}
}

func (x *Task) UpdateTask(r *UpdateTaskRequest) *Task {
	updated := x.clone()

	fm := r.FieldMask

	var fieldSet *set.Set

	if fm != nil {
		fieldSet = set.NewSet(fm.Paths...)
	} else {
		fieldSet = set.NewSet()
	}

	if fieldSet.Contains(UpdateTaskTitleField) || fieldSet.IsEmpty() {
		updated.Title = r.Title
	}
	if fieldSet.Contains(UpdateTaskDescriptionField) || fieldSet.IsEmpty() {
		updated.Description = r.Description
	}
	if fieldSet.Contains(UpdateTaskIsImportantField) || fieldSet.IsEmpty() {
		updated.IsImportant = r.IsImportant
	}
	if fieldSet.Contains(UpdateTaskIsFinishedField) || fieldSet.IsEmpty() {
		updated.IsFinished = r.IsFinished
	}
	if fieldSet.Contains(UpdateTaskTagsField) || fieldSet.IsEmpty() {
		updated.Tags = unique(r.Tags)
	}

	x.Version = xid.New().String()
	x.UpdatedAt = timestamppb.Now()

	return updated
}

func (x *Task) clone() *Task {
	tags := make([]string, len(x.Tags))
	copy(tags, x.Tags)

	createdAt := *x.CreatedAt
	updatedAt := *x.UpdatedAt

	return &Task{
		Id:          x.Id,
		Title:       x.Title,
		Description: x.Description,
		Tags:        tags,
		IsImportant: x.IsImportant,
		IsFinished:  x.IsFinished,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
		Version:     x.Version,
	}
}
