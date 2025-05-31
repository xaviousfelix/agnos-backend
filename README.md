# Agnos Hospital Middleware System 🏥

ระบบกลางสำหรับจัดการโรงพยาบาล รองรับ API การล็อกอินบุคลากร การจัดเก็บผู้ป่วย และระบบยืนยันตัวตนด้วย JWT

## 🚀 Tech Stack
- Go (Gin Framework)
- PostgreSQL
- Docker & Docker Compose
- Nginx

## 📁 โครงสร้างโปรเจกต์
- `controllers/` - จัดการ HTTP Request
- `models/` - โครงสร้างข้อมูล GORM
- `routes/` - เส้นทาง API
- `services/` - ธุรกิจ เช่น JWT
- `middleware/` - Auth Middleware
- `config/` - .env และ DB
- `tests/` - Unit Test

## 🐳 การใช้งานด้วย Docker
```bash
docker-compose up --build
