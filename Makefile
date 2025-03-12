MAIN_PATH = cmd/seelochka/main.go
CONFIG_PATH = configs/conf.yaml
BINARY_NAME = main.out
DB_PATH = db/sqlite.db
MIGRATIONS_PATH = ./db/migrations/
SWAGGER_DIRS = cmd/seelochka/,internal/http/handlers/urls/

build:
	CONFIG_PATH=${CONFIG_PATH} go build -o ${BINARY_NAME} ${MAIN_PATH}

run:
	CONFIG_PATH=${CONFIG_PATH} go run ${MAIN_PATH}

clean:
	go clean
	rm -f ${BINARY_NAME}

migrate:
	~/go/bin/migrate -database sqlite3://${DB_PATH} -path ${MIGRATIONS_PATH} $(CMD)

swagger:
	~/go/bin/swag init -d ${SWAGGER_DIRS}

.PHONY: build run clean migrate swagger
