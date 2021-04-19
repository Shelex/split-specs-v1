package storage

import (
	"errors"

	"github.com/Shelex/split-specs/entities"
)

var DB Storage

type Storage interface {
	GetProjectByID(ID string) (*entities.Project, error)
	GetFullProjectByName(name string) (entities.ProjectFull, error)
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

	StartSpec(sessionID string, machineID string, specName string) error
	EndSpec(sessionID string, machineID string) error

	//auth
	CreateUser(user entities.User) error
	GetUserByUsername(username string) (*entities.User, error)
	UpdatePassword(userID string, newPassword string) error
}

var ErrProjectNotFound = errors.New("project not found")
var ErrSessionNotFound = errors.New("project not found")
