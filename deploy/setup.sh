#!/bin/bash
# Ubuntu sunucu kurulum scripti

# Dizinleri oluştur
mkdir -p /var/db
mkdir -p /var/www/yks-frontend

# Binary ve frontend dosyalarını kopyala
cp /tmp/yks-tracker /usr/local/bin/yks-tracker
chmod +x /usr/local/bin/yks-tracker

cp /tmp/frontend/index.html /var/www/yks-frontend/
cp /tmp/frontend/css/style.css /var/www/yks-frontend/css/
cp /tmp/frontend/js/*.js /var/www/yks-frontend/js/

# Systemd servisini kur
cp /tmp/deploy/yks-backend.service /etc/systemd/system/yks-backend.service

# Servisi başlat
systemctl daemon-reload
systemctl enable yks-backend
systemctl start yks-backend

# Güvenlik duvarı
ufw allow 4000/tcp
ufw allow 4080/tcp

echo "Kurulum tamamlandı!"
