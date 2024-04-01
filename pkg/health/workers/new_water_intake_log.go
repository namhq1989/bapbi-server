package workers

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
)

func (w Workers) NewWaterIntakeLog(bgCtx context.Context, t *asynq.Task) error {
	var (
		ctx = appcontext.New(bgCtx)
		log domain.WaterIntakeLog
	)

	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type(), "payload": string(t.Payload())})

	ctx.Logger().Info("unmarshal task payload", appcontext.Fields{})
	if err := json.Unmarshal(t.Payload(), &log); err != nil {
		ctx.Logger().Error("failed to unmarshal task payload", err, appcontext.Fields{})
		return err
	}

	// find drink water profile
	ctx.Logger().Text("find drink water profile")
	dwProfile, err := w.drinkWaterProfileRepository.FindDrinkWaterProfileByUserID(ctx, log.UserID)
	if err != nil {
		ctx.Logger().Error("failed to find drink water profile", err, appcontext.Fields{})
		return err
	}
	if dwProfile == nil || !dwProfile.IsEnabled {
		ctx.Logger().Error("drink water profile not found or disabled", nil, appcontext.Fields{})
		return apperrors.Health.DrinkWaterProfileNotFound
	}

	ctx.Logger().Text("find daily hydration report")
	startOfDay := manipulation.StartOfToday()
	report, err := w.dailyHydrationReportRepository.FindDailyHydrationReportByUserID(ctx, log.UserID, startOfDay)
	if err != nil {
		ctx.Logger().Error("failed to find daily hydration report", err, appcontext.Fields{})
		return err
	}

	if report == nil {
		// report not found, create new
		ctx.Logger().Text("today's report not found, create new domain model")
		report, err = domain.NewDailyHydrationReport(log.UserID, dwProfile.DailyIntakeAmount, log.Amount, startOfDay)
		if err != nil {
			ctx.Logger().Error("failed to create new domain model", err, appcontext.Fields{})
			return err
		}

		// persist to database
		ctx.Logger().Text("persist to database")
		if err = w.dailyHydrationReportRepository.CreateDailyHydrationReport(ctx, *report); err != nil {
			ctx.Logger().Error("failed to persist to database", err, appcontext.Fields{})
			return err
		}
	} else {
		// report found, update
		ctx.Logger().Text("today's report found, add intake amount")
		if err = report.AddIntakeAmount(log.Amount); err != nil {
			ctx.Logger().Error("failed to add intake amount", err, appcontext.Fields{})
		}

		// update
		ctx.Logger().Text("update to database")
		if err = w.dailyHydrationReportRepository.UpdateDailyHydrationReport(ctx, *report); err != nil {
			ctx.Logger().Error("failed to update domain model", err, appcontext.Fields{})
			return err
		}
	}

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}
