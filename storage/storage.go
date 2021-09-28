package storage

import (
	"errors"

	"github.com/Shelex/split-specs/entities"
)

var DB Storage

type Storage interface {
	GetProjectByID(ID string) (*entities.Project, error)
	GetUserProjectIDByName(userID string, projectName string) (string, error)
	GetUserProjectIDs(userID string) ([]string, error)

	CreateProject(project entities.Project) error
	AttachProjectToUser(userID string, projectID string) error
	DeleteProject(email string, projectID string) error
	GetProjectSessions(projectID string) ([]entities.SessionWithSpecs, error)
	GetProjectUsers(projectID string) ([]string, error)

	GetSession(sessionID string) (entities.Session, error)
	GetSessionWithSpecs(sessionID string) (entities.SessionWithSpecs, error)
	CreateSession(projectName string, sessionID string, specs []entities.Spec) (*entities.Session, error)
	EndSession(sessionID string) error
	DeleteSession(email string, sessionID string) error

	GetProjectLatestSessions(projectID string, limit int) ([]*entities.Session, error)
	GetProjectLatestSession(projectID string) (*entities.Session, error)
	SetProjectLatestSession(projectName string, sessionID string) error

	CreateSpecs(sessionID string, specs []entities.Spec) error
	GetSpec(specID string) (entities.Spec, error)
	GetSpecs(sessionID string) ([]entities.Spec, error)
	StartSpec(sessionID string, machineID string, specID string) error
	EndSpec(sessionID string, machineID string, isPassed bool) error

	//auth
	CreateUser(user entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	UpdatePassword(userID string, newPassword string) error

	//api keys
	CreateApiKey(userID string, key entities.ApiKey) error
	DeleteApiKey(userID string, keyID string) error
	GetApiKeys(userID string) ([]entities.ApiKey, error)
	GetApiKey(userID string, keyID string) (entities.ApiKey, error)
}

var ErrProjectNotFound = errors.New("project not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrSpecNotFound = errors.New("spec not found")
var ErrSessionFinished = errors.New("session already finished")
var ErrApiKeyNotFound = errors.New("api key not found")
