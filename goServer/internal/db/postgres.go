package db

import (
	"context"
	"goServer/internal/config"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)
}

func Connect(cfg config.Config) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), cfg.DB_URL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	return pool
}
