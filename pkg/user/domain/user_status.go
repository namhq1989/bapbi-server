package domain

type UserStatus string

const (
	UserStatusUnknown  UserStatus = ""
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive" // for users have no activities in the last 1 month
	UserStatusClosed   UserStatus = "closed"
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusClosed:
		return string(s)
	default:
		return ""
	}
}

func (s UserStatus) IsValid() bool {
	return s != UserStatusUnknown
}

func ToUserStatus(value string) UserStatus {
	switch value {
	case UserStatusActive.String():
		return UserStatusActive
	case UserStatusInactive.String():
		return UserStatusInactive
	case UserStatusClosed.String():
		return UserStatusClosed
	default:
		return UserStatusUnknown
	}
}
