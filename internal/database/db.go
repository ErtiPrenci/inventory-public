package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB is the global instance
var DB *sql.DB

func InitDB(connStr string) {
	var err error
	// Use pgx as driver compatible with database/sql
	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Error opening connection to DB:", err)
	}

	// Keep few connections to not saturate Supabase (PgBouncer mode transactional)
	DB.SetMaxOpenConns(2)
	DB.SetMaxIdleConns(2)
	DB.SetConnMaxLifetime(30 * time.Minute)

	// Verify connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("Connection to DB successful")
}
