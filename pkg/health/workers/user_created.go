package workers

import (
	"context"

	"github.com/namhq1989/bapbi-server/internal/queue"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
)

func (w Workers) UserCreated(bgCtx context.Context, t *asynq.Task) error {
	var (
		ctx  = appcontext.New(bgCtx)
		user queue.User
	)

	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type(), "payload": string(t.Payload())})

	ctx.Logger().Info("unmarshal task payload", appcontext.Fields{})
	if err := json.Unmarshal(t.Payload(), &user); err != nil {
		ctx.Logger().Error("failed to unmarshal task payload", err, appcontext.Fields{})
		return err
	}

	// create health profile
	healthProfile, err := w.createDefaultHealthProfile(ctx, user.ID)
	if err != nil {
		return err
	}

	if err = w.createDefaultHydrationProfile(ctx, user.ID, healthProfile); err != nil {
		return err
	}

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}

func (w Workers) createDefaultHealthProfile(ctx *appcontext.AppContext, userID string) (*domain.HealthProfile, error) {
	ctx.Logger().Text("create default health profile")

	// new profile from domain
	ctx.Logger().Text("new profile from domain")
	domainProfile, err := domain.NewHealthProfile(userID, 70, 170, 7, 22)
	if err != nil {
		ctx.Logger().Error("failed to create new domain profile", err, appcontext.Fields{})
		return nil, err
	}

	// find profile in db, return error if existed
	ctx.Logger().Text("find profile in db")
	dbProfile, err := w.healthProfileRepository.FindHealthProfileByUserID(ctx, userID)
	if err != nil {
		ctx.Logger().Error("failed to find profile in db", err, appcontext.Fields{})
		return nil, err
	} else if dbProfile != nil {
		ctx.Logger().Error("profile already existed", nil, appcontext.Fields{})
		return nil, apperrors.Common.AlreadyExisted
	}

	// insert to db
	ctx.Logger().Text("insert profile to db")
	if err = w.healthProfileRepository.CreateHealthProfile(ctx, *domainProfile); err != nil {
		ctx.Logger().Error("failed to insert profile to db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("created default health profile successfully")
	return domainProfile, nil
}

func (w Workers) createDefaultHydrationProfile(ctx *appcontext.AppContext, userID string, healthProfile *domain.HealthProfile) error {
	ctx.Logger().Text("create default hydration profile")

	// find hydration profile in db
	ctx.Logger().Text("find hydration profile in db")
	hydrationProfile, err := w.hydrationProfileRepository.FindHydrationProfileByUserID(ctx, userID)
	if err != nil {
		ctx.Logger().Error("failed to find hydration profile in db", err, appcontext.Fields{})
		return err
	}

	if hydrationProfile != nil {
		ctx.Logger().Text("hydration profile already created")
		return nil
	}

	// get daily and hourly intake amount from health profile
	dailyIntakeAmount, hourlyIntakeAmount := healthProfile.GetDailyAndHourlyIntakeAmount()

	ctx.Logger().Text("create new")
	hydrationProfile, err = domain.NewHydrationProfile(userID, dailyIntakeAmount, hourlyIntakeAmount)
	if err != nil {
		ctx.Logger().Error("failed to create new hydration profile", err, appcontext.Fields{})
		return err
	}

	// save
	ctx.Logger().Text("persist hydration profile to database")
	if err = w.hydrationProfileRepository.CreateHydrationProfile(ctx, *hydrationProfile); err != nil {
		ctx.Logger().Error("failed to persist hydration profile to database", err, appcontext.Fields{})
		return err
	}

	return nil
}
