package appcontext

import (
	"context"

	"github.com/google/uuid"
	"github.com/namhq1989/bapbi-server/internal/utils/logger"
)

type contextKey int

const (
	userContextKey contextKey = iota
	sourceContextKey
	ipContextKey
)

type AppContext struct {
	requestID string
	traceID   string
	logger    *logger.Logger
	context   context.Context
}

type Fields = logger.Fields

func New(ctx context.Context) *AppContext {
	var (
		requestID = generateID()
		traceID   = generateID()
	)

	return &AppContext{
		requestID: requestID,
		traceID:   traceID,
		logger:    logger.NewLogger(logger.Fields{"requestId": requestID, "traceId": traceID}),
		context:   ctx,
	}
}

func (appCtx *AppContext) AddLogData(fields Fields) {
	appCtx.logger.AddData(fields)
}

func (appCtx *AppContext) Logger() *logger.Logger {
	return appCtx.logger

}

func (appCtx *AppContext) Context() context.Context {
	return appCtx.context
}

func (appCtx *AppContext) SetContext(ctx context.Context) {
	appCtx.context = ctx
}

func (appCtx *AppContext) SetUserID(id string) {
	appCtx.context = context.WithValue(appCtx.context, userContextKey, id)
}

func (appCtx *AppContext) GetUserID() string {
	id, ok := appCtx.context.Value(userContextKey).(string)
	if !ok {
		return ""
	}
	return id
}

func (appCtx *AppContext) SetIP(ip string) {
	appCtx.context = context.WithValue(appCtx.context, ipContextKey, ip)
}

func (appCtx *AppContext) GetIP() string {
	ip, ok := appCtx.context.Value(ipContextKey).(string)
	if !ok {
		return ""
	}
	return ip
}

func (appCtx *AppContext) SetSourceRest() {
	appCtx.context = context.WithValue(appCtx.context, sourceContextKey, "rest")
}

func (appCtx *AppContext) IsSourceRest() bool {
	source, ok := appCtx.context.Value(sourceContextKey).(string)
	if !ok {
		return false
	}
	return source == "rest"
}

func generateID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
