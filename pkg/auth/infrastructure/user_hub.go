package infrastructure

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
)

type UserHub struct {
	client userpb.UserServiceClient
}

func NewUserHub(client userpb.UserServiceClient) *UserHub {
	return &UserHub{
		client: client,
	}
}

func (r UserHub) GetOneByEmail(ctx *appcontext.AppContext, email string) (*domain.User, error) {
	resp, err := r.client.GetUserByEmail(ctx.Context(), &userpb.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		ctx.Logger().Error("failed to get user by email", err, appcontext.Fields{"email": email})
		return nil, err
	}

	user := resp.GetUser()
	if user == nil {
		return nil, nil
	}
	return &domain.User{
		ID:    user.GetId(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
	}, nil
}

func (r UserHub) GetOneByID(ctx *appcontext.AppContext, id string) (*domain.User, error) {
	resp, err := r.client.GetUserByID(ctx.Context(), &userpb.GetUserByIDRequest{
		Id: id,
	})
	if err != nil {
		ctx.Logger().Error("failed to get user by id", err, appcontext.Fields{"id": id})
		return nil, err
	}

	user := resp.GetUser()
	if user == nil {
		return nil, nil
	}
	return &domain.User{
		ID:    user.GetId(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
	}, nil
}

func (r UserHub) CreateUser(ctx *appcontext.AppContext, user domain.User) (string, error) {
	resp, err := r.client.CreateUser(ctx.Context(), &userpb.CreateUserRequest{
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		ctx.Logger().Error("failed to create user", err, appcontext.Fields{"user": user})
		return "", err
	}

	return resp.GetId(), nil
}
