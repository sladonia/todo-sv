package todo

import (
	"context"
	"sync"
)

type Storage interface {
	ByID(ctx context.Context, projectID string) (Project, error)
	AllUserProjects(ctx context.Context, userID string) ([]Project, error)
	Insert(ctx context.Context, project Project) error
	Replace(ctx context.Context, prev, curr Project) error
	Delete(ctx context.Context, projectID string) error
}

type inMemoryStorage struct {
	sync.RWMutex
	projects map[string]Project
}

func NewInMemoryStorage() Storage {
	return &inMemoryStorage{
		projects: make(map[string]Project),
	}
}

func (s *inMemoryStorage) ByID(_ context.Context, projectID string) (Project, error) {
	s.RLock()
	defer s.RUnlock()

	p, ok := s.projects[projectID]
	if !ok {
		return Project{}, ErrProjectNotFound
	}

	return p, nil
}

func (s *inMemoryStorage) AllUserProjects(_ context.Context, userID string) ([]Project, error) {
	s.RLock()
	defer s.RUnlock()

	var userProjects []Project

	for _, project := range s.projects {
		if project.OwnerID == userID {
			userProjects = append(userProjects, project)
			continue
		}

		if project.Participants.Contains(userID) {
			userProjects = append(userProjects, project)
			continue
		}
	}

	return userProjects, nil
}

func (s *inMemoryStorage) Insert(_ context.Context, project Project) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.projects[project.ID]
	if ok {
		return ErrAlreadyExists
	}

	s.projects[project.ID] = project

	return nil
}

func (s *inMemoryStorage) Replace(_ context.Context, prev, curr Project) error {
	s.Lock()
	defer s.Unlock()

	if prev.ID != curr.ID {
		return ErrIDsMismatch
	}

	_, ok := s.projects[prev.ID]
	if !ok {
		return ErrProjectNotFound
	}

	if prev.Version > curr.Version {
		return ErrVersionMismatch
	}

	s.projects[curr.ID] = curr

	return nil
}

func (s *inMemoryStorage) Delete(ctx context.Context, projectID string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.projects, projectID)

	return nil
}
