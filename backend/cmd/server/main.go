package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"enterprise/sizing-apps-backend/internal/handler"
	"enterprise/sizing-apps-backend/internal/platform/db"
	"enterprise/sizing-apps-backend/internal/service"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:   "Enterprise Sizing & Capacity Engine v1.0",
		BodyLimit: 50 * 1024 * 1024,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	baseDir, _ := os.Getwd()
	templatesDir := filepath.Join(baseDir, "templates")
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		templatesDir = filepath.Join(baseDir, "..", "templates")
	}

	fmt.Printf("🚀 Starting Sizing Backend Engine...\n")
	fmt.Printf("📁 Templates directory: %s\n", templatesDir)

	database, err := db.InitDB()
	if err != nil {
		log.Printf("⚠️ Warning: Database initialization error: %v", err)
	}

	sizingSvc := service.NewSizingService()
	sizingHandler := handler.NewSizingHandler(sizingSvc, database, templatesDir)
	sizingHandler.RegisterRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		addr := ":" + port
		fmt.Printf("⚡ Server running on http://localhost%s\n", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	<-stopChan
	fmt.Println("\n🛑 Shutting down server gracefully...")
	_ = app.Shutdown()
	fmt.Println("✅ Server stopped.")
}
