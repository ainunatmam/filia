package main

import (
	"database/sql"
	"fmt"
	"log"
	"mini-exchange/app/libraries"
	"mini-exchange/bootstrap"
	"mini-exchange/config"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pressly/goose/v3"
)

func main() {
	// Load configuration (includes .env loading)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	libraries.InitLogger()

	mysql, err := bootstrap.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
		panic(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigration(mysql)
		return
	}

	// Initialize background worker
	

	app := fiber.New(fiber.Config{
		BodyLimit:    4 * 1024 * 1024, // 4MB max body size
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Add health check endpoints before other routes
	setupHealthChecks(app, mysql)

	err = bootstrap.NewBootstrap(app, mysql, cfg).Run()
	if err != nil {
		panic(err)
	}
}

func setupHealthChecks(app *fiber.App, db *sql.DB) {
	// Liveness check - returns OK if app is running
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Readiness check - returns OK if app can handle requests (DB connected)
	app.Get("/ready", func(c *fiber.Ctx) error {
		err := db.Ping()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"error":  "database connection failed",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ready",
		})
	})
}

func runMigration(db *sql.DB) {
	if len(os.Args) < 3 {
		log.Println("Usage:")
		fmt.Println("  migrate up")
		fmt.Println("  migrate down")
		fmt.Println("  migrate status")
		fmt.Println("  migrate create <name>")
		log.Fatal("Specify migrate command")
	}

	action := os.Args[2]
	goose.SetDialect("mysql")
	switch action {
	case "up":
		if err := goose.Up(db, "./database/migration"); err != nil {
			log.Fatal(err)
		}

	case "down":
		if err := goose.Down(db, "./database/migration"); err != nil {
			log.Fatal(err)
		}

	case "status":
		if err := goose.Status(db, "./database/migration"); err != nil {
			log.Fatal(err)
		}

	case "create":
		if len(os.Args) < 4 {
			log.Fatal("Specify migration name")
		}
		name := os.Args[3]

		if err := goose.Create(db, "./database/migration", name, "sql"); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("Unknown migrate action")
	}
}
