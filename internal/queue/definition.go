package queue

type TypeHealthNames struct {
	UserCreated       string
	NewWaterIntakeLog string
}

type TypeUserNames struct {
	UserCreated string
	UserUpdated string
}

type LanguageNames struct {
	GenerateFeaturedWord     string
	GenerateWritingExercises string
}

var TypeNames = struct {
	Health   TypeHealthNames
	User     TypeUserNames
	Language LanguageNames
}{
	Health: TypeHealthNames{
		UserCreated:       "health:user.created",
		NewWaterIntakeLog: "health:hydration.newWaterIntakeLog",
	},
	User: TypeUserNames{
		UserCreated: "user:user.created",
		UserUpdated: "user:user.updated",
	},
	Language: LanguageNames{
		GenerateFeaturedWord:     "language:term.generateFeaturedWord",
		GenerateWritingExercises: "language:exercise.generateWritingExercises",
	},
}

type User struct {
	ID string `json:"id"`
}
