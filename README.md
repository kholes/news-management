# News & Topic Management API

RESTful API untuk manajemen **News** dan **Topic**, dibangun dengan **Golang (Echo + Gorm)** dan **PostgreSQL**.  
Mendukung dokumentasi API otomatis menggunakan **Swagger** dan siap dijalankan menggunakan **Docker Compose**.

## A. Fitur
- CRUD **Topic** (kategori berita).
- CRUD **News** dengan relasi ke Topic.
- API Documentation dengan Swagger.
- Unit testing dengan Go.
- CI/CD pipeline otomatis untuk build & test.
- Makefile untuk mempermudah perintah sehari-hari.


## B. Prasyarat

Pastikan sudah terinstall:
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/) *(opsional, tapi disarankan)*


## C. Memulai Project

### 1. Clone Repository
```bash
git clone https://github.com/kholes/news-management.git
cd news-management
```

### 2. Konfigurasi Environment
Buat file .env di root project:
```bash
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=news_db
APP_PORT=8080
```

## D. Menjalankan Project
Untuk menjalankan project bisa menggunakan beberapa cara:
### 1. Menggunakan Docker Compose
Build dan jalankan semua service:
```bash
docker-compose up --build
```
Hentikan service:
```bash
docker-compose down
```
### 2. Dengan Makefile
Jalankan build dengan Docker Compose:
```bash
make build
```
Jalankan dengan Docker Compose:
```bash
make run
```
Stop service:
```bash
make down
```
Jalankan test:
```bash
make test
```
## E. Dokumentasi API
Swagger docs tersedia di endpoint:
```bash
http://localhost:8080/docs/index.html
```

## F. Testing
Testing API bisa menggunakan swagger (http://localhost:8080/docs/index.html) dan make command:
```bash
make test
``` 

## G. Struktur Project
```bash
.
├── internal/
│   ├── config/
│   ├── models/
│   ├── handlers/
│   ├── routes/
│   └── database/
├── docs/              # swagger docs
├── main.go
├── .env
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── ci.yml
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## H. Makefile
```bash
# Generate swagger docs
swag:
	swag init -g main.go -o docs

# Run app
run:
	docker-compose up -d

logs:
	docker logs -f --tail 100 app

down: 
	docker-compose down

restart:
	docker-compose restart

# Build binary
build:
	docker-compose up --build -d

# Run unit tests
test:
	docker-compose exec app go test ./... -v
```