migrations:
	goose -dir db/migrations postgres "postgresql://trainingservice:root@127.0.0.1:5432/training_db?sslmode=disable" up
dropmigrations:
	goose -dir db/migrations postgres "postgresql://trainingservice:root@127.0.0.1:5432/training_db?sslmode=disable" down
run:
	go build -o training && ./training