package api

import (
	"context"
	"github.com/duclmse/fengine/fengine"
	log "github.com/duclmse/fengine/pkg/logger"
)

var _ fengine.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    fengine.Service
}

func LoggingMiddleware(svc fengine.Service, logger log.Logger) fengine.Service {
	return &loggingMiddleware{logger: logger, svc: svc}
}

func (l loggingMiddleware) Get(ctx context.Context, id string) (interface{}, error) {
	l.logger.Debug_("")
	return l.svc.Get(ctx, id)
}
