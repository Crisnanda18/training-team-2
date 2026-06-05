package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// VULNERABILITY #4: Hardcoded database credentials
/*
	Fix: Use environment variables for database credentials.
	How: Replace hardcoded connection string with environment variables and update Connect function to read from them.
*/
func Connect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Database connected successfully")
}

// VULNERABILITY #1: Raw SQL execution without parameterization
/*
	Fix: Use parameterized Queries.
	How: Add arguments as parameters to the ExecuteRawSQL function and pass them to DB.
*/
func ExecuteRawSQL(query string, args ...any) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := DB.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		entry := make(map[string]interface{})
		for i, col := range columns {
			entry[col] = values[i]
		}

		results = append(results, entry)
	}

	return results, nil
}
