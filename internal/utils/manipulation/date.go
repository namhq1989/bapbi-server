package manipulation

import "time"

func StartOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func StartOfYesterday() time.Time {
	return StartOfToday().AddDate(0, 0, -1)
}

func IsToday(t time.Time) bool {
	today := StartOfToday()
	return t.Year() == today.Year() && t.Month() == today.Month() && t.Day() == today.Day()
}

func IsYesterday(t time.Time) bool {
	yesterday := StartOfYesterday()
	return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}
