// services/author_services.go
package services

import (
    "auth-user-api/models"
    "auth-user-api/repository"
)

type AuthorService interface {
    CreateAuthor(author *models.Author) error
    GetAuthorByID(id int) (*models.Author, error)
    GetAllAuthors() ([]*models.Author, error)
    UpdateAuthor(author *models.Author) error
    DeleteAuthor(id int) error
}

type authorService struct {
    repo repository.AuthorRepository
}

func NewAuthorService(repo repository.AuthorRepository) AuthorService {
    return &authorService{repo}
}

func (s *authorService) CreateAuthor(author *models.Author) error {
    return s.repo.CreateAuthor(author)
}

func (s *authorService) GetAuthorByID(id int) (*models.Author, error) {
    return s.repo.GetAuthorByID(id)
}

func (s *authorService) GetAllAuthors() ([]*models.Author, error) {
    return s.repo.GetAllAuthors()
}

func (s *authorService) UpdateAuthor(author *models.Author) error {
    return s.repo.UpdateAuthor(author)
}

func (s *authorService) DeleteAuthor(id int) error {
    return s.repo.DeleteAuthor(id)
}
