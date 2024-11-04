// services/book_services.go
package services

import (
    "auth-user-api/models"
    "auth-user-api/repository"
)

type BookService interface {
    CreateBook(book *models.Book) error
    GetBookByID(id int) (*models.Book, error)
    GetAllBooks() ([]*models.Book, error)
    UpdateBook(book *models.Book) error
    DeleteBook(id int) error
}

type bookService struct {
    repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
    return &bookService{repo}
}

func (s *bookService) CreateBook(book *models.Book) error {
    return s.repo.CreateBook(book)
}

func (s *bookService) GetBookByID(id int) (*models.Book, error) {
    return s.repo.GetBookByID(id)
}

func (s *bookService) GetAllBooks() ([]*models.Book, error) {
    return s.repo.GetAllBooks()
}

func (s *bookService) UpdateBook(book *models.Book) error {
    return s.repo.UpdateBook(book)
}

func (s *bookService) DeleteBook(id int) error {
    return s.repo.DeleteBook(id)
}
