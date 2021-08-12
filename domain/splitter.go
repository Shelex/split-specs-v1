package domain

import (
	"errors"
	"fmt"
	"math"

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

	return nil
}

func (svc *SplitService) AddProject(userID string, projectName string, sessionID string) (string, error) {
	id, _ := gonanoid.New()

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
	latestSessions, err := svc.Repository.GetProjectLatestSessions(projectID, 5)
	if err != nil {
		return specs
	}

	var historicalSpecs []entities.Spec

	for _, session := range latestSessions {
		sessionSpecs, err := svc.Repository.GetSpecs(session.ID)
		if err != nil {
			return specs
		}
		historicalSpecs = append(historicalSpecs, sessionSpecs...)
	}

	averages := make(map[string]float64)

	for _, historicalSpec := range historicalSpecs {
		if historicalSpec.End == 0 {
			continue
		}

		// handle case when estimated duration is 0, because of milliseconds rounding
		// set it to 1 in order to separate from new specs
		// TODO: migrate to float64, timestamps should include milliseconds
		if historicalSpec.EstimatedDuration == 0 {
			historicalSpec.EstimatedDuration = 1
		}

		_, ok := averages[historicalSpec.FilePath]
		if !ok {
			averages[historicalSpec.FilePath] = float64(historicalSpec.EstimatedDuration)
		}

		average := (averages[historicalSpec.FilePath] + float64(historicalSpec.EstimatedDuration)) / 2
		averages[historicalSpec.FilePath] = average
	}

	for index, spec := range specs {
		average, ok := averages[spec.FilePath]
		if ok {
			specs[index].EstimatedDuration = int64(math.Round(average))
		}
	}

	return specs
}

func (svc *SplitService) GetProjectList(user entities.User) ([]string, error) {
	projectIds, err := svc.Repository.GetUserProjectIDs(user.ID)
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

func (svc *SplitService) Next(sessionID string, machineID string, isPreviousSpecPassed bool) (string, error) {
	if err := svc.Repository.EndSpec(sessionID, machineID, isPreviousSpecPassed); err != nil {
		if err.Error() == datastore.ErrNoSuchEntity.Error() {
			return "", storage.ErrSessionNotFound
		}
	}

	specs, err := svc.Repository.GetSpecs(sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get specs: %s", err)
	}

	if len(specs) == 0 {
		return "", fmt.Errorf("backlog for session %s is empty", sessionID)
	}

	spec := svc.CalculateNext(specs)

	if spec.FilePath == "" {
		if err := svc.Repository.EndSession(sessionID); err != nil {
			return "", fmt.Errorf("failed to finish session: %s", err)
		}
		return "", ErrSessionFinished
	}

	if err := svc.Repository.StartSpec(sessionID, machineID, spec.ID); err != nil {
		return "", fmt.Errorf("failed to start spec: %s", err)
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
