package hub

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type CreateUserHandler struct {
	userHub         domain.UserHub
	queueRepository domain.QueueRepository
}

func NewCreateUserHandler(queueRepository domain.QueueRepository, userHub domain.UserHub) CreateUserHandler {
	return CreateUserHandler{
		userHub:         userHub,
		queueRepository: queueRepository,
	}
}

func (h CreateUserHandler) CreateUser(ctx *appcontext.AppContext, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	ctx.Logger().Info("create new user", appcontext.Fields{"name": req.Name, "email": req.Email})

	// validate user
	ctx.Logger().Info("validate user", appcontext.Fields{})
	domainUser, err := domain.CreateUser(req.Name, req.Email)
	if err != nil {
		ctx.Logger().Error("failed to create domain's user", err, appcontext.Fields{})
		return nil, err
	}

	// find user with email
	ctx.Logger().Info("find user with email", appcontext.Fields{})
	user, err := h.userHub.FindOneByEmail(ctx, domainUser.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		ctx.Logger().Error("user already existed", nil, appcontext.Fields{})
		return nil, apperrors.Common.EmailAlreadyExisted
	}

	// create user
	ctx.Logger().Info("create user", appcontext.Fields{})
	user, err = h.userHub.CreateUser(ctx, *domainUser)
	if err != nil {
		return nil, err
	}

	// add queue
	ctx.Logger().Info("add queue task", appcontext.Fields{})
	err = h.queueRepository.EnqueueUserCreated(ctx, *user)
	if err != nil {
		return nil, err
	}

	// respond
	ctx.Logger().Info("done create user", appcontext.Fields{})
	return &userpb.CreateUserResponse{
		Id: user.ID,
	}, nil
}
