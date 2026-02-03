package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Database connection string from .env
	connStr := "postgres://go-ai:ipJK6TmtL8@157.66.218.138:5432/go-ai-db?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// SQL statements to drop all tables except users, roles, and schema_migrations
	sqlStatements := []string{
		"DROP TABLE IF EXISTS combo_group_item CASCADE",
		"DROP TABLE IF EXISTS combo_group_items CASCADE",
		"DROP TABLE IF EXISTS combo_group CASCADE",
		"DROP TABLE IF EXISTS combo_groups CASCADE",
		"DROP TABLE IF EXISTS menu_item_option_group CASCADE",
		"DROP TABLE IF EXISTS menu_item_option_groups CASCADE",
		"DROP TABLE IF EXISTS menu_item_option CASCADE",
		"DROP TABLE IF EXISTS menu_item_options CASCADE",
		"DROP TABLE IF EXISTS option_item CASCADE",
		"DROP TABLE IF EXISTS option_items CASCADE",
		"DROP TABLE IF EXISTS option_group CASCADE",
		"DROP TABLE IF EXISTS option_groups CASCADE",
		"DROP TABLE IF EXISTS menu_item_variant CASCADE",
		"DROP TABLE IF EXISTS menu_item_variants CASCADE",
		"DROP TABLE IF EXISTS menu_item CASCADE",
		"DROP TABLE IF EXISTS menu_items CASCADE",
		"DROP TABLE IF EXISTS topic CASCADE",
		"DROP TABLE IF EXISTS topics CASCADE",
		"DROP TABLE IF EXISTS restaurant_users CASCADE",
		"DROP TABLE IF EXISTS restaurant_hours CASCADE",
		"DROP TABLE IF EXISTS restaurant CASCADE",
		"DROP TABLE IF EXISTS restaurants CASCADE",
		"DROP TYPE IF EXISTS menu_item_type CASCADE",
	}

	fmt.Println("Starting database cleanup...")

	for _, sql := range sqlStatements {
		fmt.Printf("Executing: %s\n", sql)
		_, err := pool.Exec(context.Background(), sql)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		} else {
			fmt.Println("✓ Success")
		}
	}

	// List remaining tables
	fmt.Println("\n=== Remaining tables ===")
	rows, err := pool.Query(context.Background(), `
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
		ORDER BY tablename
	`)
	if err != nil {
		log.Fatalf("Error listing tables: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tablename string
		if err := rows.Scan(&tablename); err != nil {
			log.Fatalf("Error scanning: %v\n", err)
		}
		fmt.Printf("  - %s\n", tablename)
	}

	fmt.Println("\nCleanup complete!")
}
