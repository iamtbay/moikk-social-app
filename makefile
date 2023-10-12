POSTGRES_URL=postgres://postgres:secret@localhost:5432/moikk_social_app?sslmode=disable
CONTAINER_NAME=server-db-1
build:
	@go build -o moikk ./cmd/*.go
run: build
	@./moikk
dbmigrateup:
	@migrate -database ${POSTGRES_URL} -path db/migrations up
dbmigratedown:
	@migrate -database ${POSTGRES_URL} -path db/migrations down
stopdocker:
	docker stop ${CONTAINER_NAME}
startdocker:
	docker start ${CONTAINER_NAME}

droppoststable:
	migrate -database ${POSTGRES_URL} -path db/posts down
	migrate -database ${POSTGRES_URL} -path db/posts up