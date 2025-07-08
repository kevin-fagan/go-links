run:
	go run ./cmd/server/main.go

setup:
	go mod download
	sqlite3 golinks.db < setup.sql

