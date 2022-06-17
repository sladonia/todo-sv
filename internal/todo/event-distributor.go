package todo

import (
	"context"

	"github.com/sladonia/todo-sv/pkg/todopb"
	"go.uber.org/zap"
)

type UserEventsDistributor struct {
	workerGroupName      string
	listenSubjectPattern string
	subscriber           WorkerSubscriber
	publisher            Publisher
	log                  *zap.Logger
}

func NewUserEventsDistributor(
	workerGroupName string,
	listenSubjectPattern string,
	subscriber WorkerSubscriber,
	publisher Publisher,
	log *zap.Logger,
) *UserEventsDistributor {
	return &UserEventsDistributor{
		workerGroupName:      workerGroupName,
		listenSubjectPattern: listenSubjectPattern,
		subscriber:           subscriber,
		publisher:            publisher,
		log:                  log,
	}
}

func (l *UserEventsDistributor) Start(ctx context.Context) error {
	eventCh, err := l.subscriber.SubscribeGroup(ctx, l.listenSubjectPattern, l.workerGroupName)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-eventCh:
			if !ok {
				return nil
			}

			l.log.Debug("event distributor: new project event", zap.Any("event", event))

			project := event.GetProject()
			userIDs := project.ParticipantsIDs()

			for _, userID := range userIDs {
				err := l.publisher.Publish(ctx, todopb.NewUserEventsSubject(event.Type.String(), userID), event)
				if err != nil {
					l.log.Error("failed to publish event", zap.Error(err))
				}
			}
		}
	}
}
