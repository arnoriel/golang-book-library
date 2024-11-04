// models/loan_request.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type LoanRequest struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    BookID       int       `gorm:"not null" json:"book_id"`
    UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    RequestTime  time.Time `json:"request_time"`
    Status       string    `json:"status"` // "PENDING", "APPROVED", "REJECTED"
    RejectReason *string   `json:"reject_reason,omitempty"`
}
