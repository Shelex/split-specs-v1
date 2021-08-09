package storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Shelex/split-specs/entities"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"google.golang.org/appengine"
)

const (
	DATASTORE_PROJECT_ID = "split-specs"
	userKind             = "users"
	projectKind          = "projects"
	sessionKind          = "sessions"
	specKind             = "specs"
)

type DataStore struct {
	Client *datastore.Client
	ctx    context.Context
}

func NewDataStore() (Storage, error) {
	ctx := appengine.BackgroundContext()
	client, err := datastore.NewClient(ctx, DATASTORE_PROJECT_ID)
	if err != nil {
		return nil, err
	}

	//defer client.Close()

	DB = DataStore{
		ctx:    ctx,
		Client: client,
	}

	return DB, nil
}

func (d DataStore) CreateProject(project entities.Project) error {
	projectKey := datastore.NameKey(projectKind, project.ID, nil)
	_, err := d.Client.Put(d.ctx, projectKey, &project)
	if err != nil {
		return err
	}
	return nil
}
func (d DataStore) CreateSession(projectID string, sessionID string, specs []entities.Spec) (*entities.Session, error) {
	sessionKey := datastore.NameKey(sessionKind, sessionID, nil)

	var session *entities.Session

	err := d.CreateSpecs(sessionID, specs)
	if err != nil {
		return nil, err
	}
	session = &entities.Session{
		ID:        sessionID,
		ProjectID: projectID,
	}
	if _, err := d.Client.Put(d.ctx, sessionKey, session); err != nil {
		return nil, err
	}

	return session, err
}

func (d DataStore) GetProjectLatestSessions(projectID string, limit int) ([]*entities.Session, error) {
	sessionQuery := datastore.NewQuery(sessionKind).Filter("projectId=", projectID).Filter("end>", 0).Order("-end").Limit(limit)

	var sessions []*entities.Session

	if _, err := d.Client.GetAll(d.ctx, sessionQuery, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (d DataStore) GetProjectLatestSession(projectID string) (*entities.Session, error) {
	project, err := d.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}

	if project.LatestSession == "" {
		return nil, ErrSessionNotFound
	}

	session, err := d.GetSession(project.LatestSession)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
func (d DataStore) SetProjectLatestSession(projectID string, sessionID string) error {
	project, err := d.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	project.LatestSession = sessionID

	projectKey := datastore.NameKey(projectKind, projectID, nil)

	if _, err := d.Client.Put(d.ctx, projectKey, project); err != nil {
		return err
	}
	return nil
}

func (d DataStore) StartSpec(sessionID string, machineID string, specName string) error {
	session, err := d.GetSession(sessionID)
	if err != nil {
		return err
	}

	specs, err := d.GetSpecs(sessionID)
	if err != nil {
		return err
	}

	var startedSpec entities.Spec

	for _, spec := range specs {
		if spec.FilePath == specName {
			startedSpec = spec
			break
		}
	}

	if startedSpec.FilePath == "" {
		return nil
	}

	startedSpec.Start = time.Now().Unix()
	startedSpec.AssignedTo = machineID

	if session.Start == 0 {
		session.Start = time.Now().Unix()
	}

	tx, err := d.Client.NewTransaction(d.ctx)
	if err != nil {
		return err
	}
	sessionKey := datastore.NameKey(sessionKind, session.ID, nil)

	if _, err := tx.Put(sessionKey, &session); err != nil {
		return fmt.Errorf("failed to write session start: %s", err)
	}

	specKey := datastore.NameKey(specKind, startedSpec.ID, sessionKey)

	if _, err := tx.Put(specKey, &startedSpec); err != nil {
		return fmt.Errorf("failed to write spec start: %s", err)
	}
	if _, err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit trx start spec and session: %s", err)
	}
	return nil
}
func (d DataStore) EndSpec(sessionID string, machineID string, isPassed bool) error {
	session, _ := d.GetSession(sessionID)
	if session.ID == "" {
		return nil
	}

	specs, err := d.GetSpecs(sessionID)
	if err != nil {
		return err
	}

	var finishedSpec entities.Spec

	for _, spec := range specs {
		if spec.End == 0 && spec.Start != 0 && spec.AssignedTo == machineID {
			finishedSpec = spec
			break
		}
	}

	if finishedSpec.FilePath == "" {
		return nil
	}

	finishedSpec.End = time.Now().Unix()
	finishedSpec.EstimatedDuration = finishedSpec.End - finishedSpec.Start
	finishedSpec.Passed = isPassed

	sessionKey := datastore.NameKey(sessionKind, session.ID, nil)
	specKey := datastore.NameKey(specKind, finishedSpec.ID, sessionKey)

	if _, err := d.Client.Put(d.ctx, specKey, &finishedSpec); err != nil {
		return err
	}
	return nil
}
func (d DataStore) GetSession(sessionID string) (entities.Session, error) {
	sessionQuery := datastore.NewQuery(sessionKind).Filter("id=", sessionID).Limit(1)

	var sessions []entities.Session

	if _, err := d.Client.GetAll(d.ctx, sessionQuery, &sessions); err != nil {
		return entities.Session{}, err
	}

	if len(sessions) == 0 {
		return entities.Session{}, ErrSessionNotFound
	}

	return sessions[0], nil
}
func (d DataStore) EndSession(sessionID string) error {
	session, _ := d.GetSession(sessionID)
	if session.ID == "" {
		return ErrSessionFinished
	}
	if session.End != 0 {
		return ErrSessionFinished
	}

	session.End = time.Now().Unix()

	sessionKey := datastore.NameKey(sessionKind, sessionID, nil)

	if _, err := d.Client.Put(d.ctx, sessionKey, &session); err != nil {
		return err
	}
	return d.SetProjectLatestSession(session.ProjectID, sessionID)
}

func (d DataStore) CreateUser(user entities.User) error {
	userKey := datastore.NameKey(userKind, user.ID, nil)

	_, err := d.Client.Put(d.ctx, userKey, &user)
	if err != nil {
		return err
	}
	return nil
}

func (d DataStore) GetUserByEmail(email string) (*entities.User, error) {
	query := datastore.NewQuery(userKind).Filter("email=", email).Limit(1)

	var users []entities.User

	if _, err := d.Client.GetAll(d.ctx, query, &users); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &users[0], nil
}

func (d DataStore) GetUserProjectIDByName(userID string, projectName string) (string, error) {
	userKey := datastore.NameKey(userKind, userID, nil)

	var user entities.User

	err := d.Client.Get(d.ctx, userKey, &user)
	if len(user.ProjectIDs) == 0 {
		return "", ErrProjectNotFound
	}
	if err != nil {
		return "", err
	}

	projectKeys := make([]*datastore.Key, len(user.ProjectIDs))

	for index, projectID := range user.ProjectIDs {
		projectKeys[index] = datastore.NameKey(projectKind, projectID, nil)
	}

	projects := make([]entities.Project, len(user.ProjectIDs))

	if err := d.Client.GetMulti(d.ctx, projectKeys, projects); err != nil {
		return "", err
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project.ID, nil
		}
	}
	return "", ErrProjectNotFound
}

func (d DataStore) GetProjectByID(ID string) (*entities.Project, error) {
	projectKey := datastore.NameKey(projectKind, ID, nil)
	var project entities.Project

	err := d.Client.Get(d.ctx, projectKey, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (d DataStore) AttachProjectToUser(userID string, projectID string) error {
	userKey := datastore.NameKey(userKind, userID, nil)

	var user entities.User
	if err := d.Client.Get(d.ctx, userKey, &user); err != nil {
		return err
	}

	if user.ProjectIDs == nil {
		user.ProjectIDs = []string{projectID}
	} else {
		user.ProjectIDs = append(user.ProjectIDs, projectID)
	}
	if _, err := d.Client.Put(d.ctx, userKey, &user); err != nil {
		return err
	}
	return nil
}

func (d DataStore) UnlinkProjectFromUser(email string, projectID string) error {
	user, err := d.GetUserByEmail(email)

	if err != nil {
		return err
	}

	userHasProject, index := contains(user.ProjectIDs, projectID)

	if !userHasProject {
		return ErrProjectNotFound
	}

	user.ProjectIDs = remove(user.ProjectIDs, index)

	userKey := datastore.NameKey(userKind, user.ID, nil)

	if _, err := d.Client.Put(d.ctx, userKey, user); err != nil {
		return err
	}
	return nil
}

func (d DataStore) GetUserProjects(userID string) ([]string, error) {
	userKey := datastore.NameKey(userKind, userID, nil)

	var user entities.User
	if err := d.Client.Get(d.ctx, userKey, &user); err != nil {
		return nil, err
	}
	return user.ProjectIDs, nil
}

func (d DataStore) UpdatePassword(userID string, newPassword string) error {
	userKey := datastore.NameKey(userKind, userID, nil)

	var user entities.User
	if err := d.Client.Get(d.ctx, userKey, &user); err != nil {
		return err
	}
	user.Password = newPassword

	if _, err := d.Client.Put(d.ctx, userKey, &user); err != nil {
		return err
	}
	return nil
}

func (d DataStore) CreateSpecs(sessionID string, specs []entities.Spec) error {
	sessionKey := datastore.NameKey(sessionKind, sessionID, nil)

	specKeys := make([]*datastore.Key, len(specs))

	for index, spec := range specs {
		id, _ := gonanoid.New()
		spec.ID = id
		spec.SessionID = sessionID
		specKeys[index] = datastore.NameKey(specKind, spec.ID, sessionKey)
		specs[index] = spec
	}

	if _, err := d.Client.PutMulti(d.ctx, specKeys, specs); err != nil {
		return err
	}
	return nil
}

func (d DataStore) GetSpecs(sessionID string) ([]entities.Spec, error) {
	sessionKey := datastore.NameKey(sessionKind, sessionID, nil)
	query := datastore.NewQuery(specKind).Ancestor(sessionKey)

	specs := make([]entities.Spec, 0)

	if _, err := d.Client.GetAll(d.ctx, query, &specs); err != nil {
		return nil, err
	}
	return specs, nil
}

func (d DataStore) DeleteSession(email string, sessionID string) error {
	sessionKey := datastore.NameKey(sessionKind, sessionID, nil)

	session, err := d.GetSession(sessionID)
	if err != nil {
		return err
	}

	user, err := d.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if userHasSession, _ := contains(user.ProjectIDs, session.ProjectID); !userHasSession {
		return ErrSessionNotFound
	}

	specs, err := d.GetSpecs(sessionID)
	if err != nil {
		return err
	}

	specKeys := make([]*datastore.Key, len(specs))

	for index, spec := range specs {
		specKeys[index] = datastore.NameKey(specKind, spec.ID, sessionKey)
	}

	tx, err := d.Client.NewTransaction(d.ctx)
	if err != nil {
		return err
	}

	if err := tx.DeleteMulti(specKeys); err != nil {
		return err
	}

	if err := tx.Delete(sessionKey); err != nil {
		return err
	}

	if _, err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d DataStore) DeleteProject(email string, projectID string) error {
	if err := d.UnlinkProjectFromUser(email, projectID); err != nil {
		return err
	}

	sessions, err := d.GetProjectSessions(projectID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := d.DeleteSession(email, session.ID); err != nil {
			return err
		}
	}

	projectKey := datastore.NameKey(projectKind, projectID, nil)

	return d.Client.Delete(d.ctx, projectKey)
}

func (d DataStore) GetSessionWithSpecs(sessionID string) (entities.SessionWithSpecs, error) {
	var empty entities.SessionWithSpecs

	session, err := d.GetSession(sessionID)
	if err != nil {
		return empty, err
	}

	specs, err := d.GetSpecs(sessionID)
	if err != nil {
		return empty, err
	}

	return entities.SessionWithSpecs{
		ID:        session.ID,
		ProjectID: session.ProjectID,
		Start:     session.Start,
		End:       session.End,
		Specs:     specs,
	}, nil

}

func (d DataStore) GetProjectSessions(projectID string) ([]entities.SessionWithSpecs, error) {
	sessionQuery := datastore.NewQuery(sessionKind).Filter("projectId=", projectID)

	var sessions []entities.SessionWithSpecs

	if _, err := d.Client.GetAll(d.ctx, sessionQuery, &sessions); err != nil {
		return nil, err
	}

	for index, session := range sessions {
		sessionKey := datastore.NameKey(sessionKind, session.ID, nil)
		specsQuery := datastore.NewQuery(specKind).Ancestor(sessionKey)

		var specs []entities.Spec

		if _, err := d.Client.GetAll(d.ctx, specsQuery, &specs); err != nil {
			return nil, err
		}
		sessions[index].Specs = specs
	}
	return sessions, nil
}
