package todopb

import (
	"github.com/rs/xid"
)

func NewProjectCreatedEvent(p *Project) *Event {
	return &Event{
		Id:        xid.New().String(),
		Type:      EventType_PROJECT_CREATED,
		Project:   p,
		CreatedAt: timestampNowMilliseconds(),
	}
}

func NewProjectUpdatedEvent(p *Project) *Event {
	return &Event{
		Id:        xid.New().String(),
		Type:      EventType_PROJECT_UPDATED,
		Project:   p,
		CreatedAt: timestampNowMilliseconds(),
	}
}

func NewProjectDeletedEvent(p *Project) *Event {
	return &Event{
		Id:        xid.New().String(),
		Type:      EventType_PROJECT_DELETED,
		Project:   p,
		CreatedAt: timestampNowMilliseconds(),
	}
}
