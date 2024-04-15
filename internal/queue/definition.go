package queue

type TypeHealthNames struct {
	UserCreated       string
	NewWaterIntakeLog string
}

type TypeUserNames struct {
	UserCreated string
	UserUpdated string
}

var TypeNames = struct {
	Health TypeHealthNames
	User   TypeUserNames
}{
	Health: TypeHealthNames{
		UserCreated:       "health:user.created",
		NewWaterIntakeLog: "health:hydration.newWaterIntakeLog",
	},
	User: TypeUserNames{
		UserCreated: "user:user.created",
		UserUpdated: "user:user.updated",
	},
}

type User struct {
	ID string `json:"id"`
}
