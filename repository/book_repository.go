// repository/book_repository.go
package repository

import (
    "errors"
    "auth-user-api/models"

    "gorm.io/gorm"
)

type BookRepository interface {
    CreateBook(book *models.Book) error
    GetBookByID(id int) (*models.Book, error)
    GetAllBooks() ([]*models.Book, error)
    UpdateBook(book *models.Book) error
    DeleteBook(id int) error
}

type bookRepository struct {
    db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
    return &bookRepository{db}
}

func (r *bookRepository) CreateBook(book *models.Book) error {
    if book.Stock > book.MaxStock {
        return errors.New("stock cannot exceed max_stock")
    }
    return r.db.Create(book).Error
}

func (r *bookRepository) GetBookByID(id int) (*models.Book, error) {
    var book models.Book
    if err := r.db.Preload("Author").Preload("Publisher").First(&book, id).Error; err != nil {
        return nil, err
    }
    return &book, nil
}

func (r *bookRepository) GetAllBooks() ([]*models.Book, error) {
    var books []models.Book
    if err := r.db.Preload("Author").Preload("Publisher").Find(&books).Error; err != nil {
        return nil, err
    }

    var bookPointers []*models.Book
    for i := range books {
        bookPointers = append(bookPointers, &books[i])
    }

    return bookPointers, nil
}

func (r *bookRepository) UpdateBook(book *models.Book) error {
    if book.Stock > book.MaxStock {
        return errors.New("stock cannot exceed max_stock")
    }
    return r.db.Save(book).Error
}

func (r *bookRepository) DeleteBook(id int) error {
    result := r.db.Delete(&models.Book{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("book not found")
    }
    return nil
}
