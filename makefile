postgres:
	docker run --name eventDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it eventDB createdb --username=root --owner=root events
dropdb:
	docker exec -it eventDB dropdb events
migrateup:
	migrate -path internal/datastore/migrations -database "postgresql://root:secret@localhost:5432/events?sslmode=disable" -verbose up
migratedown:
	migrate -path internal/datastore/migrations -database "postgresql://root:secret@localhost:5432/events?sslmode=disable" -verbose down
sqlc:
	sqlc generate

PROTO_DIR=pkg/proto/event
gproto:
	protoc \
	--proto_path=. \
	--go_out=. \
	--go_opt=paths=source_relative \
	${PROTO_DIR}/*.proto \
	 --go-grpc_out=. \
	 --go-grpc_opt=paths=source_relative \
	 ${PROTO_DIR}/*.proto
ghttp:
	protoc -I . \
	--grpc-gateway_out \
	. \
	--grpc-gateway_opt \
	logtostderr=true \
	--grpc-gateway_opt \
	paths=source_relative \
	 ${PROTO_DIR}/*.proto

run:
	HOST=localhost \
	PORT=9090 \
	DB_DRIVER=postgres \
	DB_DATA_SOURCE_NAME=postgresql://root:secret@localhost:5432/events?sslmode=disable \
	GRPC_PORT=50001 \
	GRPC_HOST=localhost \
	go run .