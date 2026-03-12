package app

import (
	"application-for-kubernetes/internal/app/handler"
	"application-for-kubernetes/internal/app/middleware"
	"application-for-kubernetes/internal/config"
	"application-for-kubernetes/internal/models/core"
	"context"
	"net/http"
	"time"

	"github.com/ayayaakasvin/lightmux"
	"github.com/sirupsen/logrus"
)

type ServerApp struct {
	server *http.Server

	lmux *lightmux.LightMux

	repo core.Repository

	httpcfg *config.HTTPServer
	logger  *logrus.Logger

	serverContext context.Context
	cancel        context.CancelFunc
}

func NewServerApp(
	httpcfg *config.HTTPServer,
	logger *logrus.Logger,
	repo core.Repository,
) *ServerApp {
	ctx, cancel := context.WithCancel(context.Background())
	return &ServerApp{
		httpcfg: httpcfg,

		logger: logger,
		repo:   repo,

		serverContext: ctx,
		cancel:        cancel,
	}
}

func (s *ServerApp) Run() error {
	s.setupServer()

	s.setupLightMux()

	return s.startServer()
}

func (s *ServerApp) Shutdown() {
	s.logger.Info("Initiating graceful shutdown...")
	s.cancel()
	if s.server != nil {
		s.server.Close()
	}
}

func (s *ServerApp) startServer() error {
	s.logger.Infof("Server has been started on port: %s", s.httpcfg.Address)
	s.logger.Infof("Available handlers:\n")

	s.lmux.PrintMiddlewareInfo()
	s.lmux.PrintRoutes()

	go printServerStatus(s.serverContext, s.logger)

	// RunTLS can be run when server is hosted on domain, acts as seperate service of file storing, for my project, id chose to encapsulate servers under one docker-compose and make nginx-gateaway for my api like auth, file, user service
	// if err := s.lmux.RunTLS(s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile); err != nil {
	s.logger.Info("Running the server", s.httpcfg.Address)
	return s.lmux.RunContext(s.serverContext)
}

// setuping server by pointer, so we dont have to return any value
func (s *ServerApp) setupServer() {
	if s.server == nil {
		// s.logger.Warn("Server is nil, creating a new server pointer")
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

	mws := middleware.NewHTTPMiddlewares(s.logger)
	hndlrs := handler.NewHTTPHandlers(s.repo, s.logger)

	s.lmux.Use(mws.LoggerMiddleware)

	apiGroup := s.lmux.NewGroup("/api")

	apiGroup.NewRoute("/ping").Handle(http.MethodGet, hndlrs.PingHandler())
	apiGroup.NewRoute("/health").Handle(http.MethodGet, hndlrs.HealthHandler())

	recordsGroup := apiGroup.ContinueGroup("/records")
	rg := recordsGroup.NewRoute("")
	rg.Handle(http.MethodGet, hndlrs.GetAllRecordsHandler())
	rg.Handle(http.MethodPost, hndlrs.InsertRecordHandler())
	recordsGroup.NewRoute("/").Handle(http.MethodGet, hndlrs.GetRecordByIDHandler())

	s.logger.Info("LightMux has been set up")
}

func printServerStatus(ctx context.Context, log *logrus.Logger) {
	ticker := time.NewTicker(time.Minute * 1)

	for {
		select {
		case <-ticker.C:
			log.Info("Server is alive...")
		case <-ctx.Done():
			return
		}
	}
}
