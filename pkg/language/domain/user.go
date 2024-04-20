package domain

import (
	"strings"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

type UserHub interface {
	GetUserPlan(ctx *appcontext.AppContext, userID string) (*SubscriptionPlan, error)
}

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

func (s SubscriptionPlan) IsFree() bool {
	return s == SubscriptionPlanFree
}

func (s SubscriptionPlan) IsSupporter() bool {
	return s == SubscriptionPlanSupporter
}

func (s SubscriptionPlan) IsExceededAddTermLimitation(todayAdded int64) bool {
	if s.IsFree() {
		return todayAdded >= 10
	} else if s.IsSupporter() {
		return todayAdded >= 30
	}

	return true
}

func (s SubscriptionPlan) IsExceededSearchLimitation(todaySearched int64) bool {
	if s.IsFree() {
		return todaySearched >= 20
	} else if s.IsSupporter() {
		return todaySearched >= 50
	}

	return true
}
