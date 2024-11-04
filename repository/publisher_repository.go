package repository

import (
    "errors"
    "auth-user-api/models"
    "gorm.io/gorm"
)

type PublisherRepository interface {
    CreatePublisher(publisher *models.Publisher) error
    GetPublisherByID(id int) (*models.Publisher, error)
    GetAllPublishers() ([]*models.Publisher, error)
    UpdatePublisher(publisher *models.Publisher) error
    DeletePublisher(id int) error
}

type publisherRepository struct {
    db *gorm.DB
}

func NewPublisherRepository(db *gorm.DB) PublisherRepository {
    return &publisherRepository{db}
}

func (r *publisherRepository) CreatePublisher(publisher *models.Publisher) error {
    return r.db.Create(publisher).Error
}

func (r *publisherRepository) GetPublisherByID(id int) (*models.Publisher, error) {
    var publisher models.Publisher
    if err := r.db.First(&publisher, id).Error; err != nil {
        return nil, err
    }
    return &publisher, nil
}

func (r *publisherRepository) GetAllPublishers() ([]*models.Publisher, error) {
    var publishers []models.Publisher
    if err := r.db.Find(&publishers).Error; err != nil {
        return nil, err
    }

    var publisherPointers []*models.Publisher
    for i := range publishers {
        publisherPointers = append(publisherPointers, &publishers[i])
    }

    return publisherPointers, nil
}

func (r *publisherRepository) UpdatePublisher(publisher *models.Publisher) error {
    return r.db.Save(publisher).Error
}

func (r *publisherRepository) DeletePublisher(id int) error {
    result := r.db.Delete(&models.Publisher{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("publisher not found")
    }
    return nil
}
