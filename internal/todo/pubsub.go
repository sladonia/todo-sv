package todo

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	"github.com/sladonia/todo-sv/pkg/todopb"
)

type Publisher interface {
	Publish(ctx context.Context, subject string, ev *todopb.Event) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, subject string) (<-chan *todopb.Event, error)
}

type WorkerSubscriber interface {
	SubscribeGroup(ctx context.Context, subject, groupName string) (<-chan *todopb.Event, error)
}

type PubSub interface {
	Publisher
	Subscriber
	WorkerSubscriber
}

type nopPubSub struct{}

func (n *nopPubSub) Publish(ctx context.Context, subject string, ev *todopb.Event) error {
	return nil
}

func (n *nopPubSub) Subscribe(ctx context.Context, subject string) (<-chan *todopb.Event, error) {
	return nil, nil
}

func (n *nopPubSub) SubscribeGroup(ctx context.Context, subject, groupName string) (<-chan *todopb.Event, error) {
	return nil, nil
}

func NewNopPubSub() PubSub {
	return &nopPubSub{}
}

type NatsPubSub struct {
	conn *nats.EncodedConn
}

func NewNatsPubSub(dsn string) (PubSub, error) {
	nc, err := nats.Connect(dsn)
	if err != nil {
		return nil, err
	}

	ec, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return nil, err
	}

	return &NatsPubSub{conn: ec}, nil
}

func (ps *NatsPubSub) Publish(_ context.Context, subject string, ev *todopb.Event) error {
	return ps.conn.Publish(subject, ev)
}

func (ps *NatsPubSub) Subscribe(ctx context.Context, subject string) (<-chan *todopb.Event, error) {
	eventCh := make(chan *todopb.Event)

	_, err := ps.conn.BindRecvChan(subject, eventCh)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		close(eventCh)
	}()

	return eventCh, nil
}

func (ps *NatsPubSub) SubscribeGroup(ctx context.Context, subject, groupName string) (<-chan *todopb.Event, error) {
	eventCh := make(chan *todopb.Event)

	_, err := ps.conn.BindRecvQueueChan(subject, groupName, eventCh)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		close(eventCh)
	}()

	return eventCh, nil
}
