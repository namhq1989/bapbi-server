package dto

import "github.com/namhq1989/bapbi-server/pkg/language/domain"

type GetDataRequest struct{}

type GetDataResponse struct {
	GrammarTenses []GrammarTense `json:"grammarTenses"`
}

type GrammarTense struct {
	Code        string              `json:"code"`
	Name        LanguageTranslation `json:"name"`
	Formula     string              `json:"formula"`
	UseCase     LanguageTranslation `json:"useCase"`
	SignalWords []string            `json:"signalWords"`
}

type LanguageTranslation struct {
	English    string `json:"english"`
	Vietnamese string `json:"vietnamese"`
}

func (d GetDataResponse) FromDomain(tenses []domain.GrammarTense) []GrammarTense {
	result := make([]GrammarTense, len(tenses))
	for i, t := range tenses {
		result[i] = GrammarTense{
			Code: t.Code,
			Name: LanguageTranslation{
				English:    t.Name.English,
				Vietnamese: t.Name.Vietnamese,
			},
			Formula: t.Formula,
			UseCase: LanguageTranslation{
				English:    t.UseCase.English,
				Vietnamese: t.UseCase.Vietnamese,
			},
			SignalWords: t.SignalWords,
		}
	}
	return result
}
