package storage

import (
	"github.com/Shelex/split-test/entities"
)

type Storage interface {
	AddProjectMaybe(projectName string) error
	AddSession(projectName string, sessionID string, specs []entities.Spec) (*entities.Session, error)
	AttachSessionToProject(projectName string, sessionID string) error
	GetProjectLatestSession(projectName string) (*entities.Session, error)
	SetProjectLatestSession(projectName string, sessionID string) error
	GetFullProjectByName(name string) (entities.ProjectFull, error)
	StartSpec(sessionID string, specName string) error
	EndRunningSpec(sessionID string) error
	GetSession(sessionID string) (entities.Session, error)
	EndSession(sessionID string) error
}
