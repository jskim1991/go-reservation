start:
	go run cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go

build:
	go build -o bookings cmd/web/*.go && ./bookings