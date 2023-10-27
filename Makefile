.PHONY: migrate migrate_dock

migrate_up:
	migrate -path db/migration -database "postgresql://emtest:emtest@localhost:3004/emtest?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://emtest:emtest@localhost:3004/emtest?sslmode=disable" -verbose down

migrate_dock_up:
	docker run -v $(P)/db/migration:/migrations --network emtest_emtest migrate/migrate -path=/migrations/ -database 'postgres://emtest:emtest@postgres:5432/emtest?sslmode=disable' up

migrate_dock_down:
	docker run -v $(P)/db/migration:/migrations --network emtest_emtest migrate/migrate -path=/migrations/ -database 'postgres://emtest:emtest@postgres:5432/emtest?sslmode=disable' down -all 