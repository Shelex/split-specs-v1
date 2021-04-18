package storage

import (
	"fmt"
	"log"
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

func (i *InMem) AddProjectMaybe(projectName string) error {
	_, ok := i.projects[projectName]
	if ok {
		return nil
	}
	i.projects[projectName] = &entities.Project{}
	return nil
}

func (i *InMem) AddSession(projectName string, sessionID string, specs []entities.Spec) (*entities.Session, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("[repository]: session id cannot be empty")
	}

	if _, ok := i.sessions[sessionID]; ok {
		return nil, fmt.Errorf("[repository]: session id already in use for project %s", projectName)
	}

	session := &entities.Session{
		ID:          sessionID,
		Backlog:     specs,
		ProjectName: projectName,
	}

	i.sessions[sessionID] = session
	log.Printf("created session %s with %d specs\n", sessionID, len(specs))
	return session, nil
}

func (i *InMem) AttachSessionToProject(projectName string, sessionID string) error {
	if _, ok := i.projects[projectName]; !ok {
		return fmt.Errorf("[repository]: project %s not found", projectName)
	}
	i.projects[projectName].Sessions = append(i.projects[projectName].Sessions, sessionID)
	return nil
}

func (i *InMem) GetProjectLatestSession(projectName string) (*entities.Session, error) {
	project, ok := i.projects[projectName]
	if !ok {
		return nil, fmt.Errorf("[repository]: project %s not found", projectName)
	}

	latestSession, ok := i.sessions[project.LatestSession]
	if !ok {
		return nil, fmt.Errorf("[repository]: latest session for project %s not found", projectName)
	}

	return latestSession, nil
}

func (i *InMem) SetProjectLatestSession(projectName string, sessionID string) error {
	_, ok := i.projects[projectName]
	if !ok {
		return fmt.Errorf("[repository]: project %s not found", projectName)
	}
	i.projects[projectName].LatestSession = sessionID
	return nil
}

func (i *InMem) GetFullProjectByName(name string) (entities.ProjectFull, error) {
	var fullProject entities.ProjectFull

	project, ok := i.projects[name]
	if !ok {
		return fullProject, fmt.Errorf("[repository]: project %s not found", name)
	}

	fullProject.LatestSession = project.LatestSession

	for _, sessionID := range project.Sessions {
		session, err := i.GetSession(sessionID)
		if err != nil {
			return fullProject, fmt.Errorf("[repository]: session %s not found for %s project", sessionID, name)
		}
		fullProject.Sessions = append(fullProject.Sessions, session)
	}
	return fullProject, nil
}

func (i *InMem) GetSession(sessionID string) (entities.Session, error) {
	var empty entities.Session
	session, ok := i.sessions[sessionID]
	if !ok {
		return empty, fmt.Errorf("[repository]: session %s not found", sessionID)

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
			log.Printf("started spec %s in session %s for machine %s", spec.FilePath, sessionID, machineID)
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
			log.Printf("finished spec %s in session %s for machine %s", spec.FilePath, sessionID, machineID)
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

	if err := i.SetProjectLatestSession(session.ProjectName, sessionID); err != nil {
		return err
	}

	return nil
}
