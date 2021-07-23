package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Shelex/split-specs/api/factory"
	"github.com/Shelex/split-specs/api/graph/generated"
	"github.com/Shelex/split-specs/api/graph/model"
	"github.com/Shelex/split-specs/internal/auth"
	"github.com/Shelex/split-specs/internal/users"
	"github.com/Shelex/split-specs/pkg/jwt"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (r *mutationResolver) AddSession(ctx context.Context, session model.SessionInput) (*model.SessionInfo, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, &users.AccessDeniedError{}
	}

	id, _ := gonanoid.New()

	specs := factory.SpecFilesToSpecs(session.SpecFiles)

	if err := r.SplitService.AddSession(user.ID, session.ProjectName, id, specs); err != nil {
		return nil, err
	}

	return &model.SessionInfo{
		SessionID:   id,
		ProjectName: session.ProjectName,
	}, nil
}

func (r *mutationResolver) Register(ctx context.Context, input model.User) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", &users.InvalidEmailOrPassordError{}
	}

	id, _ := gonanoid.New()

	user := users.User{
		Password: input.Password,
		Email:    input.Email,
		ID:       id,
	}

	if !user.EmailIsValid() {
		return "", &users.InvalidEmailFormat{}
	}

	if user.Exist() {
		return "", &users.WrongEmailOrPasswordError{}
	}

	if err := user.Create(); err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.User) (string, error) {
	user := users.User{
		Email:    input.Email,
		Password: input.Password,
	}

	correct := user.Authenticate()
	if !correct {
		return "", &users.WrongEmailOrPasswordError{}
	}

	dbUser, err := r.SplitService.Repository.GetUserByEmail(input.Email)
	if err != nil {
		return "", err
	}

	user.ID = dbUser.ID

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input model.ChangePasswordInput) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", &users.AccessDeniedError{}
	}

	if err := user.ChangePassword(input.Password, input.NewPassword); err != nil {
		return "", err
	}

	return "password changed", nil
}

func (r *mutationResolver) ShareProject(ctx context.Context, email string, projectName string) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", &users.AccessDeniedError{}
	}
	if err := r.SplitService.InviteUserToProject(users.UserToEntityUser(*user), email, projectName); err != nil {
		return "", err
	}
	return fmt.Sprintf("shared project %s with %s", projectName, email), nil
}

func (r *queryResolver) NextSpec(ctx context.Context, sessionID string, options *model.NextOptions) (string, error) {
	if user := auth.ForContext(ctx); user == nil {
		return "", &users.AccessDeniedError{}
	}
	machine := "default"
	if options != nil && options.MachineID != nil {
		machine = *options.MachineID
	}

	previousSpecPassed := false
	if options != nil && options.PreviousPassed != nil {
		previousSpecPassed = *options.PreviousPassed
	}

	next, err := r.SplitService.Next(sessionID, machine, previousSpecPassed)
	if err != nil {
		return "", fmt.Errorf("failed to receive next spec: %s", err)
	}
	return next, nil
}

func (r *queryResolver) Project(ctx context.Context, name string) (*model.Project, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, &users.AccessDeniedError{}
	}

	projectID, err := r.SplitService.Repository.GetUserProjectIDByName(user.ID, name)
	if err != nil {
		return nil, err
	}

	project, err := r.SplitService.Repository.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}

	sessions := make([]*model.Session, len(project.SessionIDs))

	for index, sessionID := range project.SessionIDs {
		session, err := r.SplitService.Repository.GetSession(sessionID)
		if err != nil {
			return nil, err
		}

		specs, err := r.SplitService.Repository.GetSpecs(sessionID, session.SpecIDs)
		if err != nil {
			return nil, err
		}

		sessions[index] = factory.ProjectSessionToApiSession(session, specs)
	}

	return &model.Project{
		ProjectName:   name,
		LatestSession: &project.LatestSession,
		Sessions:      sessions,
	}, nil
}

func (r *queryResolver) Projects(ctx context.Context) ([]string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, &users.AccessDeniedError{}
	}
	return r.SplitService.GetProjectList(users.UserToEntityUser(*user))
}

func (r *queryResolver) Session(ctx context.Context, sessionID string) (*model.Session, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, &users.AccessDeniedError{}
	}
	session, err := r.SplitService.Repository.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	specs, err := r.SplitService.Repository.GetSpecs(sessionID, session.SpecIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to receive specs for session %s: %s", sessionID, err)
	}
	return factory.ProjectSessionToApiSession(session, specs), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
