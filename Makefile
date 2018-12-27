protoc:
	docker run -v `pwd`:/defs namely/protoc-all -d proto/src -l go -o .

port-migrations-up:
	docker run -v `pwd`/src/service/port/storage/postgres/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgresql://postgres:@localhost:5432/postgres?sslmode=disable up

port-migrations-down:
	docker run -v `pwd`/src/service/port/storage/postgres/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgresql://postgres:@localhost:5432/postgres?sslmode=disable down
