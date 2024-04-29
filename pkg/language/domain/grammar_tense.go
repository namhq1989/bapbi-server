package domain

import "github.com/namhq1989/bapbi-server/internal/utils/manipulation"

type GrammarTenseCode string

const (
	GrammarTenseCodeUnknown           GrammarTenseCode = ""
	GrammarTenseCodePresentSimple     GrammarTenseCode = "present_simple"
	GrammarTenseCodePresentContinuous GrammarTenseCode = "present_continuous"
	GrammarTenseCodePresentPerfect    GrammarTenseCode = "present_perfect"
	GrammarTenseCodePastSimple        GrammarTenseCode = "past_simple"
	GrammarTenseCodePastContinuous    GrammarTenseCode = "past_continuous"
	GrammarTenseCodeFutureSimple      GrammarTenseCode = "future_simple"
)

func (c GrammarTenseCode) String() string {
	switch c {
	case GrammarTenseCodePresentSimple, GrammarTenseCodePresentContinuous, GrammarTenseCodePresentPerfect, GrammarTenseCodePastSimple, GrammarTenseCodePastContinuous, GrammarTenseCodeFutureSimple:
		return string(c)
	default:
		return ""
	}
}

func (c GrammarTenseCode) IsValid() bool {
	return c != GrammarTenseCodeUnknown
}

func RandomGrammarTenseCode() GrammarTenseCode {
	randIndex := manipulation.RandomIntInRange(0, len(EnglishGrammarTenses))
	return EnglishGrammarTenses[randIndex].Code
}

func ToGrammarTenseCode(value string) GrammarTenseCode {
	switch value {
	case GrammarTenseCodePresentSimple.String():
		return GrammarTenseCodePresentSimple
	case GrammarTenseCodePresentContinuous.String():
		return GrammarTenseCodePresentContinuous
	case GrammarTenseCodePresentPerfect.String():
		return GrammarTenseCodePresentPerfect
	case GrammarTenseCodePastSimple.String():
		return GrammarTenseCodePastSimple
	case GrammarTenseCodePastContinuous.String():
		return GrammarTenseCodePastContinuous
	case GrammarTenseCodeFutureSimple.String():
		return GrammarTenseCodeFutureSimple
	default:
		return GrammarTenseCodeUnknown
	}
}

type GrammarTense struct {
	Code        GrammarTenseCode
	Name        LanguageTranslation
	Formula     string
	UseCase     LanguageTranslation
	SignalWords []string
}

var EnglishGrammarTenses = []GrammarTense{
	{
		Code: GrammarTenseCodePresentSimple,
		Name: LanguageTranslation{
			English:    "Present Simple",
			Vietnamese: "Thì hiện tại đơn",
		},
		Formula: "S + V(s/es) + O",
		UseCase: LanguageTranslation{
			English:    "Regular habits, general truths, repeated actions, fixed arrangements",
			Vietnamese: "Thói quen thường xuyên, sự thật chung, các hành động lặp đi lặp lại, các sắp xếp cố định",
		},
		SignalWords: []string{"often", "always", "never", "usually", "sometimes", "rarely", "every day", "on Mondays"},
	},
	{
		Code: GrammarTenseCodePresentContinuous,
		Name: LanguageTranslation{
			English:    "Present Continuous",
			Vietnamese: "Thì hiện tại tiếp diễn",
		},
		Formula: "S + am/is/are + V-ing + O",
		UseCase: LanguageTranslation{
			English:    "Actions currently happening, temporary actions, arrangements in the near future",
			Vietnamese: "Các hành động đang diễn ra tại thời điểm nói, các hành động tạm thời, các sắp xếp trong tương lai gần",
		},
		SignalWords: []string{"now", "at the moment", "currently", "right now", "today", "these days"},
	},
	{
		Code: GrammarTenseCodePastSimple,
		Name: LanguageTranslation{
			English:    "Past Simple",
			Vietnamese: "Thì quá khứ đơn",
		},
		Formula: "S + V-ed/V2 + O",
		UseCase: LanguageTranslation{
			English:    "Completed actions in the past, a series of past actions, past habits",
			Vietnamese: "Các hành động hoàn thành trong quá khứ, một loạt các hành động quá khứ, thói quen quá khứ",
		},
		SignalWords: []string{"yesterday", "last week", "in 1998", "the other day", "when", "ago"},
	},
	{
		Code: GrammarTenseCodePastContinuous,
		Name: LanguageTranslation{
			English:    "Past Continuous",
			Vietnamese: "Thì quá khứ tiếp diễn",
		},
		Formula: "S + was/were + V-ing + O",
		UseCase: LanguageTranslation{
			English:    "Actions that were in progress at a specific time in the past, simultaneous past actions",
			Vietnamese: "Các hành động đang diễn ra tại một thời điểm cụ thể trong quá khứ, hai hành động xảy ra đồng thời",
		},
		SignalWords: []string{"while", "when", "as", "at this time yesterday"},
	},
	{
		Code: GrammarTenseCodeFutureSimple,
		Name: LanguageTranslation{
			English:    "Future Simple",
			Vietnamese: "Thì tương lai đơn",
		},
		Formula: "S + will + V + O",
		UseCase: LanguageTranslation{
			English:    "Future decisions at the moment of speaking, predictions, promises, spontaneous decisions",
			Vietnamese: "Các quyết định về tương lai tại thời điểm nói, dự đoán, lời hứa, quyết định tức thì",
		},
		SignalWords: []string{"tomorrow", "next week", "soon", "someday", "later", "in the future"},
	},
	{
		Code: GrammarTenseCodePresentPerfect,
		Name: LanguageTranslation{
			English:    "Present Perfect",
			Vietnamese: "Thì hiện tại hoàn thành",
		},
		Formula: "S + has/have + V3 + O",
		UseCase: LanguageTranslation{
			English:    "Actions that have occurred at an unspecified time before now, actions that began in the past and continue to the present, life experiences",
			Vietnamese: "Các hành động đã xảy ra tại một thời điểm không xác định trước bây giờ, các hành động bắt đầu trong quá khứ và tiếp tục đến hiện tại, kinh nghiệm sống",
		},
		SignalWords: []string{"already", "just", "yet", "ever", "never", "recently", "so far", "until now", "since", "for"},
	},
}
