run:
	# npx @tailwindcss/cli -i ./web/assets/css/inputs.css -o ./web/assets/css/styles.css --watch
	go run ./cmd/server/main.go

setup:
	go mod download
	sqlite3 golinks-db < setup.sql

