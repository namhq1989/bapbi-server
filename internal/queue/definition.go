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
	FeaturedWord string
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
		FeaturedWord: "language:term.featuredWord",
	},
}

type User struct {
	ID string `json:"id"`
}
