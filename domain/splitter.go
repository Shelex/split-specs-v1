package domain

import (
	"errors"
	"fmt"
	"log"

	"github.com/Shelex/split-test/entities"
	"github.com/Shelex/split-test/storage"
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

func (svc *SplitService) AddSession(projectName string, sessionID string, inputSpecs []entities.Spec) error {
	if sessionID == "" {
		return fmt.Errorf("session id cannot be empty")
	}

	specs := svc.EstimateDuration(projectName, inputSpecs)

	if err := svc.Repository.AddProjectMaybe(projectName); err != nil {
		return err
	}

	if _, err := svc.Repository.AddSession(projectName, sessionID, specs); err != nil {
		return err
	}

	log.Printf("created session %s with %d specs\n", sessionID, len(specs))

	if err := svc.Repository.AttachSessionToProject(projectName, sessionID); err != nil {
		return err
	}

	log.Printf("attached session %s to project %s", sessionID, projectName)

	return nil
}

func (svc *SplitService) EstimateDuration(projectName string, specs []entities.Spec) []entities.Spec {
	latestSession, err := svc.Repository.GetProjectLatestSession(projectName)
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
	log.Printf("requesting next spec in session %s with machine %s", sessionID, machineID)

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

	log.Printf("got spec after CalculateNext to run %v for machine %s\n", spec, machineID)

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
	log.Printf("calculating next spec from %v\n", specs)
	specsToRun := getSpecsToRun(specs)
	log.Printf("got specs to run %v\n", specsToRun)

	newSpec := getNewSpec(specsToRun)
	if newSpec.FilePath != "" {
		return newSpec
	}

	log.Printf("got newSpec to run %v\n", newSpec)

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
