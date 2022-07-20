package todo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/sladonia/common-lb/logger"
	"github.com/sladonia/dockert"
	"github.com/sladonia/dockert/container"
	"github.com/sladonia/todo-sv/internal/mongodb"
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

var projectFixture1 = &todopb.Project{
	Id:           "1",
	Name:         "personal",
	OwnerId:      "1",
	Participants: []string{"2", "3"},
	Tasks:        map[string]*todopb.Task{},
	CreatedAt:    timestamppb.New(time.Now().Round(time.Millisecond)),
	UpdatedAt:    timestamppb.New(time.Now().Round(time.Millisecond)),
	Version:      "cafmaj19d3pichevllcg",
}

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
	Tasks:        map[string]*todopb.Task{},
	CreatedAt:    timestamppb.New(time.Now().Round(time.Millisecond)),
	UpdatedAt:    timestamppb.New(time.Now().Round(time.Millisecond)),
	Version:      "cag7sf19d3pr4a5t1bn0",
}

type Suite struct {
	suite.Suite

	log            *zap.Logger
	dockerPool     *dockertest.Pool
	db             *mongo.Database
	mongoContainer dockert.Container
	mongoDSN       string
	storage        Storage
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

	s.mongoContainer = container.NewMongo()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err = s.mongoContainer.Start(ctx, s.dockerPool)
	if err != nil {
		s.log.Panic("start mongo container", zap.Error(err))
	}

	s.mongoDSN = container.MongoDSN(s.mongoContainer)

	err = s.mongoContainer.WaitReady(ctx)
	if err != nil {
		s.log.Panic("mongo container start timeout", zap.Error(err))
	}

	s.db, err = mongodb.Connect(ctx, s.mongoDSN, projectDBName)
	if err != nil {
		s.log.Panic("failed to connect mongo", zap.Error(err))
	}

	s.storage = NewStorage(s.db, projectsCollectionName)
}

func (s *Suite) TearDownSuite() {
	err := s.mongoContainer.Stop()
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

func (s *Suite) TestByID() {
	ctx := context.Background()

	s.Run("success", func() {
		retrieved, err := s.storage.ByID(ctx, "2")

		s.NoError(err)
		s.Equal(projectFixtureInserted1, retrieved)
	})

	s.Run("not_found", func() {
		_, err := s.storage.ByID(ctx, "unexisting")
		s.ErrorIs(err, ErrProjectNotFound)
	})
}

func (s *Suite) TestInsertD() {
	ctx := context.Background()

	s.Run("success", func() {
		err := s.storage.Insert(ctx, projectFixture1)
		s.NoError(err)

		retrieved, err := s.storage.ByID(ctx, projectFixture1.Id)
		s.NoError(err)

		s.Equal(projectFixture1, retrieved)
	})

	s.Run("already_exists", func() {
		err := s.storage.Insert(ctx, projectFixture1)
		s.ErrorIs(err, ErrAlreadyExists)
	})
}

func (s *Suite) TestAllUserProjects() {
	ctx := context.Background()

	s.Run("success", func() {
		projects1, err := s.storage.AllUserProjects(ctx, "1")
		s.NoError(err)
		s.Len(projects1, 1)

		projects2, err := s.storage.AllUserProjects(ctx, "2")
		s.NoError(err)
		s.Len(projects2, 2)

		projects3, err := s.storage.AllUserProjects(ctx, "3")
		s.NoError(err)
		s.Len(projects3, 2)
	})

	s.Run("not_found", func() {
		projects, err := s.storage.AllUserProjects(ctx, "uexisting")

		s.NoError(err)
		s.Empty(projects)
	})
}

func (s *Suite) TestDelete() {
	ctx := context.Background()

	s.Run("success", func() {
		err := s.storage.Delete(ctx, projectFixtureInserted1.Id)
		s.NoError(err)

		_, err = s.storage.ByID(ctx, projectFixtureInserted1.Id)
		s.ErrorIs(err, ErrProjectNotFound)
	})

	s.Run("not_found", func() {
		err := s.storage.Delete(ctx, "unexisting")
		s.NoError(err)
	})
}

func (s *Suite) TestReplace() {
	ctx := context.Background()

	s.Run("success", func() {
		task := todopb.NewTask(&todopb.AddTaskRequest{
			Title:       "to do exercises",
			ProjectId:   "2",
			UserId:      "1",
			Description: "",
			Tags:        []string{"sport", "fun"},
			IsImportant: true,
		})

		updated := projectFixtureInserted1.WithTask(task)

		err := s.storage.Replace(ctx, projectFixtureInserted1, updated)
		s.NoError(err)

		retrievedProject, err := s.storage.ByID(ctx, projectFixtureInserted1.Id)
		s.NoError(err)
		s.True(retrievedProject.Version > projectFixtureInserted1.Version)
		s.Len(retrievedProject.Tasks, 1)
	})

	s.Run("project_not_found", func() {
		task := todopb.NewTask(&todopb.AddTaskRequest{
			Title:       "to do exercises",
			ProjectId:   "2",
			UserId:      "1",
			Description: "",
			Tags:        []string{"sport", "fun"},
			IsImportant: true,
		})

		updated := projectFixture1.WithTask(task)

		err := s.storage.Replace(ctx, projectFixture1, updated)
		s.ErrorIs(err, ErrProjectNotFound)
	})

	s.Run("ids_mismatch", func() {
		task := todopb.NewTask(&todopb.AddTaskRequest{
			Title:       "to do exercises",
			ProjectId:   "2",
			UserId:      "1",
			Description: "",
			Tags:        []string{"sport", "fun"},
			IsImportant: true,
		})

		updated := projectFixtureInserted1.WithTask(task)

		err := s.storage.Replace(ctx, projectFixture1, updated)
		s.ErrorIs(err, ErrIDsMismatch)
	})

	s.Run("success", func() {
		task := todopb.NewTask(&todopb.AddTaskRequest{
			Title:       "to do exercises",
			ProjectId:   "4",
			UserId:      "3",
			Description: "",
			Tags:        []string{"sport", "fun"},
			IsImportant: true,
		})

		updated := projectFixtureInserted2.WithTask(task)

		err := s.storage.Replace(ctx, updated, projectFixtureInserted2)
		s.ErrorIs(err, ErrVersionMismatch)
	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
