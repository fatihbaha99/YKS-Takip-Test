package main

import (
	"log"
	"os"
	"path/filepath"

	"yks-tracker/backend"
	"yks-tracker/bot"
	"yks-tracker/database"
	"yks-tracker/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		exe, _ := os.Executable()
		dbPath = filepath.Join(filepath.Dir(exe), "uygulama.db")
	}

	if err := database.Init(dbPath); err != nil {
		log.Fatalf("Veritabanı başlatılamadı: %v", err)
	}
	defer database.Close()

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Migration hatası: %v", err)
	}

	frontendPath := os.Getenv("FRONTEND_PATH")
	if frontendPath == "" {
		frontendPath = "./frontend"
	}

	frontendApp := fiber.New()
	backendApp := fiber.New()

	frontendApp.Use(logger.New())
	backendApp.Use(logger.New())
	backendApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	go func() {
		log.Println("[Frontend] Port 4000'de yayında - statik dosyalar")
		frontendApp.Static("/", frontendPath)
		if err := frontendApp.Listen(":4000"); err != nil {
			log.Fatalf("Frontend sunucu hatası: %v", err)
		}
	}()

	bot.Start(os.Getenv("TELEGRAM_BOT_TOKEN"))
	go scheduler.StartNightlyCleanup()

	backend.SetupRoutes(backendApp)

	go func() {
		log.Println("[Backend] Port 4080'de yayında - REST API")
		if err := backendApp.Listen(":4080"); err != nil {
			log.Fatalf("Backend sunucu hatası: %v", err)
		}
	}()

	select {}
}
