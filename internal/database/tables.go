package database

var Tables = struct {
	User      string
	AuthToken string

	// health
	HealthProfile              string
	HealthHydrationProfile     string
	HealthWaterIntakeLog       string
	HealthHydrationDailyReport string

	// language
	LanguageTerm string
}{
	User:      "user.users",
	AuthToken: "auth.tokens",

	HealthProfile:              "health.profiles",
	HealthHydrationProfile:     "health.hydration.profiles",
	HealthWaterIntakeLog:       "health.hydration.waterIntakeLogs",
	HealthHydrationDailyReport: "health.hydration.dailyReports",

	LanguageTerm: "language.terms",
}
