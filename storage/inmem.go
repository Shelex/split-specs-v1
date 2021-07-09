package storage

import (
	"fmt"
	"time"

	"github.com/Shelex/split-specs/entities"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type InMem struct {
	sessions map[string]*entities.Session
	projects map[string]*entities.Project
	users    map[string]*entities.User
	specs    map[string]*entities.Spec
}

func NewInMemStorage() (Storage, error) {
	DB = &InMem{
		sessions: map[string]*entities.Session{},
		projects: map[string]*entities.Project{},
		users:    map[string]*entities.User{},
		specs:    map[string]*entities.Spec{},
	}
	return DB, nil
}

func (i *InMem) CreateUser(userInput entities.User) error {
	i.users[userInput.ID] = &userInput
	return nil
}

func (i *InMem) UpdatePassword(userID string, newPassword string) error {
	i.users[userID].Password = newPassword
	return nil
}

func (i *InMem) GetUserByEmail(email string) (*entities.User, error) {
	for _, user := range i.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (i *InMem) GetUserProjects(userID string) ([]string, error) {
	return i.users[userID].ProjectIDs, nil
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
	if _, ok := i.sessions[sessionID]; ok {
		return nil, fmt.Errorf("[repository]: session id already in use for project %s", projectID)
	}

	specIds, err := i.CreateSpecs(sessionID, specs)
	if err != nil {
		return nil, fmt.Errorf("failed to create specs")
	}

	session := &entities.Session{
		ID:        sessionID,
		SpecIDs:   specIds,
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

	specs, err := i.GetSpecs(sessionID, session.SpecIDs)
	if err != nil {
		return err
	}

	for _, spec := range specs {
		if spec.FilePath == specName {
			if session.Start == 0 {
				i.sessions[sessionID].Start = time.Now().Unix()
			}
			i.specs[spec.ID].Start = time.Now().Unix()
			i.specs[spec.ID].AssignedTo = machineID
			return nil
		}
	}
	return nil
}

func (i *InMem) EndSpec(sessionID string, machineID string, isPassed bool) error {
	session, err := i.GetSession(sessionID)
	if err != nil {
		return err
	}

	specs, err := i.GetSpecs(sessionID, session.SpecIDs)
	if err != nil {
		return err
	}

	for _, spec := range specs {
		if spec.End == 0 && spec.Start != 0 && spec.AssignedTo == machineID {
			i.specs[spec.ID].End = time.Now().Unix()
			i.specs[spec.ID].EstimatedDuration = i.specs[spec.ID].End - i.specs[spec.ID].Start
			i.specs[spec.ID].Passed = isPassed
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

func (i *InMem) CreateSpecs(sessionID string, specs []entities.Spec) ([]string, error) {
	ids := make([]string, len(specs))
	for index, spec := range specs {
		id, _ := gonanoid.New()
		spec.ID = id
		spec.SessionID = sessionID
		i.specs[spec.ID] = &spec
		ids[index] = spec.ID
	}
	return ids, nil
}

func (i *InMem) GetSpecs(sessionID string, ids []string) ([]entities.Spec, error) {
	specs := make([]entities.Spec, len(ids))
	for index, id := range ids {
		spec, ok := i.specs[id]
		if !ok {
			return nil, ErrSpecNotFound
		}
		specs[index] = *spec
	}
	return specs, nil
}
