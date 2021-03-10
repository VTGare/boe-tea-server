package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/VTGare/boe-tea-backend/internal/config"
	"github.com/VTGare/boe-tea-backend/internal/database"
	"github.com/VTGare/boe-tea-backend/pkg/artworks"
	"github.com/VTGare/boe-tea-backend/pkg/guilds"
	"github.com/VTGare/boe-tea-backend/pkg/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func createFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		//Limit body size to 1 MB
		BodyLimit:    1 * 1024 * 1024,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).JSON(fiber.Map{
				"status":  code,
				"message": err.Error(),
			})

			if err != nil {
				return ctx.Status(500).SendString("Internal Server Error")
			}

			return nil
		},
	})

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	return app
}

func main() {
	var (
		configPath string
		environ    string
		logger     *zap.Logger
		err        error
	)

	flag.StringVar(&configPath, "config", "config.json", "Path to JSON configuration file.")
	flag.StringVar(&environ, "env", "dev", "Environment type. Accepts one of: dev, prod")
	flag.Parse()

	switch environ {
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	default:
		log.Fatal("Unknown environment type: " + environ)
	}

	if err != nil {
		fmt.Printf("Error creating a logger: %v", err)
		os.Exit(1)
	}

	cfg, err := config.FromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.ConnectionURI(), cfg.Database.DefaultDB)
	if err != nil {
		log.Fatal(err)
	}

	app := createFiber()
	env := &handlers.Env{
		Guilds:   guilds.NewService(db, logger.Sugar()),
		Artworks: artworks.NewService(db, logger.Sugar()),
	}

	api := app.Group("/api")
	handlers.Register(api, env)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Print("Shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Print(err)
	}

	if err := db.Close(); err != nil {
		log.Print(err)
	}

	os.Exit(0)
}
