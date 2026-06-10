package backend

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"yks-tracker/database"
	"yks-tracker/models"

	"github.com/gofiber/fiber/v2"
)

var tytSubjects = []string{"Türkçe", "Tarih", "Coğrafya", "Felsefe", "Din Kültürü", "Matematik", "Geometri", "Fizik", "Kimya", "Biyoloji"}
var aytSubjects = []string{"Türk Dili ve Edebiyatı", "Tarih-1", "Coğrafya-1", "Tarih-2", "Coğrafya-2", "Felsefe Grubu", "Din Kültürü", "Matematik", "Geometri", "Fizik", "Kimya", "Biyoloji"}

func CreateExam(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var req struct {
		ExamType string `json:"exam_type"`
		Title    string `json:"title"`
		ExamDate string `json:"exam_date"`
		Results  []struct {
			Subject string `json:"subject"`
			Correct int    `json:"correct"`
			Wrong   int    `json:"wrong"`
		} `json:"results"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	if req.ExamType != "TYT" && req.ExamType != "AYT" {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz sınav türü"})
	}
	if req.ExamDate == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Tarih zorunlu"})
	}

	now := time.Now()

	result, err := database.DB.Exec(
		`INSERT INTO exams (user_id, exam_type, title, exam_date, created_at) VALUES (?, ?, ?, ?, ?)`,
		userID, req.ExamType, req.Title, req.ExamDate, now,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Deneme kaydedilemedi"})
	}

	examID, _ := result.LastInsertId()

	for _, r := range req.Results {
		net := float64(r.Correct) - float64(r.Wrong)*0.25
		_, err := database.DB.Exec(
			`INSERT INTO exam_results (exam_id, subject, correct, wrong, net) VALUES (?, ?, ?, ?, ?)`,
			examID, r.Subject, r.Correct, r.Wrong, net,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ders sonuçları kaydedilemedi"})
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"message":  "Deneme kaydedildi",
		"exam_id":  examID,
	})
}

func ListExams(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	examType := c.Query("type", "AYT")
	sort := c.Query("sort", "date_desc")

	if examType != "TYT" && examType != "AYT" {
		examType = "AYT"
	}

	orderClause := "ORDER BY e.exam_date DESC"
	switch sort {
	case "date_asc":
		orderClause = "ORDER BY e.exam_date ASC"
	case "net_desc":
		orderClause = "ORDER BY total_net DESC"
	case "net_asc":
		orderClause = "ORDER BY total_net ASC"
	}

	query := fmt.Sprintf(`
		SELECT e.id, e.exam_type, e.title, e.exam_date, e.created_at,
			COALESCE(SUM(er.net), 0) as total_net
		FROM exams e
		LEFT JOIN exam_results er ON er.exam_id = e.id
		WHERE e.user_id = ? AND e.exam_type = ?
		GROUP BY e.id
		%s`, orderClause)

	rows, err := database.DB.Query(query, userID, examType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Liste alınamadı"})
	}
	defer rows.Close()

	type ExamSummary struct {
		ID        int64     `json:"id"`
		ExamType  string    `json:"exam_type"`
		Title     string    `json:"title"`
		ExamDate  string    `json:"exam_date"`
		CreatedAt time.Time `json:"created_at"`
		TotalNet  float64   `json:"total_net"`
	}

	var exams []ExamSummary
	for rows.Next() {
		var e ExamSummary
		if err := rows.Scan(&e.ID, &e.ExamType, &e.Title, &e.ExamDate, &e.CreatedAt, &e.TotalNet); err != nil {
			continue
		}
		exams = append(exams, e)
	}

	if exams == nil {
		exams = []ExamSummary{}
	}

	return c.JSON(fiber.Map{"exams": exams})
}

func GetExam(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	examID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz ID"})
	}

	var exam models.Exam
	err = database.DB.QueryRow(
		`SELECT id, user_id, exam_type, title, exam_date, created_at FROM exams WHERE id = ? AND user_id = ?`,
		examID, userID,
	).Scan(&exam.ID, &exam.UserID, &exam.ExamType, &exam.Title, &exam.ExamDate, &exam.CreatedAt)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Deneme bulunamadı"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Veritabanı hatası"})
	}

	rrows, err := database.DB.Query(
		`SELECT id, exam_id, subject, correct, wrong, net FROM exam_results WHERE exam_id = ? ORDER BY id`,
		examID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sonuçlar alınamadı"})
	}
	defer rrows.Close()

	for rrows.Next() {
		var r models.ExamResult
		if err := rrows.Scan(&r.ID, &r.ExamID, &r.Subject, &r.Correct, &r.Wrong, &r.Net); err != nil {
			continue
		}
		exam.Results = append(exam.Results, r)
	}

	if exam.Results == nil {
		exam.Results = []models.ExamResult{}
	}

	return c.JSON(fiber.Map{"exam": exam})
}

func DeleteExam(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	examID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz ID"})
	}

	result, err := database.DB.Exec(
		`DELETE FROM exams WHERE id = ? AND user_id = ?`,
		examID, userID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Silinemedi"})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Deneme bulunamadı"})
	}

	return c.JSON(fiber.Map{"message": "Deneme silindi"})
}

func GetExamSubjects(c *fiber.Ctx) error {
	examType := c.Query("type", "AYT")
	if examType == "TYT" {
		return c.JSON(fiber.Map{"subjects": tytSubjects})
	}
	return c.JSON(fiber.Map{"subjects": aytSubjects})
}
