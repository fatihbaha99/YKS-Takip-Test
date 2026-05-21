package backend

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"yks-tracker/database"
	"yks-tracker/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Tüm alanlar zorunlu"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Şifre işlenemedi"})
	}

	now := time.Now()
	purgeAt := now.AddDate(1, 0, 0)

	result, err := database.DB.Exec(
		`INSERT INTO users (name, email, password_hash, created_at, purge_at) VALUES (?, ?, ?, ?, ?)`,
		req.Name, req.Email, string(hash), now, purgeAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return c.Status(409).JSON(fiber.Map{"error": "Bu email zaten kayıtlı"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Kayıt hatası"})
	}

	userID, _ := result.LastInsertId()

	token, err := GenerateToken(userID, req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token oluşturulamadı"})
	}

	return c.Status(201).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    userID,
			"name":  req.Name,
			"email": req.Email,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	var user models.User
	err := database.DB.QueryRow(
		`SELECT id, name, email, password_hash FROM users WHERE email = ? AND active = 1`,
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return c.Status(401).JSON(fiber.Map{"error": "Email veya şifre hatalı"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Veritabanı hatası"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Email veya şifre hatalı"})
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token oluşturulamadı"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var user models.User
	err := database.DB.QueryRow(
		`SELECT id, name, email, telegram_chat_id, activation_code, created_at, purge_at FROM users WHERE id = ?`,
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.TelegramChatID, &user.ActivationCode, &user.CreatedAt, &user.PurgeAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
	}

	return c.JSON(fiber.Map{"user": user})
}

func GenerateActivationCode(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	code := generateCode()

	_, err := database.DB.Exec(`UPDATE users SET activation_code = ? WHERE id = ?`, code, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kod oluşturulamadı"})
	}

	return c.JSON(fiber.Map{"activation_code": code})
}

func DisconnectTelegram(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	result, err := database.DB.Exec(
		`UPDATE users SET telegram_chat_id = 0, activation_code = '' WHERE id = ?`,
		userID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Bağlantı kesilemedi"})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
	}

	return c.JSON(fiber.Map{"message": "Telegram bağlantısı kesildi"})
}

func generateCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GenerateAPIKey(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("YKS Takip Sistemi API v1 - %d", time.Now().Year()),
	})
}
