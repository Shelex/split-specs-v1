package storage

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/Shelex/split-test/entities"
)

type DataStore struct {
	Client *datastore.Client
}

func NewDataStore() (Storage, error) {
	store := DataStore{}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "test-splitter")

	if err != nil {
		return store, err
	}

	//defer client.Close()

	store.Client = client

	return store, nil
}

func (d DataStore) AddProjectMaybe(projectName string) error {
	return errors.New("not implemented")
}
func (d DataStore) AddSession(projectName string, sessionID string, specs []entities.Spec) (*entities.Session, error) {
	return nil, errors.New("not implemented")
}
func (d DataStore) AttachSessionToProject(projectName string, sessionID string) error {
	return errors.New("not implemented")
}
func (d DataStore) GetProjectLatestSession(projectName string) (*entities.Session, error) {
	return nil, errors.New("not implemented")
}
func (d DataStore) SetProjectLatestSession(projectName string, sessionID string) error {
	return errors.New("not implemented")
}
func (d DataStore) GetFullProjectByName(name string) (entities.ProjectFull, error) {
	return entities.ProjectFull{}, errors.New("not implemented")
}
func (d DataStore) StartSpec(sessionID string, specName string) error {
	return errors.New("not implemented")
}
func (d DataStore) EndRunningSpec(sessionID string) error {
	return errors.New("not implemented")
}
func (d DataStore) GetSession(sessionID string) (entities.Session, error) {
	return entities.Session{}, errors.New("not implemented")
}
func (d DataStore) EndSession(sessionID string) error {
	return errors.New("not implemented")
}
