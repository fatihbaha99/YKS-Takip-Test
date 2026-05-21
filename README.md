# YKS Takip

Go Fiber + SQLite ile yazılmış, Telegram bot destekli YKS çalışma takip sistemi.

## Özellikler

- TYT/AYT konu bazlı çalışma kaydı
- Aralıklı tekrar (spaced repetition) sistemi
- Telegram bot ile çalışma kaydı ve sorgulama
- Günlük/haftalık takvim görünümü
- Kullanıcı yönetimi ve JWT kimlik doğrulama
- Gece temizliği (1 yıl geçmiş veriler -> Excel -> e-posta -> silme)

## Kurulum

1. `.env` dosyası oluşturun (`.env.example` içindeki alanları doldurun)
2. `go build -o yks-tracker .`
3. Binary'yi çalıştırın

## Teknolojiler

- **Backend:** Go + Fiber
- **Veritabanı:** SQLite (modernc.org/sqlite)
- **Bot:** Telegram Bot API (polling modu)
- **Frontend:** Pure HTML/CSS/JS

## Lisans

MIT
