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
	LanguageTerm                string
	LanguageUserTerm            string
	LanguageUserActionHistory   string
	LanguageWritingExercise     string
	LanguageUserWritingExercise string
	LanguageUserTermExercise    string
}{
	User:      "user.users",
	AuthToken: "auth.tokens",

	HealthProfile:              "health.profiles",
	HealthHydrationProfile:     "health.hydration.profiles",
	HealthWaterIntakeLog:       "health.hydration.waterIntakeLogs",
	HealthHydrationDailyReport: "health.hydration.dailyReports",

	LanguageTerm:                "language.terms",
	LanguageUserTerm:            "language.userTerms",
	LanguageUserActionHistory:   "language.userActionHistories",
	LanguageWritingExercise:     "language.writingExercises",
	LanguageUserWritingExercise: "language.userWritingExercises",
	LanguageUserTermExercise:    "language.userTermExercises",
}
