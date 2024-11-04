// models/user.go

package models

import (
    "time"

    "gorm.io/gorm"
)

type User struct {
    ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    Username  string         `gorm:"unique;not null" json:"username"`
    Email     string         `gorm:"unique;not null" json:"email"`
    Password  string         `gorm:"not null" json:"-"`
    Role      int            `gorm:"not null;default:2" json:"role"` // 1 untuk admin, 2 untuk member
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
