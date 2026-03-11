package bootstrap

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"wallet-api/app/libraries"
	"wallet-api/app/middleware"
	"wallet-api/app/routes"
	"wallet-api/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
)

type bootstrap struct {
	c       *fiber.App
	mysqlDB *sql.DB
	cfg     *config.Config
}

func NewBootstrap(c *fiber.App, mysqlDB *sql.DB, cfg *config.Config) *bootstrap {
	return &bootstrap{
		c:       c,
		mysqlDB: mysqlDB,
		cfg:     cfg,
	}
}

func (b *bootstrap) Run() error {
	// Initialize middleware stack
	b.initMiddleware()

	// Initialize JWT from config
	middleware.InitJWT(b.cfg.JWT.Secret, b.cfg.JWT.GetJWTExpiration())

	// Initialize Rate Limiter from config
	middleware.InitRateLimiter(b.cfg.RateLimit.Max, b.cfg.RateLimit.Expiration)

	// Initialize routes
	routes.NewRouter(b.c, b.mysqlDB).Init()

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		libraries.Logger.Info().Msg("Shutting down server gracefully...")

		// Close database connection
		if b.mysqlDB != nil {
			if err := b.mysqlDB.Close(); err != nil {
				libraries.Logger.Error().Err(err).Msg("Error closing database connection")
			} else {
				libraries.Logger.Info().Msg("Database connection closed")
			}
		}

		if err := b.c.Shutdown(); err != nil {
			libraries.Logger.Error().Err(err).Msg("Error during shutdown")
		}
	}()

	// Start server with configured port
	libraries.Logger.Info().Str("addr", b.cfg.App.Port).Msg("Starting server...")
	return b.c.Listen(b.cfg.App.Port)
}

func (b *bootstrap) initMiddleware() {
	// Security headers
	b.c.Use(helmet.New())

	// CORS
	b.c.Use(cors.New())

	// Request ID for tracing
	b.c.Use(requestid.New())

	// Structured logging
	b.c.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// Recovery from panics
	b.c.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Rate Limiting
	b.c.Use(middleware.RateLimiter())

	// Input Sanitization
	b.c.Use(middleware.SanitizeInput())

	// Centralized Error Handler (should be last)
	b.c.Use(middleware.ErrorHandler())
}
