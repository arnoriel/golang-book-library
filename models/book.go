// models/book.go
package models

import "gorm.io/gorm"

type Book struct {
    ID          int     `gorm:"primaryKey" json:"id"`
    Title       string  `json:"title"`
    AuthorID    int     `json:"author_id"`
    PublisherID int     `json:"publisher_id"`
    Summary     string  `json:"summary"`
    Stock       int     `json:"stock"`
    MaxStock    int     `json:"max_stock"`
    Author      Author  `gorm:"foreignKey:AuthorID"`
    Publisher   Publisher `gorm:"foreignKey:PublisherID"`
    gorm.Model
}
