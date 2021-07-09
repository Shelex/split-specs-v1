package storage

import (
	"errors"

	"github.com/Shelex/split-specs/entities"
)

var DB Storage

type Storage interface {
	GetProjectByID(ID string) (*entities.Project, error)
	GetUserProjectIDByName(userID string, projectName string) (string, error)
	GetUserProjects(userID string) ([]string, error)

	CreateProject(project entities.Project) error
	AttachProjectToUser(userID string, projectID string) error

	GetSession(sessionID string) (entities.Session, error)
	CreateSession(projectName string, sessionID string, specs []entities.Spec) (*entities.Session, error)
	AttachSessionToProject(projectName string, sessionID string) error
	EndSession(sessionID string) error

	GetProjectLatestSession(projectID string) (*entities.Session, error)
	SetProjectLatestSession(projectName string, sessionID string) error

	CreateSpecs(sessionID string, specs []entities.Spec) ([]string, error)
	GetSpecs(sessionID string, ids []string) ([]entities.Spec, error)
	StartSpec(sessionID string, machineID string, specName string) error
	EndSpec(sessionID string, machineID string) error

	//auth
	CreateUser(user entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	UpdatePassword(userID string, newPassword string) error
}

var ErrProjectNotFound = errors.New("project not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrSpecNotFound = errors.New("spec not found")
var ErrSessionFinished = errors.New("session already finished")
