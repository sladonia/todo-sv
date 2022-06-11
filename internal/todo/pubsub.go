package todo

import (
	"context"
	"sync"

	"github.com/sladonia/todo-sv/pkg/todopb"
)

type Publisher interface {
	Publish(ctx context.Context, project *todopb.Project) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, userID string) (<-chan *todopb.Project, error)
}

type PubSub interface {
	Publisher
	Subscriber
}

type inMemoryPubSub struct {
	sync.RWMutex
	subscribers map[string]chan *todopb.Project
}

func NewInMemoryPubSub() PubSub {
	return &inMemoryPubSub{
		subscribers: make(map[string]chan *todopb.Project),
	}
}

func (ps *inMemoryPubSub) Subscribe(ctx context.Context, userID string) (<-chan *todopb.Project, error) {
	ps.Lock()
	defer ps.Unlock()

	go func() {
		ps.Lock()
		defer ps.Unlock()

		<-ctx.Done()
		delete(ps.subscribers, userID)
	}()

	ch, ok := ps.subscribers[userID]
	if ok {
		return ch, nil
	}

	ch = make(chan *todopb.Project)
	ps.subscribers[userID] = ch

	return ch, nil
}

func (ps *inMemoryPubSub) Publish(_ context.Context, project *todopb.Project) error {
	ps.RLock()
	defer ps.RUnlock()

	for _, userID := range project.ParticipantsIDs() {
		ch, ok := ps.subscribers[userID]
		if !ok {
			continue
		}

		ch <- project
	}

	return nil
}
