package handler

import (
	"application-for-kubernetes/internal/domain"

	"github.com/sirupsen/logrus"
)

type Handlers struct {
	cache  domain.Cache
	logger *logrus.Logger
}

func NewHTTPHandlers(cache domain.Cache, logger *logrus.Logger) *Handlers {
	return &Handlers{
		cache:  cache,

		logger: logger,
	}
}
