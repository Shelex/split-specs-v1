package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Shelex/split-specs/api/factory"
	"github.com/Shelex/split-specs/api/graph/generated"
	"github.com/Shelex/split-specs/api/graph/model"
	"github.com/Shelex/split-specs/internal/auth"
	"github.com/Shelex/split-specs/internal/users"
	"github.com/Shelex/split-specs/pkg/jwt"
	uuid "github.com/satori/go.uuid"
)

func (r *mutationResolver) AddSession(ctx context.Context, session model.SessionInput) (*model.SessionInfo, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, &users.AccessDeniedError{}
	}

	id := uuid.NewV4().String()

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
	if input.Username == "" || input.Password == "" {
		return "", &users.InvalidUsernameOrPassordError{}
	}

	user := users.User{
		Username: input.Username,
		Password: input.Password,
		ID:       uuid.NewV4().String(),
	}

	if user.Exist() {
		return "", &users.WrongUsernameOrPasswordError{}
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
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		return "", &users.WrongUsernameOrPasswordError{}
	}

	dbUser, err := r.SplitService.Repository.GetUserByUsername(input.Username)
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

func (r *mutationResolver) InviteUser(ctx context.Context, username string, projectName string) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", &users.AccessDeniedError{}
	}
	if err := r.SplitService.InviteUserToProject(users.UserToEntityUser(*user), username, projectName); err != nil {
		return "", err
	}
	return "invited " + username, nil
}

func (r *queryResolver) NextSpec(ctx context.Context, sessionID string, machineID *string) (string, error) {
	if user := auth.ForContext(ctx); user == nil {
		return "", &users.AccessDeniedError{}
	}
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
		sessions[index] = factory.ProjectSessionToApiSession(session)
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
