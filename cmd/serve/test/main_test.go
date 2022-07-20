package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/sladonia/common-lb/dockertool"
	"github.com/sladonia/common-lb/dockertool/dockermongo"
	"github.com/sladonia/common-lb/logger"
	"github.com/sladonia/todo-sv/internal/mongodb"
	"github.com/sladonia/todo-sv/internal/todo"
	"github.com/sladonia/todo-sv/pkg/todopb"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	projectDBName          = "todo_test"
	projectsCollectionName = "projects_test"
)

var projectFixtureInserted1 = &todopb.Project{
	Id:           "2",
	Name:         "to-buy",
	OwnerId:      "1",
	Participants: []string{"2", "3"},
	Tasks:        map[string]*todopb.Task{},
	CreatedAt:    timestamppb.New(time.Now().Round(time.Millisecond)),
	UpdatedAt:    timestamppb.New(time.Now().Round(time.Millisecond)),
	Version:      "cag6rg19d3prkkb9fuag",
}

var projectFixtureInserted2 = &todopb.Project{
	Id:           "3",
	Name:         "different",
	OwnerId:      "2",
	Participants: []string{"3"},
	Tasks: map[string]*todopb.Task{
		"1": {
			Id:          "1",
			Title:       "pay bill",
			Description: "pay the god damn bill already!",
			Tags:        []string{"home", "bills"},
			IsImportant: true,
			IsFinished:  false,
			CreatedAt:   timestamppb.New(time.Now().Round(time.Millisecond)),
			UpdatedAt:   timestamppb.New(time.Now().Round(time.Millisecond)),
			Version:     "cai6enp9d3pjf0mq7se0",
		},
	},
	CreatedAt: timestamppb.New(time.Now().Round(time.Millisecond)),
	UpdatedAt: timestamppb.New(time.Now().Round(time.Millisecond)),
	Version:   "cag7sf19d3pr4a5t1bn0",
}

type Suite struct {
	suite.Suite

	log            *zap.Logger
	dockerPool     *dockertest.Pool
	db             *mongo.Database
	mongoContainer dockertool.Container
	natsContainer  dockertool.Container
	mongoDSN       string
	storage        todo.Storage
	pubSub         todo.PubSub
	service        todopb.ToDoServiceServer
}

func (s *Suite) SetupSuite() {
	var err error

	s.log, err = logger.NewZap("debug")
	if err != nil {
		panic(fmt.Sprintf("init logger: %s", err.Error()))
	}

	s.dockerPool, err = dockertest.NewPool("")
	if err != nil {
		s.log.Panic("init docker pool", zap.Error(err))
	}

	s.mongoContainer = dockermongo.NewDefault(s.log)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err = s.mongoContainer.Start(ctx, s.dockerPool)
	if err != nil {
		s.log.Panic("start mongo container", zap.Error(err))
	}

	s.mongoDSN = s.mongoContainer.DSN()

	err = s.mongoContainer.WaitReady(ctx)
	if err != nil {
		s.log.Panic("mongo container start timeout", zap.Error(err))
	}

	s.db, err = mongodb.Connect(ctx, s.mongoDSN, projectDBName)
	if err != nil {
		s.log.Panic("failed to connect mongo", zap.Error(err))
	}

	s.storage = todo.NewStorage(s.db, projectsCollectionName)
	s.pubSub = todo.NewNopPubSub()
	s.service = todo.NewService(s.log, s.storage, s.pubSub)
}

func (s *Suite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := s.mongoContainer.Stop(ctx)
	if err != nil {
		s.log.Panic("stop container", zap.Error(err))
	}
}

func (s *Suite) SetupTest() {
	err := s.storage.Insert(context.Background(), projectFixtureInserted1)
	if err != nil {
		s.log.Panic("failed to insert fixture", zap.Error(err))
	}

	err = s.storage.Insert(context.Background(), projectFixtureInserted2)
	if err != nil {
		s.log.Panic("failed to insert fixture", zap.Error(err))
	}
}

func (s *Suite) TearDownTest() {
	_, err := s.db.Collection(projectsCollectionName).DeleteMany(context.Background(), bson.M{})
	if err != nil {
		s.log.Panic("failed to delete projects", zap.Error(err))
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
