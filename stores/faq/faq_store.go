package faq

import (
	"bankapi/constants"
	"bankapi/models"
	"bankapi/utils"
	"context"
	"strings"
	"time"

	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type FAQStore interface {
	GetFaqList(ctx context.Context, authValues *models.AuthValues, platformType string) (interface{}, error)
}

type faqStore struct {
	logger *commonSrv.LoggerService
}

func NewFAQStore(logger *commonSrv.LoggerService) FAQStore {
	return &faqStore{
		logger: logger,
	}
}

func (f *faqStore) GetFaqList(ctx context.Context, authValues *models.AuthValues, platformType string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.FAQ,
		RequestURI: "/api/faq/list",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	var appOnly, webOnly bool

	if strings.ToLower(platformType) == "ios" || strings.ToLower(platformType) == "android" {
		appOnly = true
	} else if strings.ToLower(platformType) == "web" {
		webOnly = true
	} else {
		appOnly = true
	}

	faqData, err := models.GetAllCategoriesWithFaqs(appOnly, webOnly)
	if err != nil {
		logData.Message = "GetFaqList: Error getting faq list"
		f.logger.LogError(logData)
		return nil, err
	}

	logData.Message = "GetFaqList: Faq list found"
	logData.EndTime = time.Now()
	f.logger.LogInfo(logData)

	return faqData, nil
}
