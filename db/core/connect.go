package core

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func CreateConnection() *sqlx.DB {
	connectionString := createConnectionString()
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatal("Error opening database connection")
	}
	slog.Info("Created postgresql connection")
	return db
}

func createConnectionString() string {
	dbName, dbUser, dbPassword, dbHost, dbPort := readEnvVariables()
	result := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	return result
}

func readEnvVariables() (dbName, dbUser, dbPassword, dbHost, dbPort string) {
	if err := godotenv.Load(); err != nil {
		slog.Error("error reading environment variables")
	}
	return readEnvVariable("DATABASE_NAME"),
		readEnvVariable("DATABASE_USER"),
		readEnvVariable("DATABASE_PASSWORD"),
		readEnvVariable("DATABASE_HOST"),
		readEnvVariable("DATABASE_PORT")
}

func readEnvVariable(variableName string) string {
	result, exists := os.LookupEnv(variableName)
	if !exists {
		log.Fatalf("Can't load env variable: %v", variableName)
	}
	return result
}
