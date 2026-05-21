package backend

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"yks-tracker/database"
	"yks-tracker/scheduler"

	"github.com/gofiber/fiber/v2"
)

type StudyRequest struct {
	Subject   string `json:"subject"`
	Topic     string `json:"topic"`
	StudyType string `json:"study_type"`
	Stars     int    `json:"stars"`
	Correct   int    `json:"correct"`
	Wrong     int    `json:"wrong"`
}

func RecordStudy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var req StudyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	if req.Subject == "" || req.Topic == "" || req.StudyType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Eksik alanlar var"})
	}

	var net float64
	switch req.StudyType {
	case "goz_gezdir":
		if req.Stars < 1 || req.Stars > 5 {
			return c.Status(400).JSON(fiber.Map{"error": "Yıldız 1-5 arası olmalı"})
		}
	case "test_coz", "genel_test":
		if req.Correct < 0 || req.Wrong < 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Doğru/yanlış negatif olamaz"})
		}
		net = float64(req.Correct) - float64(req.Wrong)*0.25
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz çalışma tipi"})
	}

	today := time.Now().Format("2006-01-02")

	_, err := database.DB.Exec(
		`INSERT INTO study_sessions (user_id, subject, topic, study_type, stars, correct, wrong, net, study_date)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Subject, req.Topic, req.StudyType, req.Stars, req.Correct, req.Wrong, net, today,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kayıt hatası"})
	}

	createSpacedRepetitionTodos(userID, req.Subject, req.Topic)

	return c.Status(201).JSON(fiber.Map{"message": "Çalışma kaydedildi", "net": net})
}

func createSpacedRepetitionTodos(userID int64, subject, topic string) {
	now := time.Now()
	intervals := []struct {
		days int
		typ  string
	}{
		{1, "goz_gezdir"},
		{7, "test_coz"},
		{30, "genel_test"},
	}

	for _, iv := range intervals {
		dueDate := now.AddDate(0, 0, iv.days).Format("2006-01-02")
		database.DB.Exec(
			`INSERT INTO todos (user_id, subject, topic, todo_type, due_date)
			 VALUES (?, ?, ?, ?, ?)`,
			userID, subject, topic, iv.typ, dueDate,
		)
	}
}

func GetTodayTodos(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	today := time.Now().Format("2006-01-02")

	rows, err := database.DB.Query(
		`SELECT id, subject, topic, todo_type, due_date, completed
		 FROM todos WHERE user_id = ? AND due_date <= ? AND completed = 0
		 ORDER BY due_date ASC`,
		userID, today,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sorgu hatası"})
	}
	defer rows.Close()

	type TodoItem struct {
		ID       int64  `json:"id"`
		Subject  string `json:"subject"`
		Topic    string `json:"topic"`
		TodoType string `json:"todo_type"`
		DueDate  string `json:"due_date"`
	}

	var todos []TodoItem
	for rows.Next() {
		var t TodoItem
		var completed int
		if err := rows.Scan(&t.ID, &t.Subject, &t.Topic, &t.TodoType, &t.DueDate, &completed); err != nil {
			continue
		}
		todos = append(todos, t)
	}

	if todos == nil {
		todos = []TodoItem{}
	}

	return c.JSON(fiber.Map{"todos": todos, "date": today})
}

func GetCalendarData(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	year := c.Query("year")
	m := c.Query("month")
	monthInt, _ := strconv.Atoi(m)
	month := fmt.Sprintf("%02d", monthInt)

	studyRows, err := database.DB.Query(`
		SELECT study_date, COUNT(*) as session_count,
			SUM(CASE WHEN study_type = 'goz_gezdir' THEN 1 ELSE 0 END) as review_count,
			SUM(CASE WHEN study_type IN ('test_coz','genel_test') THEN 1 ELSE 0 END) as test_count,
			COALESCE(SUM(net), 0) as total_net
		FROM study_sessions
		WHERE user_id = ? AND strftime('%Y', study_date) = ? AND strftime('%m', study_date) = ?
		GROUP BY study_date ORDER BY study_date ASC
	`, userID, year, month)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sorgu hatası"})
	}

	type DayData struct {
		Date         string `json:"date"`
		Sessions     int    `json:"sessions"`
		Reviews      int    `json:"reviews"`
		Tests        int    `json:"tests"`
		TotalNet     float64 `json:"total_net"`
		Todos        int    `json:"todos"`
	}

	var days []DayData
	for studyRows.Next() {
		var d DayData
		if err := studyRows.Scan(&d.Date, &d.Sessions, &d.Reviews, &d.Tests, &d.TotalNet); err != nil {
			continue
		}
		days = append(days, d)
	}
	studyRows.Close()

	todoRows, err := database.DB.Query(`
		SELECT due_date, COUNT(*) FROM todos
		WHERE user_id = ? AND strftime('%Y', due_date) = ? AND strftime('%m', due_date) = ?
		GROUP BY due_date ORDER BY due_date ASC
	`, userID, year, month)
	if err == nil {
		for todoRows.Next() {
			var date string
			var count int
			if err := todoRows.Scan(&date, &count); err != nil {
				continue
			}
			found := false
			for i := range days {
				if days[i].Date == date {
					days[i].Todos = count
					found = true
					break
				}
			}
			if !found {
				days = append(days, DayData{Date: date, Todos: count})
			}
		}
		todoRows.Close()
	}

	if days == nil {
		days = []DayData{}
	}

	return c.JSON(fiber.Map{"days": days, "year": year, "month": month})
}

func GetDaySessions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	date := c.Query("date")

	rows, err := database.DB.Query(`
		SELECT id, subject, topic, study_type, stars, correct, wrong, net, created_at
		FROM study_sessions
		WHERE user_id = ? AND study_date = ?
		ORDER BY created_at ASC
	`, userID, date)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sorgu hatası"})
	}
	defer rows.Close()

	type Session struct {
		ID        int64   `json:"id"`
		Subject   string  `json:"subject"`
		Topic     string  `json:"topic"`
		StudyType string  `json:"study_type"`
		Stars     int     `json:"stars"`
		Correct   int     `json:"correct"`
		Wrong     int     `json:"wrong"`
		Net       float64 `json:"net"`
		CreatedAt string  `json:"created_at"`
	}

	var sessions []Session
	for rows.Next() {
		var s Session
		var createdAt string
		if err := rows.Scan(&s.ID, &s.Subject, &s.Topic, &s.StudyType, &s.Stars, &s.Correct, &s.Wrong, &s.Net, &createdAt); err != nil {
			continue
		}
		sessions = append(sessions, s)
	}

	if sessions == nil {
		sessions = []Session{}
	}

	return c.JSON(fiber.Map{"sessions": sessions, "date": date})
}

func GetDayTodos(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	date := c.Query("date")

	rows, err := database.DB.Query(`
		SELECT id, subject, topic, todo_type, completed
		FROM todos
		WHERE user_id = ? AND due_date = ?
		ORDER BY completed ASC, created_at ASC
	`, userID, date)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sorgu hatası"})
	}
	defer rows.Close()

	type TodoItem struct {
		ID        int64  `json:"id"`
		Subject   string `json:"subject"`
		Topic     string `json:"topic"`
		TodoType  string `json:"todo_type"`
		Completed bool   `json:"completed"`
	}

	var todos []TodoItem
	for rows.Next() {
		var t TodoItem
		var c int
		if err := rows.Scan(&t.ID, &t.Subject, &t.Topic, &t.TodoType, &c); err != nil {
			continue
		}
		t.Completed = c == 1
		todos = append(todos, t)
	}

	if todos == nil {
		todos = []TodoItem{}
	}

	return c.JSON(fiber.Map{"todos": todos, "date": date})
}

func DownloadExcel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	fileName, err := scheduler.GenerateUserExcel(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Excel oluşturulamadı"})
	}

	c.SendFile(fileName)
	os.Remove(fileName)
	return nil
}

func CompleteTodo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	id := c.Params("id")

	result, err := database.DB.Exec(
		`UPDATE todos SET completed = 1 WHERE id = ? AND user_id = ?`,
		id, userID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Güncelleme hatası"})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Görev bulunamadı"})
	}

	return c.JSON(fiber.Map{"message": "Görev tamamlandı"})
}
