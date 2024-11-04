package services

import (
    "auth-user-api/models"
    "auth-user-api/repository"
)

type PublisherService interface {
    CreatePublisher(publisher *models.Publisher) error
    GetPublisherByID(id int) (*models.Publisher, error)
    GetAllPublishers() ([]*models.Publisher, error)
    UpdatePublisher(publisher *models.Publisher) error
    DeletePublisher(id int) error
}

type publisherService struct {
    repo repository.PublisherRepository
}

func NewPublisherService(repo repository.PublisherRepository) PublisherService {
    return &publisherService{repo}
}

func (s *publisherService) CreatePublisher(publisher *models.Publisher) error {
    return s.repo.CreatePublisher(publisher)
}

func (s *publisherService) GetPublisherByID(id int) (*models.Publisher, error) {
    return s.repo.GetPublisherByID(id)
}

func (s *publisherService) GetAllPublishers() ([]*models.Publisher, error) {
    return s.repo.GetAllPublishers()
}

func (s *publisherService) UpdatePublisher(publisher *models.Publisher) error {
    return s.repo.UpdatePublisher(publisher)
}

func (s *publisherService) DeletePublisher(id int) error {
    return s.repo.DeletePublisher(id)
}
