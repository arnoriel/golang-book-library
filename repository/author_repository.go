// repository/author_repository.go
package repository

import (
    "errors"
    "auth-user-api/models"
    "gorm.io/gorm"
)

type AuthorRepository interface {
    CreateAuthor(author *models.Author) error
    GetAuthorByID(id int) (*models.Author, error)
    GetAllAuthors() ([]*models.Author, error)
    UpdateAuthor(author *models.Author) error
    DeleteAuthor(id int) error
}

type authorRepository struct {
    db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
    return &authorRepository{db}
}

func (r *authorRepository) CreateAuthor(author *models.Author) error {
    return r.db.Create(author).Error
}

func (r *authorRepository) GetAuthorByID(id int) (*models.Author, error) {
    var author models.Author
    if err := r.db.First(&author, id).Error; err != nil {
        return nil, err
    }
    return &author, nil
}

func (r *authorRepository) GetAllAuthors() ([]*models.Author, error) {
    var authors []models.Author
    if err := r.db.Find(&authors).Error; err != nil {
        return nil, err
    }

    var authorPointers []*models.Author
    for i := range authors {
        authorPointers = append(authorPointers, &authors[i])
    }

    return authorPointers, nil
}

func (r *authorRepository) UpdateAuthor(author *models.Author) error {
    return r.db.Save(author).Error
}

func (r *authorRepository) DeleteAuthor(id int) error {
    result := r.db.Delete(&models.Author{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("author not found")
    }
    return nil
}
