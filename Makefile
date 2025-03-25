MAIN_PATH = cmd/seelochka/main.go
BINARY_NAME = server
DB_PATH = db/sqlite.db
MIGRATIONS_PATH = ./db/migrations/
SWAGGER_DIRS = cmd/seelochka/,internal/http/handlers/url/save/,internal/http/handlers/url/redirect/

build:
	go build -o ${BINARY_NAME} ${MAIN_PATH}

run:
	go run ${MAIN_PATH}

run_release:
	./server

clean:
	go clean
	rm -f ${BINARY_NAME}

migrate:
	~/go/bin/migrate -database sqlite3://${DB_PATH} -path ${MIGRATIONS_PATH} $(CMD)

swagger:
	~/go/bin/swag init -d ${SWAGGER_DIRS}

.PHONY: build run clean migrate swagger
