// models/loan_record.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type LoanRecord struct {
    ID         uint       `gorm:"primaryKey" json:"id"`
    BookID     int        `gorm:"not null" json:"book_id"`
    UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
    LoanDate   time.Time  `json:"loan_date"`
    DueDate    time.Time  `json:"due_date"`
    Returned   bool       `json:"returned"`
    ReturnDate *time.Time `json:"return_date,omitempty"`
}
