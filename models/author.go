package models

import "gorm.io/gorm"

type Author struct {
    ID        int    `gorm:"primaryKey" json:"id"`
    Name      string `json:"name"`
    gorm.Model
}
