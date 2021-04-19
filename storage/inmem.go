package storage

import (
	"fmt"
	"time"

	"github.com/Shelex/split-specs/entities"
)

type InMem struct {
	sessions map[string]*entities.Session
	projects map[string]*entities.Project
	users    map[string]*entities.User
}

func NewInMemStorage() (Storage, error) {
	DB = &InMem{
		sessions: map[string]*entities.Session{},
		projects: map[string]*entities.Project{},
		users:    map[string]*entities.User{},
	}
	return DB, nil
}

func (i *InMem) CreateUser(userInput entities.User) error {
	i.users[userInput.ID] = &userInput
	return nil
}

func (i *InMem) GetUserByUsername(username string) (*entities.User, error) {
	for _, user := range i.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (i *InMem) GetUserProjectIDByName(userID string, projectName string) (string, error) {
	user, ok := i.users[userID]
	if !ok {
		return "", fmt.Errorf("user not found")
	}
	for _, ID := range user.ProjectIDs {
		project, err := i.GetProjectByID(ID)
		if err == nil && project.Name == projectName {
			return project.ID, nil
		}
	}
	return "", ErrProjectNotFound
}

func (i *InMem) GetProjectByID(ID string) (*entities.Project, error) {
	project, ok := i.projects[ID]
	if !ok {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

func (i *InMem) CreateProject(project entities.Project) error {
	i.projects[project.ID] = &project
	return nil
}

func (i *InMem) AttachProjectToUser(userID string, projectID string) error {
	i.users[userID].ProjectIDs = append(i.users[userID].ProjectIDs, projectID)
	return nil
}

func (i *InMem) CreateSession(projectID string, sessionID string, specs []entities.Spec) (*entities.Session, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("[repository]: session id cannot be empty")
	}

	if _, ok := i.sessions[sessionID]; ok {
		return nil, fmt.Errorf("[repository]: session id already in use for project %s", projectID)
	}

	session := &entities.Session{
		ID:        sessionID,
		Backlog:   specs,
		ProjectID: projectID,
	}

	i.sessions[sessionID] = session
	return session, nil
}

func (i *InMem) AttachSessionToProject(projectID string, sessionID string) error {
	if _, ok := i.projects[projectID]; !ok {
		return ErrProjectNotFound
	}
	i.projects[projectID].SessionIDs = append(i.projects[projectID].SessionIDs, sessionID)
	return nil
}

func (i *InMem) GetProjectLatestSession(projectID string) (*entities.Session, error) {
	project, ok := i.projects[projectID]
	if !ok {
		return nil, ErrProjectNotFound
	}

	latestSession, ok := i.sessions[project.LatestSession]
	if !ok {
		return nil, fmt.Errorf("latest session for project %s not found", projectID)
	}

	return latestSession, nil
}

func (i *InMem) SetProjectLatestSession(projectID string, sessionID string) error {
	_, ok := i.projects[projectID]
	if !ok {
		return ErrProjectNotFound
	}
	i.projects[projectID].LatestSession = sessionID
	return nil
}

func (i *InMem) GetFullProjectByName(name string) (entities.ProjectFull, error) {
	var fullProject entities.ProjectFull

	project, ok := i.projects[name]
	if !ok {
		return fullProject, ErrProjectNotFound
	}

	fullProject.LatestSession = project.LatestSession

	for _, sessionID := range project.SessionIDs {
		session, err := i.GetSession(sessionID)
		if err != nil {
			return fullProject, fmt.Errorf("session %s not found for %s project", sessionID, name)
		}
		fullProject.Sessions = append(fullProject.Sessions, session)
	}
	return fullProject, nil
}

func (i *InMem) GetSession(sessionID string) (entities.Session, error) {
	var empty entities.Session
	session, ok := i.sessions[sessionID]
	if !ok {
		return empty, fmt.Errorf("session %s not found", sessionID)

	}
	return *session, nil
}

func (i *InMem) StartSpec(sessionID string, machineID string, specName string) error {
	session, err := i.GetSession(sessionID)
	if err != nil {
		return err
	}

	for index, spec := range session.Backlog {
		if spec.FilePath == specName {
			if session.Start == 0 {
				i.sessions[sessionID].Start = time.Now().Unix()
			}
			i.sessions[sessionID].Backlog[index].Start = time.Now().Unix()
			i.sessions[sessionID].Backlog[index].AssignedTo = machineID
			return nil
		}
	}
	return nil
}

func (i *InMem) EndSpec(sessionID string, machineID string) error {
	session, err := i.GetSession(sessionID)
	if err != nil {
		return err
	}
	for index, spec := range session.Backlog {
		if spec.End == 0 && spec.Start != 0 && spec.AssignedTo == machineID {
			backlogItem := i.sessions[sessionID].Backlog[index]
			backlogItem.End = time.Now().Unix()
			backlogItem.EstimatedDuration = backlogItem.End - backlogItem.Start
			i.sessions[sessionID].Backlog[index] = backlogItem
			return nil
		}
	}
	return nil
}

func (i *InMem) EndSession(sessionID string) error {
	session, err := i.GetSession(sessionID)
	if err != nil {
		return err
	}

	i.sessions[sessionID].End = time.Now().Unix()

	if err := i.SetProjectLatestSession(session.ProjectID, sessionID); err != nil {
		return err
	}

	return nil
}
