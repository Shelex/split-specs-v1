package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Shelex/split-specs/api/factory"
	"github.com/Shelex/split-specs/api/graph/generated"
	"github.com/Shelex/split-specs/api/graph/model"
	"github.com/Shelex/split-specs/entities"
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

func (r *mutationResolver) DeleteSession(ctx context.Context, sessionID string) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", &users.AccessDeniedError{}
	}

	if err := r.SplitService.Repository.DeleteSession(user.Email, sessionID); err != nil {
		return "", err
	}
	return "session deleted", nil
}

func (r *mutationResolver) DeleteProject(ctx context.Context, projectName string) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", &users.AccessDeniedError{}
	}

	projectID, err := r.SplitService.Repository.GetUserProjectIDByName(user.ID, projectName)
	if err != nil {
		return "", err
	}

	if err := r.SplitService.Repository.DeleteProject(user.Email, projectID); err != nil {
		return "", err
	}

	return "project deleted", nil
}

func (r *mutationResolver) AddAPIKey(ctx context.Context, name string, expireAt int) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", &users.AccessDeniedError{}
	}

	id, _ := gonanoid.New()

	apiKey := entities.ApiKey{
		ID:       id,
		UserID:   user.ID,
		Name:     name,
		ExpireAt: int64(expireAt),
	}

	if err := r.SplitService.Repository.CreateApiKey(user.ID, apiKey); err != nil {
		return "", err
	}

	token, err := jwt.GenerateApiKey(*user, apiKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) DeleteAPIKey(ctx context.Context, keyID string) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", &users.AccessDeniedError{}
	}

	err := r.SplitService.Repository.DeleteApiKey(user.ID, keyID)
	if err != nil {
		return "", err
	}

	return "apiKey deleted", nil
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

func (r *queryResolver) Project(ctx context.Context, name string, pagination *model.Pagination) (*model.Project, error) {
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

	sessions, total, err := r.SplitService.Repository.GetProjectSessions(projectID, factory.ApiPaginationToPagination(pagination))
	if err != nil {
		return nil, err
	}

	return &model.Project{
		ProjectName:   name,
		LatestSession: &project.LatestSession,
		Sessions:      factory.ProjectSessionsToApiSessions(sessions),
		TotalSessions: total,
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
	session, err := r.SplitService.Repository.GetSessionWithSpecs(sessionID)
	if err != nil {
		return nil, err
	}

	return factory.ProjectSessionToApiSession(session), nil
}

func (r *queryResolver) GetAPIKeys(ctx context.Context) ([]*model.APIKey, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, &users.AccessDeniedError{}
	}

	keys, err := r.SplitService.Repository.GetApiKeys(user.ID)
	if err != nil {
		return nil, err
	}

	return factory.ApiKeysToApi(keys), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
