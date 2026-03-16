package middleware

import (
	"application-for-kubernetes/internal/domain"

	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	cache  domain.Cache
	logger *logrus.Logger
}

func NewHTTPMiddlewares(cache domain.Cache, logger *logrus.Logger) *Middlewares {
	return &Middlewares{
		cache:  cache,
		logger: logger,
	}
}
