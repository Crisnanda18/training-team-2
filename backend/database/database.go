package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// VULNERABILITY #4: Hardcoded database credentials
func Connect() {
	// These credentials should come from environment variables!
	dsn := "host=localhost user=taskuser password=taskpass123 dbname=securetask port=5432 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Database connected successfully")
}

// VULNERABILITY #1: Raw SQL execution without parameterization
func ExecuteRawSQL(query string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// This allows SQL injection!
	rows, err := DB.Raw(query).Rows()
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
