package workers

import (
	"context"

	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
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

	// find hydration profile
	ctx.Logger().Text("find hydration profile")
	hydrationProfile, err := w.hydrationProfileRepository.FindHydrationProfileByUserID(ctx, log.UserID)
	if err != nil {
		ctx.Logger().Error("failed to find hydration profile", err, appcontext.Fields{})
		return err
	}
	if hydrationProfile == nil || !hydrationProfile.IsEnabled {
		ctx.Logger().Error("hydration profile not found or disabled", nil, appcontext.Fields{})
		return apperrors.Health.HydrationProfileNotFound
	}

	// TASK

	ctx.Logger().Text("update hydration daily report")
	report, err := w.updateHydrationDailyReport(ctx, log, *hydrationProfile)
	if err != nil {
		ctx.Logger().Error("failed to update hydration daily report", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("update hydration profile")
	if err = w.updateHydrationProfile(ctx, log, *hydrationProfile, *report); err != nil {
		ctx.Logger().Error("failed to update hydration profile", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}

func (w Workers) updateHydrationProfile(ctx *appcontext.AppContext, log domain.WaterIntakeLog, profile domain.HydrationProfile, report domain.HydrationDailyReport) error {
	ctx.Logger().Text("set highest intake amount")
	if err := profile.SetHighestIntakeAmount(log.Amount); err != nil {
		ctx.Logger().Error("failed to set highest intake amount", err, appcontext.Fields{})
		return err
	}

	// if today's goal is achieved, update streak
	if report.IsAchieved {
		ctx.Logger().Text("today's goal is achieved, update hydration streak")
		if manipulation.IsToday(profile.CurrentStreakDate) {
			ctx.Logger().Text("latest streak is today's streak, do nothing")
			return nil
		} else if manipulation.IsYesterday(profile.CurrentStreakDate) {
			ctx.Logger().Text("latest streak is yesterday's streak, increase streak")
			profile.IncreaseStreak()
		} else {
			ctx.Logger().Text("latest streak is far from today, reset streak")
			profile.ResetStreak()
		}
	} else {
		ctx.Logger().Text("today's goal is not achieved, skip calculate the streak")
	}

	return w.hydrationProfileRepository.UpdateHydrationProfile(ctx, profile)
}

func (w Workers) updateHydrationDailyReport(ctx *appcontext.AppContext, log domain.WaterIntakeLog, hydrationProfile domain.HydrationProfile) (*domain.HydrationDailyReport, error) {
	ctx.Logger().Text("find hydration daily report")

	startOfDay := manipulation.StartOfToday()
	report, err := w.hydrationDailyReportRepository.FindHydrationDailyReportByUserID(ctx, log.UserID, startOfDay)
	if err != nil {
		ctx.Logger().Error("failed to find hydration daily report", err, appcontext.Fields{})
		return nil, err
	}

	if report == nil {
		// report not found, create new
		ctx.Logger().Text("today's report not found, create new domain model")
		report, err = domain.NewHydrationDailyReport(log.UserID, hydrationProfile.DailyIntakeAmount, log.Amount, startOfDay)
		if err != nil {
			ctx.Logger().Error("failed to create new domain model", err, appcontext.Fields{})
			return nil, err
		}

		// persist to database
		ctx.Logger().Text("persist hydration daily report to database")
		if err = w.hydrationDailyReportRepository.CreateHydrationDailyReport(ctx, *report); err != nil {
			ctx.Logger().Error("failed to persist hydration daily report to database", err, appcontext.Fields{})
			return nil, err
		}
	} else {
		// report found, update
		ctx.Logger().Text("today's report found, add intake amount")
		if err = report.AddIntakeAmount(log.Amount); err != nil {
			ctx.Logger().Error("failed to add intake amount", err, appcontext.Fields{})
		}

		// update
		ctx.Logger().Text("update to database")
		if err = w.hydrationDailyReportRepository.UpdateHydrationDailyReport(ctx, *report); err != nil {
			ctx.Logger().Error("failed to update domain model", err, appcontext.Fields{})
			return nil, err
		}
	}

	return report, nil
}
