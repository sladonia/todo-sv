package todopb

import (
	"time"

	"github.com/rs/xid"
	"github.com/sladonia/todo-sv/pkg/set"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	UpdateProjectNameField         = "name"
	UpdateProjectOwnerIDField      = "owner_id"
	UpdateProjectParticipantsField = "participants"
)

func NewProject(r *CreateProjectRequest) *Project {
	now := timestampNowMilliseconds()

	return &Project{
		Id:           xid.New().String(),
		Name:         r.Name,
		OwnerId:      r.OwnerId,
		Participants: unique(r.Participants),
		CreatedAt:    now,
		UpdatedAt:    now,
		Version:      xid.New().String(),
	}
}

func (x *Project) WithTask(task *Task) *Project {
	updated := *x

	updated.Tasks[task.Id] = task
	updated.Version = xid.New().String()
	updated.UpdatedAt = timestampNowMilliseconds()

	return &updated
}

func (x *Project) WithoutTask(taskID string) *Project {
	updated := *x

	delete(updated.Tasks, taskID)
	updated.Version = xid.New().String()
	updated.UpdatedAt = timestampNowMilliseconds()

	return &updated
}

func (x *Project) Update(r *UpdateProjectRequest) *Project {
	updated := *x

	fm := r.FieldMask
	fieldsSet := set.NewSet(fm.Paths...)

	if fieldsSet.Contains(UpdateProjectNameField) || len(fm.Paths) == 0 {
		updated.Name = r.Name
	}
	if fieldsSet.Contains(UpdateProjectOwnerIDField) || len(fm.Paths) == 0 {
		updated.OwnerId = r.OwnerId
	}
	if fieldsSet.Contains(UpdateProjectParticipantsField) || len(fm.Paths) == 0 {
		updated.Participants = unique(r.Participants)
	}

	updated.Version = xid.New().String()
	updated.UpdatedAt = timestampNowMilliseconds()

	return &updated
}

func (x *Project) CanEdit(userID string) bool {
	if x.OwnerId == userID {
		return true
	}

	for _, participant := range x.Participants {
		if participant == userID {
			return true
		}
	}

	return false
}

func (x *Project) IsOwner(userID string) bool {
	if x.OwnerId == userID {
		return true
	}

	return false
}

func unique(prev []string) []string {
	unique := make(map[string]struct{})

	for _, s := range prev {
		unique[s] = struct{}{}
	}

	updated := make([]string, len(unique))

	i := 0
	for s := range unique {
		updated[i] = s
		i++
	}

	return updated
}

func timestampNowMilliseconds() *timestamppb.Timestamp {
	timeNow := time.Now().Round(time.Millisecond)
	return timestamppb.New(timeNow)
}
