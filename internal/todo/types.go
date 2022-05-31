package todo

import (
	"time"

	"github.com/rs/xid"
)

const (
	ProjectIDField           = "id"
	ProjectNameField         = "name"
	ProjectOwnerIDField      = "owner_id"
	ProjectParticipantsField = "participants"
	ProjectTasksField        = "tasks"

	UpdateProjectNameField         = "name"
	UpdateProjectOwnerIDField      = "owner_id"
	UpdateProjectParticipantsField = "participants"

	UpdateTaskTitleField       = "title"
	UpdateTaskDescriptionField = "description"
	UpdateTaskIsImportantField = "is_important"
	UpdateTaskIsFinishedField  = "is_finished"
	UpdateTaskTagsField        = "tags"
)

type Project struct {
	ID           string          `bson:"_id"`
	Version      string          `bson:"version"`
	Name         string          `bson:"name"`
	OwnerID      string          `bson:"owner_id"`
	Participants *Set            `bson:"participants"`
	Tasks        map[string]Task `bson:"tasks"`
	CreatedAt    time.Time       `bson:"created_at"`
	UpdatedAt    time.Time       `bson:"updated_at"`
}

type Task struct {
	ID          string    `bson:"id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Tags        *Set      `bson:"tags"`
	IsImportant bool      `bson:"is_important"`
	IsFinished  bool      `bson:"is_finished"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdateAt    time.Time `bson:"update_at"`
	Version     string    `bson:"version"`
}

func NewProject(name, ownerID string, participants []string) Project {
	now := time.Now()

	return Project{
		ID:           xid.New().String(),
		Version:      xid.New().String(),
		Name:         name,
		OwnerID:      ownerID,
		Participants: NewSet(participants...),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (p Project) CanEdit(userID string) bool {
	if p.OwnerID == userID {
		return true
	}

	if p.Participants.Contains(userID) {
		return true
	}

	return false
}

func (p Project) IsOwner(userID string) bool {
	if p.OwnerID == userID {
		return true
	}

	return false
}

func (p Project) WithTask(t Task) Project {
	p.Tasks[t.ID] = t
	p.Version = xid.New().String()
	p.UpdatedAt = time.Now()

	return p
}

func (p Project) WithoutTask(taskID string) Project {
	delete(p.Tasks, taskID)
	p.Version = xid.New().String()
	p.UpdatedAt = time.Now()

	return p
}

func (p Project) WithName(name string) Project {
	p.Name = name
	p.Version = xid.New().String()
	p.UpdatedAt = time.Now()

	return p
}

func (p Project) WithOwnerID(ownerID string) Project {
	p.OwnerID = ownerID
	p.Version = xid.New().String()
	p.UpdatedAt = time.Now()

	return p
}

func (p Project) WithParticipants(participantsIDs []string) Project {
	p.Participants = NewSet(participantsIDs...)
	p.Version = xid.New().String()
	p.UpdatedAt = time.Now()

	return p
}

func (t Task) WithTitle(title string) Task {
	t.Title = title
	t.Version = xid.New().String()
	t.UpdateAt = time.Now()

	return t
}

func (t Task) WithDescription(description string) Task {
	t.Title = description
	t.Version = xid.New().String()
	t.UpdateAt = time.Now()

	return t
}

func (t Task) WithTags(tags []string) Task {
	t.Tags = NewSet(tags...)
	t.Version = xid.New().String()
	t.UpdateAt = time.Now()

	return t
}

func (t Task) WithFinished(isFinished bool) Task {
	t.IsFinished = isFinished
	t.Version = xid.New().String()
	t.UpdateAt = time.Now()

	return t
}

func (t Task) WithImportant(isImportant bool) Task {
	t.IsFinished = isImportant
	t.Version = xid.New().String()
	t.UpdateAt = time.Now()

	return t
}
