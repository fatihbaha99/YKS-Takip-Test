package backend

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/", GenerateAPIKey)

	auth := api.Group("/auth")
	auth.Post("/register", Register)
	auth.Post("/login", Login)

	protected := api.Group("")
	protected.Use(AuthRequired)

	protected.Get("/profile", GetProfile)
	protected.Post("/auth/activation-code", GenerateActivationCode)
	protected.Post("/auth/disconnect-telegram", DisconnectTelegram)

	protected.Post("/study", RecordStudy)
	protected.Get("/todos/today", GetTodayTodos)
	protected.Post("/todos/:id/complete", CompleteTodo)
	protected.Get("/calendar", GetCalendarData)
	protected.Get("/calendar/day", GetDaySessions)
	protected.Get("/calendar/day/todos", GetDayTodos)
	protected.Get("/export/excel", DownloadExcel)
}
