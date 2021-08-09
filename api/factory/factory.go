package factory

import (
	"github.com/Shelex/split-specs/api/graph/model"
	"github.com/Shelex/split-specs/entities"
)

func SpecFilesToSpecs(files []*model.SpecFile) []entities.Spec {
	specs := make([]entities.Spec, len(files))

	if len(files) > 0 {
		for i, f := range files {
			specs[i] = entities.Spec{
				FilePath: f.FilePath,
				Tests:    f.Tests,
			}
		}
	}

	return specs
}

func ProjectSessionsToApiSessions(sessions []entities.SessionWithSpecs) []*model.Session {
	apiSessions := make([]*model.Session, len(sessions))
	for i, session := range sessions {
		apiSessions[i] = ProjectSessionToApiSession(session)
	}
	return apiSessions
}

func ProjectSessionToApiSession(session entities.SessionWithSpecs) *model.Session {
	return &model.Session{
		ID:      session.ID,
		Start:   int(session.Start),
		End:     int(session.End),
		Backlog: specsToApiSpecs(session.Specs),
	}
}

func specsToApiSpecs(specs []entities.Spec) []*model.Spec {
	apiSpecs := make([]*model.Spec, len(specs))

	for i, spec := range specs {
		apiSpecs[i] = specToApiSpec(spec)
	}
	return apiSpecs
}

func specToApiSpec(spec entities.Spec) *model.Spec {
	return &model.Spec{
		File:              spec.FilePath,
		EstimatedDuration: int(spec.EstimatedDuration),
		Start:             int(spec.Start),
		End:               int(spec.End),
		Passed:            spec.Passed,
		AssignedTo:        spec.AssignedTo,
	}
}
