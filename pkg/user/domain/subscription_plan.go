package domain

import "strings"

type SubscriptionPlan string

const (
	SubscriptionPlanUnknown   SubscriptionPlan = ""
	SubscriptionPlanFree      SubscriptionPlan = "free"
	SubscriptionPlanSupporter SubscriptionPlan = "supporter"
)

func ToSubscriptionPlan(value string) SubscriptionPlan {
	switch strings.ToLower(value) {
	case SubscriptionPlanFree.String():
		return SubscriptionPlanFree
	case SubscriptionPlanSupporter.String():
		return SubscriptionPlanSupporter
	default:
		return SubscriptionPlanUnknown
	}
}

func (s SubscriptionPlan) String() string {
	switch s {
	case SubscriptionPlanFree, SubscriptionPlanSupporter:
		return string(s)
	default:
		return ""
	}
}

func (s SubscriptionPlan) IsValid() bool {
	return s != SubscriptionPlanUnknown
}
