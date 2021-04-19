package domain

import (
	"errors"
	"fmt"

	"github.com/Shelex/split-specs/entities"
	"github.com/Shelex/split-specs/storage"
	uuid "github.com/satori/go.uuid"
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
		if errors.Is(err, storage.ErrProjectNotFound) {
			newID, err := svc.AddProject(userID, projectName)
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

func (svc *SplitService) AddProject(userID string, projectName string) (string, error) {
	id := uuid.NewV4().String()
	if err := svc.Repository.CreateProject(entities.Project{
		ID:   id,
		Name: projectName,
	}); err != nil {
		return "", err
	}

	if err := svc.Repository.AttachProjectToUser(userID, id); err != nil {
		return "", err
	}

	return id, nil
}

func (svc *SplitService) EstimateDuration(projectID string, specs []entities.Spec) []entities.Spec {
	latestSession, err := svc.Repository.GetProjectLatestSession(projectID)
	if err != nil {
		return specs
	}

	for idx, spec := range specs {
		for _, history := range latestSession.Backlog {
			if history.FilePath == spec.FilePath {
				specs[idx].EstimatedDuration = history.EstimatedDuration
			}
		}
	}
	return specs
}

func (svc *SplitService) GetProject(name string) (entities.ProjectFull, error) {
	var empty entities.ProjectFull

	project, err := svc.Repository.GetFullProjectByName(name)

	if err != nil {
		return empty, err
	}

	return project, nil
}

func (svc *SplitService) Next(sessionID string, machineID string) (string, error) {
	if err := svc.Repository.EndSpec(sessionID, machineID); err != nil {
		return "", err
	}

	session, err := svc.Repository.GetSession(sessionID)
	if err != nil {
		return "", err
	}
	if len(session.Backlog) == 0 {
		return "", fmt.Errorf("backlog for session %s is empty", sessionID)
	}

	spec := svc.CalculateNext(session.Backlog)

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
