package todo

import (
	"context"
	"errors"

	"github.com/sladonia/todo-sv/pkg/todopb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStorage struct {
	db      *mongo.Database
	colName string
}

func NewStorage(db *mongo.Database, colName string) Storage {
	return &mongoStorage{
		db:      db,
		colName: colName,
	}
}

func (s *mongoStorage) ByID(ctx context.Context, projectID string) (*todopb.Project, error) {
	var projectBSON ProjectBSON

	res := s.collection().FindOne(ctx, bson.M{"_id": projectID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, ErrProjectNotFound
		}

		return nil, res.Err()
	}

	err := res.Decode(&projectBSON)
	if err != nil {
		return nil, err
	}

	return projectBSON.Project(), nil
}

func (s *mongoStorage) AllUserProjects(ctx context.Context, userID string) ([]*todopb.Project, error) {
	cur, err := s.collection().Find(
		ctx,
		bson.D{
			{"$or", bson.A{
				bson.M{"owner_id": userID},
				bson.M{"participants": userID},
			}}},
	)
	if err != nil {
		return nil, err
	}

	var projectsBSON []ProjectBSON

	err = cur.All(ctx, &projectsBSON)
	if err != nil {
		return nil, err
	}

	var projects []*todopb.Project

	for _, projBSON := range projectsBSON {
		projects = append(projects, projBSON.Project())
	}

	return projects, nil
}

func (s *mongoStorage) Insert(ctx context.Context, project *todopb.Project) error {
	projectBSON := NewProjectBSON(project)

	_, err := s.collection().InsertOne(ctx, projectBSON)
	if err != nil {
		if IsDuplicateKeyError(err) {
			return ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (s *mongoStorage) Replace(ctx context.Context, prev, curr *todopb.Project) error {
	if prev.Id != curr.Id {
		return ErrIDsMismatch
	}

	if prev.Version > curr.Version {
		return ErrVersionMismatch
	}

	curBSON := NewProjectBSON(curr)

	res := s.collection().FindOneAndReplace(ctx, bson.M{"_id": prev.Id}, curBSON)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return ErrProjectNotFound
		}

		return res.Err()
	}

	return nil
}

func (s *mongoStorage) Delete(ctx context.Context, projectID string) error {
	_, err := s.collection().DeleteOne(ctx, bson.M{"_id": projectID})
	return err
}

func (s *mongoStorage) collection() *mongo.Collection {
	return s.db.Collection(s.colName)
}

//TODO: add unique indexes
