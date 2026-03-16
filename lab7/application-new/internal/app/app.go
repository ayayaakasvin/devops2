package app

import (
	"application-for-kubernetes/internal/app/handler"
	"application-for-kubernetes/internal/app/middleware"
	"application-for-kubernetes/internal/config"
	"application-for-kubernetes/internal/domain"
	"context"
	"net/http"

	"github.com/ayayaakasvin/lightmux"
	"github.com/sirupsen/logrus"
)

type ServerApp struct {
	server *http.Server

	lmux *lightmux.LightMux

	cache domain.Cache

	httpcfg *config.HTTPServer
	logger  *logrus.Logger

	serverContext context.Context
	cancel        context.CancelFunc
}

func NewServerApp(
	httpcfg *config.HTTPServer,
	logger *logrus.Logger,
	cache domain.Cache,
) *ServerApp {
	ctx, cancel := context.WithCancel(context.Background())
	return &ServerApp{
		httpcfg: httpcfg,

		logger: logger,
		cache:   cache,

		serverContext: ctx,
		cancel:        cancel,
	}
}

func (s *ServerApp) Start(ctx context.Context) error {
	s.setupServer()
	s.setupLightMux()

	return func() error {
		s.logger.Infof("Server has been started on port: %s", s.httpcfg.Address)

		return s.lmux.Run(ctx)
	}()
}

func (s *ServerApp) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *ServerApp) Shutdown() {
	s.logger.Info("Initiating graceful shutdown...")
	s.cancel()
	if s.server != nil {
		s.server.Close()
	}
}

func (s *ServerApp) setupServer() {
	if s.server == nil {
		s.logger.Warn("Server is nil, creating a new server pointer")
		s.server = &http.Server{}
	}

	s.server.Addr = s.httpcfg.Address
	s.logger.Infof("Server address set to: %s", s.server.Addr)
	s.server.IdleTimeout = s.httpcfg.IdleTimeout
	s.server.ReadTimeout = s.httpcfg.Timeout
	s.server.WriteTimeout = s.httpcfg.Timeout

	s.logger.Info("Server has been set up")
}

func (s *ServerApp) setupLightMux() {
	s.lmux = lightmux.NewLightMux(s.server)

	mws := middleware.NewHTTPMiddlewares(s.cache, s.logger)
	hndlrs := handler.NewHTTPHandlers(s.cache, s.logger)

	s.lmux.Use(mws.RecoverMiddleware)
	s.lmux.Use(mws.LoggerMiddleware)
	s.lmux.Use(mws.CacheWriteMiddleware)

	apiGroup := s.lmux.NewGroup("/api")
	apiGroup.NewRoute("/redis/info").Handle(http.MethodGet, hndlrs.InfoHandler())
	apiGroup.NewRoute("/health").Handle(http.MethodGet, hndlrs.HealthHandler())

	rgGroup := apiGroup.ContinueGroup("/records")
	rgGroup.NewRoute("/").Handle(http.MethodGet, hndlrs.GetAllRecordsHandler())

	s.logger.Info("LightMux has been set up")
	s.logger.Infof("Available handlers:\n")
	s.lmux.PrintMiddlewareInfo()
	s.lmux.PrintRoutes()
}
