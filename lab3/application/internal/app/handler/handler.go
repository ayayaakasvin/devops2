package handler

import (
	"application-for-kubernetes/internal/models/core"

	"github.com/sirupsen/logrus"
)

type Handlers struct {
	repo  core.Repository
	logger *logrus.Logger
}

func NewHTTPHandlers(repo core.Repository, logger *logrus.Logger) *Handlers {
	return &Handlers{
		repo:  repo,

		logger: logger,
	}
}
