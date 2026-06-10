package scheduler

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"yks-tracker/database"

	"github.com/xuri/excelize/v2"
)

var turkeyLocScheduler = time.FixedZone("TRT", 3*60*60)

func StartNightlyCleanup() {
	for {
		now := time.Now().In(turkeyLocScheduler)
		next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.AddDate(0, 0, 1)
		}
		time.Sleep(time.Until(next))
		cleanupExpiredUsers()
	}
}

func cleanupExpiredUsers() {
	log.Println("[Scheduler] Süresi dolan kullanıcılar temizleniyor...")

	rows, err := database.DB.Query(`
		SELECT id, name, email FROM users WHERE purge_at <= datetime('now') AND active = 1
	`)
	if err != nil {
		log.Printf("[Scheduler] Sorgu hatası: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			continue
		}
		if err := exportAndDeleteUser(id, name, email); err != nil {
			log.Printf("[Scheduler] Kullanıcı %d temizlenemedi: %v", id, err)
		} else {
			log.Printf("[Scheduler] Kullanıcı %d (%s) temizlendi", id, email)
		}
	}
}

func exportAndDeleteUser(userID int64, name, email string) error {
	fileName, err := GenerateUserExcel(userID)
	if err != nil {
		return err
	}

	sendEmail(email, name, fileName)
	deleteUserData(userID)
	os.Remove(fileName)
	return nil
}

func GenerateUserExcel(userID int64) (string, error) {
	studyRows, err := database.DB.Query(`
		SELECT subject, topic, study_type, stars, correct, wrong, net, study_date, created_at
		FROM study_sessions WHERE user_id = ? ORDER BY created_at ASC
	`, userID)
	if err != nil {
		return "", err
	}
	defer studyRows.Close()

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Çalışma Geçmişi"
	f.SetSheetName("Sheet1", sheet)
	f.SetCellValue(sheet, "A1", "Ders")
	f.SetCellValue(sheet, "B1", "Konu")
	f.SetCellValue(sheet, "C1", "Çalışma Tipi")
	f.SetCellValue(sheet, "D1", "Yıldız")
	f.SetCellValue(sheet, "E1", "Doğru")
	f.SetCellValue(sheet, "F1", "Yanlış")
	f.SetCellValue(sheet, "G1", "Net")
	f.SetCellValue(sheet, "H1", "Çalışma Tarihi")
	f.SetCellValue(sheet, "I1", "Kayıt Tarihi")

	rowIdx := 2
	for studyRows.Next() {
		var subject, topic, studyType, studyDate, createdAt string
		var stars, correct, wrong int
		var net float64
		if err := studyRows.Scan(&subject, &topic, &studyType, &stars, &correct, &wrong, &net, &studyDate, &createdAt); err != nil {
			continue
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), subject)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), topic)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx), studyType)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIdx), stars)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIdx), correct)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIdx), wrong)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIdx), net)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", rowIdx), studyDate)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", rowIdx), createdAt)
		rowIdx++
	}

	todoRows, err := database.DB.Query(`
		SELECT subject, topic, todo_type, due_date, completed, created_at
		FROM todos WHERE user_id = ? ORDER BY created_at ASC
	`, userID)
	if err != nil {
		return "", err
	}
	defer todoRows.Close()

	sheet2 := "Görevler"
	f.NewSheet(sheet2)
	f.SetCellValue(sheet2, "A1", "Ders")
	f.SetCellValue(sheet2, "B1", "Konu")
	f.SetCellValue(sheet2, "C1", "Görev Tipi")
	f.SetCellValue(sheet2, "D1", "Son Tarih")
	f.SetCellValue(sheet2, "E1", "Tamamlandı")
	f.SetCellValue(sheet2, "F1", "Oluşturulma")

	rowIdx = 2
	for todoRows.Next() {
		var subject, topic, todoType, dueDate, createdAt string
		var completed int
		if err := todoRows.Scan(&subject, &topic, &todoType, &dueDate, &completed, &createdAt); err != nil {
			continue
		}
		f.SetCellValue(sheet2, fmt.Sprintf("A%d", rowIdx), subject)
		f.SetCellValue(sheet2, fmt.Sprintf("B%d", rowIdx), topic)
		f.SetCellValue(sheet2, fmt.Sprintf("C%d", rowIdx), todoType)
		f.SetCellValue(sheet2, fmt.Sprintf("D%d", rowIdx), dueDate)
		f.SetCellValue(sheet2, fmt.Sprintf("E%d", rowIdx), completed)
		f.SetCellValue(sheet2, fmt.Sprintf("F%d", rowIdx), createdAt)
		rowIdx++
	}

	fileName := fmt.Sprintf("yks_veri_%d_%s.xlsx", userID, time.Now().In(turkeyLocScheduler).Format("20060102"))
	if err := f.SaveAs(fileName); err != nil {
		return "", err
	}

	return fileName, nil
}

func sendEmail(to, name, fileName string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpUser == "" {
		log.Printf("[Scheduler] SMTP ayarları eksik, %s için e-posta gönderilemedi", to)
		return
	}

	subject := "YKS Takip - Verileriniz"
	body := fmt.Sprintf("Sayın %s,\n\nYKS Takip uygulamasındaki verileriniz 1 yılı doldurduğu için silinmiştir.\nVerilerinizin yedeği ekteki Excel dosyasında sunulmuştur.\n\nTeşekkür ederiz.", name)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", smtpUser, to, subject, body)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	if err := smtp.SendMail(addr, auth, smtpUser, []string{to}, []byte(msg)); err != nil {
		log.Printf("[Scheduler] E-posta gönderilemedi %s: %v", to, err)
		return
	}
	log.Printf("[Scheduler] E-posta gönderildi: %s", to)
}

func deleteUserData(userID int64) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM study_sessions WHERE user_id = ?", userID)
	tx.Exec("DELETE FROM todos WHERE user_id = ?", userID)
	tx.Exec("DELETE FROM users WHERE id = ?", userID)

	return tx.Commit()
}
