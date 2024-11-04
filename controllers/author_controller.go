package controllers

import (
    "net/http"
    "strconv"
    "auth-user-api/services"
    "auth-user-api/models"
    "auth-user-api/domains"
    "github.com/labstack/echo/v4"
)

type AuthorController struct {
    service services.AuthorService
}

func NewAuthorController(service services.AuthorService) *AuthorController {
    return &AuthorController{service}
}

// CreateAuthor handles creating a new author
func (c *AuthorController) CreateAuthor(ctx echo.Context) error {
    author := new(models.Author)
    if err := ctx.Bind(author); err != nil {
        response := domains.NewErrorResponse("400", "Invalid input", err.Error())
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := c.service.CreateAuthor(author); err != nil {
        response := domains.NewErrorResponse("500", "Failed to create author", err.Error())
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    data := domains.AuthorResponse{
        ID:        author.ID,
        Name:      author.Name,
        CreatedAt: author.CreatedAt.String(),
        UpdatedAt: author.UpdatedAt.String(),
        DeletedAt: nil,
    }
    response := domains.NewSuccessResponseWithData("200", "Author created successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// GetAuthorByID retrieves an author by their ID
func (c *AuthorController) GetAuthorByID(ctx echo.Context) error {
    id, _ := strconv.Atoi(ctx.Param("id"))
    author, err := c.service.GetAuthorByID(id)
    if err != nil {
        response := domains.NewErrorResponse("404", "Author not found", err.Error())
        return ctx.JSON(http.StatusNotFound, response)
    }

    data := domains.AuthorResponse{
        ID:        author.ID,
        Name:      author.Name,
        CreatedAt: author.CreatedAt.String(),
        UpdatedAt: author.UpdatedAt.String(),
        DeletedAt: nil,
    }
    response := domains.NewSuccessResponseWithData("200", "Author retrieved successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// GetAllAuthors retrieves all authors
func (c *AuthorController) GetAllAuthors(ctx echo.Context) error {
    authors, err := c.service.GetAllAuthors()
    if err != nil {
        response := domains.NewErrorResponse("500", "Failed to retrieve authors", err.Error())
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    authorData := make([]domains.AuthorResponse, len(authors))
    for i, author := range authors {
        authorData[i] = domains.AuthorResponse{
            ID:        author.ID,
            Name:      author.Name,
            CreatedAt: author.CreatedAt.String(),
            UpdatedAt: author.UpdatedAt.String(),
            DeletedAt: nil,
        }
    }
    response := domains.NewSuccessResponseWithData("200", "Authors retrieved successfully", authorData)
    return ctx.JSON(http.StatusOK, response)
}

// UpdateAuthor handles updating an author's details
func (c *AuthorController) UpdateAuthor(ctx echo.Context) error {
    id, _ := strconv.Atoi(ctx.Param("id"))
    author, err := c.service.GetAuthorByID(id)
    if err != nil {
        response := domains.NewErrorResponse("404", "Author not found", err.Error())
        return ctx.JSON(http.StatusNotFound, response)
    }

    if err := ctx.Bind(author); err != nil {
        response := domains.NewErrorResponse("400", "Invalid input", err.Error())
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := c.service.UpdateAuthor(author); err != nil {
        response := domains.NewErrorResponse("500", "Failed to update author", err.Error())
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    data := domains.AuthorResponse{
        ID:        author.ID,
        Name:      author.Name,
        CreatedAt: author.CreatedAt.String(),
        UpdatedAt: author.UpdatedAt.String(),
        DeletedAt: nil,
    }
    response := domains.NewSuccessResponseWithData("200", "Author updated successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// DeleteAuthor handles deleting an author by ID
func (c *AuthorController) DeleteAuthor(ctx echo.Context) error {
    id, _ := strconv.Atoi(ctx.Param("id"))
    if err := c.service.DeleteAuthor(id); err != nil {
        response := domains.NewErrorResponse("404", "Failed to delete author", err.Error())
        return ctx.JSON(http.StatusNotFound, response)
    }

    data := domains.DeleteResponse{UserID: strconv.Itoa(id)}
    response := domains.NewSuccessResponseWithData("200", "Author deleted successfully", data)
    return ctx.JSON(http.StatusOK, response)
}
