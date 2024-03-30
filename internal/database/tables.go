package database

var Tables = struct {
	User      string
	AuthToken string
}{
	User:      "users",
	AuthToken: "auth-tokens",
}
