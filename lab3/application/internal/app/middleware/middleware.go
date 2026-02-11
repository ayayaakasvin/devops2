package middleware

import (
	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	logger     *logrus.Logger
}

func NewHTTPMiddlewares(logger *logrus.Logger) *Middlewares {	
	return &Middlewares{
		logger:     logger,
	}
}
