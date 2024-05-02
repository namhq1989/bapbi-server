package service

import "github.com/namhq1989/bapbi-server/pkg/language/domain"

type LanguageService struct {
	userTermRepository          domain.UserTermRepository
	userActionHistoryRepository domain.UserActionHistoryRepository
	userHub                     domain.UserHub
}

func NewLanguageService(userTermRepository domain.UserTermRepository, userActionHistoryRepository domain.UserActionHistoryRepository, userHub domain.UserHub) LanguageService {
	return LanguageService{
		userTermRepository:          userTermRepository,
		userActionHistoryRepository: userActionHistoryRepository,
		userHub:                     userHub,
	}
}
