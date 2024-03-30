package database

var Tables = struct {
	User      string
	AuthToken string

	// health
	HealthProfile           string
	HealthDrinkWaterProfile string
	HealthWaterIntakeLog    string
}{
	User:      "user.users",
	AuthToken: "auth.tokens",

	HealthProfile:           "health-profiles",
	HealthDrinkWaterProfile: "health-drinkWater-profiles",
	HealthWaterIntakeLog:    "health-drinkWater-waterIntakeLogs",
}
