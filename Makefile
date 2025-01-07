migration_up: 
	migrate -path database/migration/ -database "postgresql://user:password@localhost:5432/dating?sslmode=disable" -verbose up

migration_down: 
	migrate -path database/migration/ -database "postgresql://user:password@localhost:5432/dating?sslmode=disable" -verbose down
