package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shelex/split-specs/api/factory"
	"github.com/Shelex/split-specs/api/graph/generated"
	"github.com/Shelex/split-specs/api/graph/model"
	uuid "github.com/satori/go.uuid"
)

func (r *mutationResolver) AddSession(ctx context.Context, session model.SessionInput) (*model.SessionInfo, error) {
	jsonized, _ := json.Marshal(session)
	log.Println(string(jsonized))

	id := uuid.NewV4().String()

	specs := factory.SpecFilesToSpecs(session.SpecFiles)

	if err := r.SplitService.AddSession(session.ProjectName, id, specs); err != nil {
		return nil, err
	}

	return &model.SessionInfo{
		SessionID:   id,
		ProjectName: session.ProjectName,
	}, nil
}

func (r *queryResolver) NextSpec(ctx context.Context, sessionID string, machineID *string) (string, error) {
	machine := "default"
	if machineID != nil {
		machine = *machineID
	}

	next, err := r.SplitService.Next(sessionID, machine)
	if err != nil {
		return "", err
	}
	return next, nil
}

func (r *queryResolver) Project(ctx context.Context, name string) (*model.Project, error) {
	project, err := r.SplitService.GetProject(name)
	if err != nil {
		return nil, err
	}
	return &model.Project{
		ProjectName:   name,
		LatestSession: &project.LatestSession,
		Sessions:      factory.ProjectSessionsToApiSessions(project.Sessions),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
