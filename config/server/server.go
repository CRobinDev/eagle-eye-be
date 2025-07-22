package server

import (
	"context"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/config/gemini"
	"github.com/CRobinDev/karsa/config/midtrans"
	"github.com/CRobinDev/karsa/config/router"
	authHandler "github.com/CRobinDev/karsa/internal/app/auth/handler"
	authRepository "github.com/CRobinDev/karsa/internal/app/auth/repository"
	authService "github.com/CRobinDev/karsa/internal/app/auth/service"
	customerHandler "github.com/CRobinDev/karsa/internal/app/customer/handler"
	customerRepository "github.com/CRobinDev/karsa/internal/app/customer/repository"
	customerService "github.com/CRobinDev/karsa/internal/app/customer/service"
	"github.com/CRobinDev/karsa/internal/app/detection/handler"
	"github.com/CRobinDev/karsa/internal/app/detection/repository"
	"github.com/CRobinDev/karsa/internal/app/detection/service"
	paymentHandler "github.com/CRobinDev/karsa/internal/app/payment/handler"
	paymentRepository "github.com/CRobinDev/karsa/internal/app/payment/repository"
	paymentService "github.com/CRobinDev/karsa/internal/app/payment/service"
	userHandler "github.com/CRobinDev/karsa/internal/app/user/handler"
	userRepository "github.com/CRobinDev/karsa/internal/app/user/repository"
	userService "github.com/CRobinDev/karsa/internal/app/user/service"
	"github.com/CRobinDev/karsa/internal/infra/database"
	_gemini "github.com/CRobinDev/karsa/pkg/gemini"
	"github.com/CRobinDev/karsa/pkg/jwt"
	"github.com/CRobinDev/karsa/pkg/log"
	_midtrans "github.com/CRobinDev/karsa/pkg/midtrans"
	"github.com/CRobinDev/karsa/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"google.golang.org/genai"
)

type Handler interface {
	SetEndpoint(router fiber.Router)
}

type ServerConfig struct {
	router   *fiber.App
	handlers []Handler
	logger   *logrus.Logger
	db       *sqlx.DB
	val      validator.Validator
	midtrans *snap.Client
	gemini   *genai.Client
	// cache    *bigcache.BigCache
	// oauth *google.GoogleOAuth
}

func Init() *ServerConfig {
	env.NewEnv()
	jwt.InitJWT()
	logger := log.NewLogger()
	app := router.NewFiber(logger)
	db := database.NewPostgresPool(logger)
	val := validator.NewValidator()
	midtrans := midtrans.NewMidtrans()
	gemini := gemini.NewGemini()

	return &ServerConfig{
		router:   app,
		logger:   logger,
		db:       db,
		val:      val,
		midtrans: midtrans,
		gemini:   gemini,
	}
}

func (s *ServerConfig) RegisterHandler(handler Handler) {
	s.handlers = append(s.handlers, handler)
}

func (s *ServerConfig) DependencyInjection() {
	midtransService := _midtrans.NewMidtransService(s.midtrans)
	geminiService := _gemini.NewGeminiService(s.gemini, s.logger)

	userRepostory := userRepository.NewUserRepository(s.db)
	userService := userService.NewUserService(userRepostory, s.logger)
	userHandler := userHandler.NewUserHandler(userService, s.val)

	authRepository := authRepository.NewAuthRepository(s.db)
	authService := authService.NewAuthService(authRepository, s.logger)
	authHandler := authHandler.NewAuthHandler(authService, s.val)

	paymentRepository := paymentRepository.NewPaymentRepository(s.db)
	paymentService := paymentService.NewPaymentService(paymentRepository, midtransService, s.logger)
	paymentHandler := paymentHandler.NewPaymentHandler(paymentService, s.val)

	customerRepository := customerRepository.NewCustomerRepository(s.db)
	customerService := customerService.NewCustomerService(customerRepository, userRepostory, paymentService, s.logger)
	customerHandler := customerHandler.NewCustomerHandler(customerService, s.val)

	detectionRepository := repository.NewDetectionRepository(s.db)
	detectionService := service.NewDetectionService(detectionRepository, geminiService, s.logger)
	detectionHandler := handler.NewDetectionHandler(detectionService, customerService, s.val)

	s.handlers = []Handler{
		userHandler,
		authHandler,
		customerHandler,
		paymentHandler,
		detectionHandler,
	}
}

func (s *ServerConfig) Start() {
	s.DependencyInjection()
	s.HealthCheck()
	rootPath := s.router.Group("/api/v1")

	for _, handler := range s.handlers {
		handler.SetEndpoint(rootPath)
	}

	if err := s.router.Listen(":" + env.GetEnv().Port); err != nil {
		s.logger.Fatal("Failed to start server: ", err)
	}
}

func (s *ServerConfig) Shutdown(ctx context.Context) {
	if err := s.router.Shutdown(); err != nil {
		s.logger.Error("Failed to shutdown server: ", err)
	}

	if err := s.db.Close(); err != nil {
		s.logger.Error("Failed to close database connection: ", err)
	}

	if err := s.logger.Writer().Close(); err != nil {
		s.logger.Error("Failed to close logger writer: ", err)
	}

	s.logger.Info("Server shutdown gracefully")
}

func (s *ServerConfig) HealthCheck() {
	s.router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
