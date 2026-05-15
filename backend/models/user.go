package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password"` // VULNERABILITY #5: No hashing, plain text storage
	Name      string    `json:"name"`
	Role      string    `json:"role"` // admin or user
	Bio       string    `json:"bio"`  // VULNERABILITY #3: XSS if not sanitized
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VULNERABILITY #2: Returning password in JSON responses
// Should have `json:"-"` on Password field
