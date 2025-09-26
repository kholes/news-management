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
