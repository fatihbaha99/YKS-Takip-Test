package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"yks-tracker/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI
var turkeyLoc = time.FixedZone("TRT", 3*60*60)

var subjectsData = map[string]map[string][]string{
	"TYT": {
		"Türkçe":     {"Sözcükte Anlam", "Cümlede Anlam", "Paragraf", "Ses Bilgisi", "Yazım Kuralları", "Noktalama İşaretleri", "Sözcük Türleri", "Fiiller", "Cümlenin Ögeleri", "Cümle Çeşitleri", "Anlatım Bozukluğu"},
		"Matematik":  {"Temel Kavramlar", "Sayı Basamakları", "Bölme-Bölünebilme", "EBOB-EKOK", "Rasyonel Sayılar", "Basit Eşitsizlikler", "Mutlak Değer", "Üslü Sayılar", "Köklü Sayılar", "Çarpanlara Ayırma", "Oran-Orantı", "Denklem Çözme", "Problemler", "Kümeler", "Fonksiyonlar", "Permütasyon", "Kombinasyon", "Olasılık", "İstatistik"},
		"Fizik":      {"Fizik Bilimine Giriş", "Madde ve Özellikleri", "Hareket ve Kuvvet", "Enerji", "Isı ve Sıcaklık", "Elektrik", "Manyetizma", "Dalgalar", "Optik"},
		"Kimya":      {"Kimya Bilimi", "Atom ve Periyodik Sistem", "Kimyasal Türler Arası Etkileşimler", "Maddenin Halleri", "Doğa ve Kimya", "Asitler ve Bazlar", "Kimya Her Yerde"},
		"Biyoloji":   {"Yaşam Bilimi Biyoloji", "Canlıların Yapısında Bulunan Organik Bileşikler", "Hücre", "Canlıların Çeşitliliği ve Sınıflandırılması", "Hücre Bölünmeleri", "Kalıtım", "Ekoloji"},
		"Tarih":      {"Tarih ve Zaman", "İnsanlığın İlk Dönemleri", "Orta Çağ'da Dünya", "İlk ve Orta Çağ'da Türkler", "İslamiyet'in Doğuşu", "Türklerin İslamiyet'i Kabulü", "Beylikten Devlete", "Dünya Gücü Osmanlı", "Değişen Dünya ve Avrupa", "Osmanlı Kültür ve Medeniyeti"},
		"Coğrafya":   {"Doğa ve İnsan", "Dünya'nın Şekli ve Hareketleri", "Coğrafi Koordinat Sistemi", "Harita Bilgisi", "Atmosfer ve İklim", "Yeryüzündeki İklim Tipleri", "Topoğrafya ve Kayaçlar", "İç Kuvvetler", "Dış Kuvvetler", "Nüfus", "Göç", "Türkiye'de Yerleşme"},
		"Felsefe":    {"Felsefenin Alanı", "Bilgi Felsefesi", "Varlık Felsefesi", "Ahlak Felsefesi", "Sanat Felsefesi", "Din Felsefesi", "Siyaset Felsefesi", "Bilim Felsefesi"},
		"Din Kültürü": {"İnsan ve Din", "Allah İnancı", "İbadet", "Kur'an'da Bazı Kavramlar", "İslam ve Toplum", "Haklar ve Özgürlükler"},
	},
	"AYT": {
		"Matematik":           {"Trigonometri", "Logaritma", "Diziler", "Limit", "Türev", "İntegral", "Matris-Determinant", "Karmaşık Sayılar", "Polinomlar", "Çember ve Daire", "Katı Cisimler", "Analitik Geometri", "Olasılık ve İstatistik"},
		"Türk Dili ve Edebiyatı": {"Türk Edebiyatının Dönemleri", "İslamiyet Öncesi Türk Edebiyatı", "İslami Dönem Türk Edebiyatı", "Halk Edebiyatı", "Divan Edebiyatı", "Tanzimat Edebiyatı", "Servetifünun Edebiyatı", "Fecriati Edebiyatı", "Milli Edebiyat", "Cumhuriyet Dönemi Edebiyatı", "Günümüz Türk Edebiyatı"},
		"Fizik":               {"Vektörler", "Bağıl Hareket", "Dinamik", "İş-Güç-Enerji", "Atışlar", "Dönme Hareketi", "Elektrik Alan", "Manyetik Alan", "İndüksiyon", "Alternatif Akım", "Dalga Mekaniği", "Atom Fiziği", "Radyoaktivite"},
		"Kimya":               {"Modern Atom Teorisi", "Gazlar", "Sıvı Çözeltiler", "Kimyasal Tepkimelerde Enerji", "Tepkime Hızları", "Kimyasal Denge", "Asit-Baz Dengesi", "Çözünürlük Dengesi", "Elektrokimya"},
		"Biyoloji":            {"Sinir Sistemi", "Endokrin Sistem", "Duyu Sistemleri", "Destek ve Hareket Sistemi", "Dolaşım Sistemi", "Solunum Sistemi", "Boşaltım Sistemi", "Üreme Sistemi", "Bağışıklık Sistemi", "Bitki Biyolojisi", "Canlılar ve Çevre", "DNA ve Genetik", "Evrim"},
		"Tarih":               {"Tarih Bilimi", "İlk Çağ Uygarlıkları", "Türk Tarihi", "Osmanlı Tarihi", "Avrupa Tarihi", "Yakın Çağ Dünya Tarihi", "Türkiye Cumhuriyeti Tarihi", "Soğuk Savaş Dönemi", "Küreselleşen Dünya"},
		"Coğrafya":            {"Doğal Sistemler", "Beşeri Sistemler", "Bölgeler ve Ülkeler", "Küresel Ortam", "Türkiye'nin Coğrafi Konumu", "Türkiye'de İklim", "Türkiye'de Yer Şekilleri", "Türkiye'de Nüfus", "Türkiye'de Tarım", "Türkiye'de Sanayi"},
		"Felsefe":             {"Psikoloji Bilimi", "Davranış ve Süreçleri", "Öğrenme", "Bellek", "Düşünme ve Dil", "Sosyoloji Bilimi", "Toplumsal Yapı", "Toplumsal Değişme", "Kültür", "Toplumsal Kurumlar"},
		"Din Kültürü":         {"Vahiy ve Akıl", "İnanç ve İbadet", "İslam Ahlakı", "Kur'an'da İnsan", "Hak ve Sorumluluk", "İslam ve Bilim"},
	},
}

func listSubjects(exam string) string {
	subjects := subjectsData[exam]
	keys := make([]string, 0, len(subjects))
	for k := range subjects {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func listTopics(exam, subject string) string {
	topics := subjectsData[exam][subject]
	return strings.Join(topics, ", ")
}

type StudyFlow struct {
	Step      int
	Exam      string
	Subject   string
	Topic     string
	StudyType string
	Stars     int
	Correct   int
	Wrong     int
}

var (
	studyFlows   = make(map[int64]*StudyFlow)
	studyFlowsMu sync.Mutex
)

func Start(token string) {
	if token == "" {
		token = os.Getenv("TELEGRAM_BOT_TOKEN")
	}
	if token == "" {
		log.Println("[Bot] TELEGRAM_BOT_TOKEN bulunamadı, bot başlatılmadı")
		return
	}

	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("[Bot] Başlatılamadı: %v", err)
		return
	}

	log.Printf("[Bot] %s olarak başlatıldı", Bot.Self.UserName)

	go startPolling()
	go startDailyReminder()
}

func startPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := update.Message
		chatID := msg.Chat.ID
		text := strings.TrimSpace(msg.Text)

		studyFlowsMu.Lock()
		_, inFlow := studyFlows[chatID]
		studyFlowsMu.Unlock()

	if inFlow {
		if msg.IsCommand() && text == "/iptal" {
			cancelStudyFlow(chatID)
			continue
		}
		if !msg.IsCommand() {
			handleStudyFlowInput(chatID, text)
			continue
		}
	}

		if msg.IsCommand() {
			handleCommand(msg)
		}
	}
}

func handleCommand(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := msg.Text
	parts := strings.Fields(text)
	cmd := parts[0]

	if !isUserActive(chatID) && cmd != "/start" && !strings.HasPrefix(cmd, "/aktifkod") {
		sendMessage(chatID, "Önce /aktifkod KOD ile botu aktifleştirmelisin.")
		return
	}

	switch cmd {
	case "/start":
		sendMessage(chatID, "Merhaba! YKS Takip botuna hoş geldin.\nWeb panelinden aldığın aktivasyon kodunu /aktifkod KOD şeklinde göndererek botu aktifleştirebilirsin.\nBağlantını kesmek için /botbaglantikes yazabilirsin.\nTüm komutlar için /help yaz.")
	case "/help":
		sendMessage(chatID, "Kullanılabilir komutlar:\n\n/start - Hoş geldin mesajı\n/help - Bu yardım mesajı\n/aktifkod KOD - Botu web panel ile eşleştir\n/botbaglantikes - Bot bağlantısını kes\n/iptal - Yarıda kalan işlemi iptal et\n/bugun - Bugünkü görevleri listele\n/haftalik - 7 günlük görev takvimi\n/tamamla ID - Bir görevi tamamla\n/calistim - Adım adım çalışma kaydet")
	case "/iptal":
		cancelStudyFlow(chatID)
		sendMessage(chatID, "Mevcut işlem iptal edildi.")
	case "/botbaglantikes":
		deactivateUser(chatID)
	case "/bugun":
		sendTodayTodos(chatID)
	case "/haftalik":
		sendWeeklyTodos(chatID)
	case "/tamamla":
		if len(parts) < 2 {
			sendMessage(chatID, "Kullanım: /tamamla GOREV_ID (ID'yi /bugun ile görebilirsin)")
			return
		}
		id, _ := strconv.Atoi(parts[1])
		completeTodoTelegram(chatID, int64(id))
	case "/calistim":
		startStudyFlow(chatID)
	default:
		if len(text) > 10 && text[:9] == "/aktifkod" {
			code := text[10:]
			activateUser(chatID, code)
		} else {
			sendMessage(chatID, "Tanınmayan komut. /help yazarak tüm komutları görebilirsin.")
		}
	}
}

func isUserActive(chatID int64) bool {
	var count int
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE telegram_chat_id = ? AND active = 1`, chatID).Scan(&count)
	return count > 0
}

func startStudyFlow(chatID int64) {
	studyFlowsMu.Lock()
	studyFlows[chatID] = &StudyFlow{Step: 1}
	studyFlowsMu.Unlock()
	sendMessage(chatID, "📚 Yeni çalışma kaydı\n\nSınav türünü yaz (TYT veya AYT):")
}

func cancelStudyFlow(chatID int64) {
	studyFlowsMu.Lock()
	delete(studyFlows, chatID)
	studyFlowsMu.Unlock()
	sendMessage(chatID, "Çalışma kaydı iptal edildi.")
}

func handleStudyFlowInput(chatID int64, text string) {
	studyFlowsMu.Lock()
	flow, exists := studyFlows[chatID]
	studyFlowsMu.Unlock()
	if !exists {
		return
	}

	switch flow.Step {
	case 1:
		text = strings.ToUpper(text)
		if text != "TYT" && text != "AYT" {
			sendMessage(chatID, "Sadece TYT veya AYT yazabilirsin. Tekrar dene:")
			return
		}
		flow.Exam = text
		flow.Step = 2
		subjects := listSubjects(text)
		sendMessage(chatID, "Ders seç:\n"+subjects)

	case 2:
		examTopics, ok := subjectsData[flow.Exam]
		if !ok {
			sendMessage(chatID, "Geçersiz sınav türü. Tekrar başlat.")
			cancelStudyFlow(chatID)
			return
		}
		found := false
		for s := range examTopics {
			if strings.EqualFold(s, text) {
				flow.Subject = s
				found = true
				break
			}
		}
		if !found {
			subjects := listSubjects(flow.Exam)
			sendMessage(chatID, "Geçersiz ders. Şunlardan birini yaz:\n"+subjects)
			return
		}
		flow.Step = 3
		topics := listTopics(flow.Exam, flow.Subject)
		sendMessage(chatID, "Konu seç:\n"+topics)

	case 3:
		topics := subjectsData[flow.Exam][flow.Subject]
		found := false
		for _, t := range topics {
			if strings.EqualFold(t, text) {
				flow.Topic = t
				found = true
				break
			}
		}
		if !found {
			sendMessage(chatID, "Geçersiz konu. Şunlardan birini yaz:\n"+listTopics(flow.Exam, flow.Subject))
			return
		}
		flow.Step = 4
		sendMessage(chatID, "Çalışma tipini yaz:\n- goz (Göz Gezdirme)\n- test (Test Çözümü)\n- genel (Genel Test)")

	case 4:
		switch text {
		case "goz":
			flow.StudyType = "goz_gezdir"
			flow.Step = 5
			sendMessage(chatID, "Başarı seviyen kaç yıldız? (1-5 arası bir sayı yaz)")
		case "test":
			flow.StudyType = "test_coz"
			flow.Step = 6
			sendMessage(chatID, "Doğru ve yanlış sayılarını yaz (örn: 12 3)")
		case "genel":
			flow.StudyType = "genel_test"
			flow.Step = 6
			sendMessage(chatID, "Doğru ve yanlış sayılarını yaz (örn: 12 3)")
		default:
			sendMessage(chatID, "Geçersiz. goz, test veya genel yaz:")
		}

	case 5:
		stars, err := strconv.Atoi(text)
		if err != nil || stars < 1 || stars > 5 {
			sendMessage(chatID, "Lütfen 1 ile 5 arasında bir sayı yaz:")
			return
		}
		flow.Stars = stars
		flow.Step = 7
		saveStudyFromFlow(chatID, flow)

	case 6:
		parts := strings.Fields(text)
		if len(parts) < 2 {
			sendMessage(chatID, "Doğru ve yanlış sayılarını boşlukla ayırarak yaz (örn: 12 3):")
			return
		}
		correct, err1 := strconv.Atoi(parts[0])
		wrong, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || correct < 0 || wrong < 0 {
			sendMessage(chatID, "Geçersiz sayı. Doğru ve yanlış sayılarını yaz (örn: 12 3):")
			return
		}
		flow.Correct = correct
		flow.Wrong = wrong
		flow.Step = 7
		saveStudyFromFlow(chatID, flow)
	}
}

func saveStudyFromFlow(chatID int64, flow *StudyFlow) {
	today := time.Now().In(turkeyLoc).Format("2006-01-02")
	var net float64
	if flow.StudyType != "goz_gezdir" {
		net = float64(flow.Correct) - float64(flow.Wrong)*0.25
	}

	_, err := database.DB.Exec(`
		INSERT INTO study_sessions (user_id, subject, topic, study_type, stars, correct, wrong, net, study_date)
		SELECT id, ?, ?, ?, ?, ?, ?, ?, ?
		FROM users WHERE telegram_chat_id = ?
	`, flow.Subject, flow.Topic, flow.StudyType, flow.Stars, flow.Correct, flow.Wrong, net, today, chatID)

	if err != nil {
		sendMessage(chatID, "Kayıt hatası: "+err.Error())
		studyFlowsMu.Lock()
		delete(studyFlows, chatID)
		studyFlowsMu.Unlock()
		return
	}

	var userID int64
	database.DB.QueryRow(`SELECT id FROM users WHERE telegram_chat_id = ?`, chatID).Scan(&userID)
	if userID > 0 && flow.StudyType == "goz_gezdir" {
		createSpacedRepetitionTodosFromBot(userID, flow.Subject, flow.Topic)
	}

	studyFlowsMu.Lock()
	delete(studyFlows, chatID)
	studyFlowsMu.Unlock()

	msg := fmt.Sprintf("✓ Kaydedildi: %s > %s (%s)\n", flow.Subject, flow.Topic, flow.StudyType)
	if flow.StudyType == "goz_gezdir" {
		msg += fmt.Sprintf("⭐ %d yıldız\n+1, +7, +30 gün için tekrar görevleri oluşturuldu.", flow.Stars)
	} else {
		msg += fmt.Sprintf("%d doğru, %d yanlış = %.2f net", flow.Correct, flow.Wrong, net)
	}
	sendMessage(chatID, msg)
}

func sendTodayTodos(chatID int64) {
	today := time.Now().In(turkeyLoc).Format("2006-01-02")
	rows, err := database.DB.Query(`
		SELECT t.id, t.subject, t.topic, t.todo_type FROM todos t
		JOIN users u ON u.id = t.user_id
		WHERE u.telegram_chat_id = ? AND t.due_date <= ? AND t.completed = 0
		ORDER BY t.due_date ASC
	`, chatID, today)
	if err != nil {
		sendMessage(chatID, "Bir hata oluştu.")
		return
	}
	defer rows.Close()

	labelMap := map[string]string{"goz_gezdir": "Göz Gezdir", "test_coz": "Test Çöz", "genel_test": "Genel Test"}
	var msg string
	count := 0
	for rows.Next() {
		var id int64
		var subject, topic, todoType string
		if err := rows.Scan(&id, &subject, &topic, &todoType); err != nil {
			continue
		}
		count++
		label := labelMap[todoType]
		if label == "" {
			label = todoType
		}
		msg += fmt.Sprintf("[%d] %s - %s (%s)\n", id, subject, topic, label)
	}

	if count == 0 {
		sendMessage(chatID, "Bugün için yapılması gereken görev bulunmuyor. Tebrikler!")
		return
	}

	sendMessage(chatID, "Bugün yapman gerekenler:\n\n"+msg+"\nTamamlamak için: /tamamla ID")
}

func sendWeeklyTodos(chatID int64) {
	today := time.Now().In(turkeyLoc)
	weekLater := today.AddDate(0, 0, 7)
	todayStr := today.Format("2006-01-02")
	weekLaterStr := weekLater.Format("2006-01-02")

	rows, err := database.DB.Query(`
		SELECT t.id, t.subject, t.topic, t.todo_type, t.due_date, t.completed
		FROM todos t
		JOIN users u ON u.id = t.user_id
		WHERE u.telegram_chat_id = ? AND t.due_date >= ? AND t.due_date <= ?
		ORDER BY t.due_date ASC, t.completed ASC
	`, chatID, todayStr, weekLaterStr)
	if err != nil {
		sendMessage(chatID, "Bir hata oluştu.")
		return
	}
	defer rows.Close()

	type Task struct {
		ID        int64
		Subject   string
		Topic     string
		TodoType  string
		DueDate   string
		Completed bool
	}

	tasksByDate := make(map[string][]Task)
	var dateOrder []string

	for rows.Next() {
		var t Task
		var completed int
		if err := rows.Scan(&t.ID, &t.Subject, &t.Topic, &t.TodoType, &t.DueDate, &completed); err != nil {
			continue
		}
		t.Completed = completed == 1
		if _, exists := tasksByDate[t.DueDate]; !exists {
			dateOrder = append(dateOrder, t.DueDate)
		}
		tasksByDate[t.DueDate] = append(tasksByDate[t.DueDate], t)
	}

	if len(dateOrder) == 0 {
		sendMessage(chatID, "Önümüzdeki 7 günde görev bulunmuyor. Tebrikler!")
		return
	}

	labelMap := map[string]string{"goz_gezdir": "Göz Gezdir", "test_coz": "Test Çöz", "genel_test": "Genel Test"}
	dayNames := map[string]string{
		"Monday": "Pazartesi", "Tuesday": "Salı", "Wednesday": "Çarşamba",
		"Thursday": "Perşembe", "Friday": "Cuma", "Saturday": "Cumartesi", "Sunday": "Pazar",
	}

	msg := "📅 Haftalık Takvim\n\n"
	for _, date := range dateOrder {
		t, _ := time.Parse("2006-01-02", date)
		dayName := dayNames[t.Weekday().String()]
		msg += fmt.Sprintf("── %s %s ──\n", dayName, date)
		for _, task := range tasksByDate[date] {
			label := labelMap[task.TodoType]
			if label == "" {
				label = task.TodoType
			}
			status := ""
			if task.Completed {
				status = " ✓"
			}
			msg += fmt.Sprintf("[%d] %s - %s (%s)%s\n", task.ID, task.Subject, task.Topic, label, status)
		}
		msg += "\n"
	}
	msg += "Tamamlamak için: /tamamla ID"

	sendMessage(chatID, msg)
}

func completeTodoTelegram(chatID int64, todoID int64) {
	result, err := database.DB.Exec(`
		UPDATE todos SET completed = 1
		WHERE id = ? AND user_id = (SELECT id FROM users WHERE telegram_chat_id = ?)
	`, todoID, chatID)
	if err != nil {
		sendMessage(chatID, "Bir hata oluştu.")
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		sendMessage(chatID, "Görev bulunamadı veya zaten tamamlanmış.")
		return
	}
	sendMessage(chatID, "✓ Görev tamamlandı!")
}

func createSpacedRepetitionTodosFromBot(userID int64, subject, topic string) {
	now := time.Now().In(turkeyLoc)
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
			`INSERT INTO todos (user_id, subject, topic, todo_type, due_date) VALUES (?, ?, ?, ?, ?)`,
			userID, subject, topic, iv.typ, dueDate,
		)
	}
}

func activateUser(chatID int64, code string) {
	result, err := database.DB.Exec(
		`UPDATE users SET telegram_chat_id = ? WHERE activation_code = ? AND telegram_chat_id = 0`,
		chatID, code,
	)
	if err != nil {
		sendMessage(chatID, "Bir hata oluştu, tekrar dene.")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		sendMessage(chatID, "Geçersiz veya daha önce kullanılmış aktivasyon kodu.")
		return
	}

	sendMessage(chatID, "Tebrikler! Bot başarıyla aktifleştirildi. Artık günlük ders hatırlatmaları alacaksın.\nBağlantını kesmek için /botbaglantikes yazabilirsin.")
}

func deactivateUser(chatID int64) {
	result, err := database.DB.Exec(
		`UPDATE users SET telegram_chat_id = 0, activation_code = '' WHERE telegram_chat_id = ?`,
		chatID,
	)
	if err != nil {
		sendMessage(chatID, "Bir hata oluştu.")
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		sendMessage(chatID, "Zaten aktif bir bağlantın bulunmuyor.")
		return
	}
	sendMessage(chatID, "Bağlantın kesildi. Artık bildirim almayacaksın.\nTekrar bağlanmak için web panelden yeni kod alıp /aktifkod KOD yazabilirsin.")
}

func sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	Bot.Send(msg)
}

func startDailyReminder() {
	for {
		now := time.Now().In(turkeyLoc)
		today := now.Format("2006-01-02")

		rows, err := database.DB.Query(`
			SELECT id, telegram_chat_id, reminder_hour, reminder_minute, last_reminder_date
			FROM users WHERE telegram_chat_id > 0 AND active = 1
		`)
		if err == nil {
			for rows.Next() {
				var userID, chatID int64
				var hour, minute int
				var lastDate string
				if err := rows.Scan(&userID, &chatID, &hour, &minute, &lastDate); err != nil {
					continue
				}
				if now.Hour() == hour && now.Minute() == minute && lastDate != today {
					sendDailyTodosToUser(chatID)
					database.DB.Exec(`UPDATE users SET last_reminder_date = ? WHERE id = ?`, today, userID)
				}
			}
			rows.Close()
		}

		time.Sleep(30 * time.Second)
	}
}

func sendDailyTodosToUser(chatID int64) {
	today := time.Now().In(turkeyLoc).Format("2006-01-02")

	trows, err := database.DB.Query(`
		SELECT t.id, t.subject, t.topic, t.todo_type FROM todos t
		JOIN users u ON u.id = t.user_id
		WHERE u.telegram_chat_id = ? AND t.due_date <= ? AND t.completed = 0
		ORDER BY t.due_date ASC
	`, chatID, today)
	if err != nil {
		log.Printf("[Bot] Günlük görev sorgu hatası: %v", err)
		return
	}
	defer trows.Close()

	type Task struct {
		ID       int64
		Subject  string
		Topic    string
		TodoType string
	}

	var tasks []Task
	for trows.Next() {
		var t Task
		if err := trows.Scan(&t.ID, &t.Subject, &t.Topic, &t.TodoType); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	if len(tasks) == 0 {
		return
	}

	labelMap := map[string]string{
		"goz_gezdir": "Göz Gezdir",
		"test_coz":   "Test Çöz",
		"genel_test": "Genel Test",
	}

	msg := "Günaydın! Bugün yapman gerekenler:\n\n"
	for _, t := range tasks {
		label := labelMap[t.TodoType]
		if label == "" {
			label = t.TodoType
		}
		msg += fmt.Sprintf("[%d] %s - %s (%s)\n", t.ID, t.Subject, t.Topic, label)
	}
	msg += "\nTamamlamak için: /tamamla ID"
	sendMessage(chatID, msg)
}
