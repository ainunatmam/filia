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
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	libraries.InitLogger()

	app := fiber.New(fiber.Config{
		BodyLimit:    4 * 1024 * 1024,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	err = bootstrap.NewBootstrap(app, cfg).Run()
	if err != nil {
		panic(err)
	}
}

func setupHealthChecks(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
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
