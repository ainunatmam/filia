package bootstrap

import (
	"mini-exchange/app/libraries"
	"mini-exchange/app/middleware"
	"mini-exchange/app/routes"
	"mini-exchange/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
)

type bootstrap struct {
	c   *fiber.App
	cfg *config.Config
}

func NewBootstrap(c *fiber.App, cfg *config.Config) *bootstrap {
	return &bootstrap{
		c:   c,
		cfg: cfg,
	}
}

func (b *bootstrap) Run() error {
	b.initMiddleware()
	middleware.InitJWT(b.cfg.JWT.Secret, b.cfg.JWT.GetJWTExpiration())
	middleware.InitRateLimiter(b.cfg.RateLimit.Max, b.cfg.RateLimit.Expiration)
	routes.NewRouter(b.c, b.cfg).Init()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		libraries.Logger.Info().Msg("Shutting down server gracefully...")

		if err := b.c.Shutdown(); err != nil {
			libraries.Logger.Error().Err(err).Msg("Error during shutdown")
		}
	}()

	libraries.Logger.Info().Str("addr", b.cfg.App.Port).Msg("Starting server...")
	return b.c.Listen(b.cfg.App.Port)
}

func (b *bootstrap) initMiddleware() {
	b.c.Use(helmet.New())
	b.c.Use(cors.New())
	b.c.Use(requestid.New())
	b.c.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	b.c.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	b.c.Use(middleware.RateLimiter())
	b.c.Use(middleware.SanitizeInput())
	b.c.Use(middleware.ErrorHandler())
}
