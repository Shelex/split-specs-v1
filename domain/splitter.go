package domain

import (
	"errors"
	"fmt"

	"github.com/Shelex/split-specs/entities"
	"github.com/Shelex/split-specs/storage"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"google.golang.org/appengine/datastore"
)

var ErrSessionFinished = errors.New("session finished")

type SplitService struct {
	Repository storage.Storage
}

func NewSplitService(repo storage.Storage) SplitService {
	return SplitService{
		Repository: repo,
	}
}

func (svc *SplitService) AddSession(userID string, projectName string, sessionID string, inputSpecs []entities.Spec) error {
	if sessionID == "" {
		return fmt.Errorf("session id cannot be empty")
	}

	projectID, err := svc.Repository.GetUserProjectIDByName(userID, projectName)

	if err != nil {
		if err.Error() == storage.ErrProjectNotFound.Error() || err.Error() == datastore.ErrNoSuchEntity.Error() {
			newID, err := svc.AddProject(userID, projectName, sessionID)
			if err != nil {
				return err
			}
			projectID = newID
		} else {
			return err
		}
	}

	specs := svc.EstimateDuration(projectID, inputSpecs)

	if _, err := svc.Repository.CreateSession(projectID, sessionID, specs); err != nil {
		return err
	}

	if err := svc.Repository.AttachSessionToProject(projectID, sessionID); err != nil {
		return err
	}

	return nil
}

func (svc *SplitService) AddProject(userID string, projectName string, sessionID string) (string, error) {
	id, _ := gonanoid.New()

	if err := svc.Repository.CreateProject(entities.Project{
		ID:         id,
		Name:       projectName,
		SessionIDs: []string{sessionID},
	}); err != nil {
		return "", err
	}

	if err := svc.Repository.AttachProjectToUser(userID, id); err != nil {
		return "", err
	}

	return id, nil
}

func (svc *SplitService) InviteUserToProject(user entities.User, guest string, projectName string) error {
	projectID, err := svc.Repository.GetUserProjectIDByName(user.ID, projectName)
	if err != nil {
		return fmt.Errorf("failed to share project")
	}

	guestUser, err := svc.Repository.GetUserByEmail(guest)
	if err != nil {
		return fmt.Errorf("failed to share project")
	}

	if _, err := svc.Repository.GetUserProjectIDByName(guestUser.ID, projectName); err != nil {
		return fmt.Errorf("user already has project with such name")
	}

	return svc.Repository.AttachProjectToUser(guestUser.ID, projectID)
}

func (svc *SplitService) EstimateDuration(projectID string, specs []entities.Spec) []entities.Spec {
	latestSession, err := svc.Repository.GetProjectLatestSession(projectID)
	if err != nil {
		return specs
	}
	historySpecs, err := svc.Repository.GetSpecs(latestSession.ID, latestSession.SpecIDs)
	if err != nil {
		return specs
	}

	for idx, spec := range specs {
		for _, history := range historySpecs {
			if history.FilePath == spec.FilePath {
				specs[idx].EstimatedDuration = history.EstimatedDuration
			}
		}
	}
	return specs
}

func (svc *SplitService) GetProjectList(user entities.User) ([]string, error) {
	projectIds, err := svc.Repository.GetUserProjects(user.ID)
	if err != nil {
		return []string{}, err
	}

	projects := make([]string, len(projectIds))

	for index, id := range projectIds {
		project, err := svc.Repository.GetProjectByID(id)
		if err != nil {
			return []string{}, err
		}
		projects[index] = project.Name
	}
	return projects, nil
}

func (svc *SplitService) Next(sessionID string, machineID string) (string, error) {
	if err := svc.Repository.EndSpec(sessionID, machineID); err != nil {
		if err.Error() != datastore.ErrNoSuchEntity.Error() {
			return "", storage.ErrSessionNotFound
		}
	}

	session, err := svc.Repository.GetSession(sessionID)
	if err.Error() != datastore.ErrNoSuchEntity.Error() {
		return "", storage.ErrSessionNotFound
	}
	if err != nil {
		return "", err
	}
	if len(session.SpecIDs) == 0 {
		return "", fmt.Errorf("backlog for session %s is empty", sessionID)
	}

	specs, err := svc.Repository.GetSpecs(sessionID, session.SpecIDs)
	if err != nil {
		return "", err
	}

	spec := svc.CalculateNext(specs)

	if spec.FilePath == "" {
		if err := svc.Repository.EndSession(sessionID); err != nil {
			return "", err
		}
		return "", ErrSessionFinished
	}

	if err := svc.Repository.StartSpec(sessionID, machineID, spec.FilePath); err != nil {
		return "", err
	}

	return spec.FilePath, nil
}

func (svc *SplitService) CalculateNext(specs []entities.Spec) entities.Spec {
	specsToRun := getSpecsToRun(specs)

	newSpec := getNewSpec(specsToRun)
	if newSpec.FilePath != "" {
		return newSpec
	}

	return getLongestSpec(specsToRun)
}

func getLongestSpec(specs []entities.Spec) entities.Spec {
	longestSpec := entities.Spec{}

	for _, spec := range specs {
		if spec.EstimatedDuration > longestSpec.EstimatedDuration {
			longestSpec = spec
		}
	}

	return longestSpec
}

func getSpecsToRun(specs []entities.Spec) []entities.Spec {
	filtered := make([]entities.Spec, 0)
	for _, spec := range specs {
		if spec.Start == 0 {
			filtered = append(filtered, spec)
		}
	}
	return filtered
}

func getNewSpec(specs []entities.Spec) entities.Spec {
	for _, spec := range specs {
		if spec.EstimatedDuration == 0 {
			return spec
		}
	}
	return entities.Spec{}
}
