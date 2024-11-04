// controllers/publisher_controller.go
package controllers

import (
    "net/http"
    "strconv"
    "auth-user-api/models"
    "auth-user-api/services"
    "auth-user-api/domains"
    "github.com/labstack/echo/v4"
)

type PublisherController struct {
    service services.PublisherService
}

func NewPublisherController(service services.PublisherService) *PublisherController {
    return &PublisherController{service}
}

// CreatePublisher handles creating a new publisher
func (c *PublisherController) CreatePublisher(ctx echo.Context) error {
    publisher := new(models.Publisher)
    if err := ctx.Bind(publisher); err != nil {
        response := domains.NewErrorResponse("400", "Invalid input", err.Error())
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := c.service.CreatePublisher(publisher); err != nil {
        response := domains.NewErrorResponse("500", "Failed to create publisher", err.Error())
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    data := domains.PublisherResponse{
        ID:        publisher.ID,
        Name:      publisher.Name,
        CreatedAt: publisher.CreatedAt.String(),
        UpdatedAt: publisher.UpdatedAt.String(),
        DeletedAt: nil,
    }
    response := domains.NewSuccessResponseWithData("200", "Publisher created successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// GetPublisherByID retrieves a publisher by their ID
func (c *PublisherController) GetPublisherByID(ctx echo.Context) error {
    id, _ := strconv.Atoi(ctx.Param("id"))
    publisher, err := c.service.GetPublisherByID(id)
    if err != nil {
        response := domains.NewErrorResponse("404", "Publisher not found", err.Error())
        return ctx.JSON(http.StatusNotFound, response)
    }

    data := domains.PublisherResponse{
        ID:        publisher.ID,
        Name:      publisher.Name,
        CreatedAt: publisher.CreatedAt.String(),
        UpdatedAt: publisher.UpdatedAt.String(),
        DeletedAt: nil,
    }
    response := domains.NewSuccessResponseWithData("200", "Publisher retrieved successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// GetAllPublishers retrieves all publishers
func (c *PublisherController) GetAllPublishers(ctx echo.Context) error {
    publishers, err := c.service.GetAllPublishers()
    if err != nil {
        response := domains.NewErrorResponse("500", "Failed to retrieve publishers", err.Error())
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    publisherData := make([]domains.PublisherResponse, len(publishers))
    for i, publisher := range publishers {
        publisherData[i] = domains.PublisherResponse{
            ID:        publisher.ID,
            Name:      publisher.Name,
            CreatedAt: publisher.CreatedAt.String(),
            UpdatedAt: publisher.UpdatedAt.String(),
            DeletedAt: nil,
        }
    }
    response := domains.NewSuccessResponseWithData("200", "Publishers retrieved successfully", publisherData)
    return ctx.JSON(http.StatusOK, response)
}

// UpdatePublisher handles updating a publisher's details
func (c *PublisherController) UpdatePublisher(ctx echo.Context) error {
    id, _ := strconv.Atoi(ctx.Param("id"))
    publisher, err := c.service.GetPublisherByID(id)
    if err != nil {
        response := domains.NewErrorResponse("404", "Publisher not found", err.Error())
        return ctx.JSON(http.StatusNotFound, response)
    }

    if err := ctx.Bind(publisher); err != nil {
        response := domains.NewErrorResponse("400", "Invalid input", err.Error())
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := c.service.UpdatePublisher(publisher); err != nil {
        response := domains.NewErrorResponse("500", "Failed to update publisher", err.Error())
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    data := domains.PublisherResponse{
        ID:        publisher.ID,
        Name:      publisher.Name,
        CreatedAt: publisher.CreatedAt.String(),
        UpdatedAt: publisher.UpdatedAt.String(),
        DeletedAt: nil,
    }
    response := domains.NewSuccessResponseWithData("200", "Publisher updated successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// DeletePublisher handles deleting a publisher by ID
func (c *PublisherController) DeletePublisher(ctx echo.Context) error {
    id, _ := strconv.Atoi(ctx.Param("id"))
    if err := c.service.DeletePublisher(id); err != nil {
        response := domains.NewErrorResponse("404", "Failed to delete publisher", err.Error())
        return ctx.JSON(http.StatusNotFound, response)
    }

    data := domains.DeleteResponse{UserID: strconv.Itoa(id)}
    response := domains.NewSuccessResponseWithData("200", "Publisher deleted successfully", data)
    return ctx.JSON(http.StatusOK, response)
}
