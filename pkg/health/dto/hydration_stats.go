package dto

import "github.com/namhq1989/bapbi-server/internal/utils/httprespond"

type HydrationStatsRequest struct{}

type HydrationStatsResponse struct {
	TodayIntakes             []HydrationStatsTodayIntake `json:"todayIntakes"`
	TodayProgress            HydrationStatsTodayProgress `json:"todayProgress"`
	LongestStreakValue       int                         `json:"longestStreakValue"`
	LongestStreakAt          *httprespond.TimeResponse   `json:"longestStreakAt"`
	HighestIntakeAmountValue int                         `json:"highestIntakeAmountValue"`
	HighestIntakeAmountAt    *httprespond.TimeResponse   `json:"highestIntakeAmountAt"`
}

type HydrationStatsTodayIntake struct {
	Amount    int                       `json:"amount"`
	IntakeAt  *httprespond.TimeResponse `json:"intakeAt"`
	CreatedAt *httprespond.TimeResponse `json:"createdAt"`
}

type HydrationStatsTodayProgress struct {
	Goal       int  `json:"goal"`
	Completed  int  `json:"completed"`
	IsAchieved bool `json:"isAchieved"`
}
