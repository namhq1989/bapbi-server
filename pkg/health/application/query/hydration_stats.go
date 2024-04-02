package query

import (
	"fmt"
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type HydrationStatsHandler struct {
	hydrationProfileRepository     domain.HydrationProfileRepository
	hydrationDailyReportRepository domain.HydrationDailyReportRepository
	waterIntakeLogRepository       domain.WaterIntakeLogRepository
}

func NewHydrationStatsHandler(
	hydrationProfileRepository domain.HydrationProfileRepository,
	hydrationDailyReportRepository domain.HydrationDailyReportRepository,
	waterIntakeLogRepository domain.WaterIntakeLogRepository,
) HydrationStatsHandler {
	return HydrationStatsHandler{
		hydrationProfileRepository:     hydrationProfileRepository,
		hydrationDailyReportRepository: hydrationDailyReportRepository,
		waterIntakeLogRepository:       waterIntakeLogRepository,
	}
}

func (h HydrationStatsHandler) HydrationStats(ctx *appcontext.AppContext, performerID string, _ dto.HydrationStatsRequest) (*dto.HydrationStatsResponse, error) {
	ctx.Logger().Info("get hydration stats", appcontext.Fields{"performerID": performerID})

	// find hydration profile
	profile := h.getHydrationProfile(ctx, performerID)
	if profile == nil {
		return nil, apperrors.Health.HydrationProfileNotFound
	}

	var result = &dto.HydrationStatsResponse{
		TodayIntakes:             h.getWaterIntakeLogsByUserID(ctx, performerID),
		TodayProgress:            dto.HydrationStatsTodayProgress{},
		LongestStreakValue:       profile.LongestSuccessStreakValue,
		LongestStreakAt:          httprespond.NewTimeResponse(profile.LongestSuccessStreakAt),
		HighestIntakeAmountValue: profile.HighestIntakeAmountValue,
		HighestIntakeAmountAt:    httprespond.NewTimeResponse(profile.HighestIntakeAmountAt),
	}

	// find today report
	report := h.getTodayHydrationReport(ctx, performerID)
	if report != nil {
		result.TodayProgress.Goal = report.GoalAmount
		result.TodayProgress.Completed = report.IntakeAmount
		result.TodayProgress.IsAchieved = report.IsAchieved
	}

	return result, nil
}

func (h HydrationStatsHandler) getHydrationProfile(ctx *appcontext.AppContext, performerID string) *domain.HydrationProfile {
	ctx.Logger().Text("get hydration profile")

	profile, err := h.hydrationProfileRepository.FindHydrationProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find hydration profile", err, appcontext.Fields{})
		return nil
	}
	return profile
}

func (h HydrationStatsHandler) getTodayHydrationReport(ctx *appcontext.AppContext, performerID string) *domain.HydrationDailyReport {
	ctx.Logger().Text("get today hydration report")

	report, err := h.hydrationDailyReportRepository.FindHydrationDailyReportByUserID(ctx, performerID, manipulation.StartOfToday())
	if err != nil {
		ctx.Logger().Error("failed to find hydration report", err, appcontext.Fields{})
		return nil
	}
	return report
}

func (h HydrationStatsHandler) getWaterIntakeLogsByUserID(ctx *appcontext.AppContext, performerID string) []dto.HydrationStatsTodayIntake {
	ctx.Logger().Text("get today water intake logs")
	var (
		result = make([]dto.HydrationStatsTodayIntake, 0)
		filter = domain.WaterIntakeLogFilter{
			From: manipulation.StartOfToday(),
			To:   time.Now(),
		}
	)

	ctx.Logger().Text("find water intake logs in database")
	logs, err := h.waterIntakeLogRepository.FindWaterIntakeLogsByUserID(ctx, performerID, filter)
	if err != nil {
		ctx.Logger().Error("failed to find water intake logs", err, appcontext.Fields{})
		return result
	}

	if len(logs) == 0 {
		ctx.Logger().Text("no water intake logs found")
		return result
	}

	ctx.Logger().Text(fmt.Sprintf("found %d water intake logs", len(logs)))
	for _, log := range logs {
		result = append(result, dto.HydrationStatsTodayIntake{
			Amount:    log.Amount,
			IntakeAt:  httprespond.NewTimeResponse(log.IntakeAt),
			CreatedAt: httprespond.NewTimeResponse(log.CreatedAt),
		})
	}

	return result
}
