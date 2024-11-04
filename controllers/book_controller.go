// controllers/book_controller.go

package controllers

import (
    "net/http"
    "strconv"
    "time"
    "auth-user-api/domains"
    "auth-user-api/models"
    "auth-user-api/services"
    "github.com/labstack/echo/v4"
)

type BookController struct {
    bookService      services.BookService
    authorService    services.AuthorService
    publisherService services.PublisherService
}

func NewBookController(bookService services.BookService, authorService services.AuthorService, publisherService services.PublisherService) *BookController {
    return &BookController{
        bookService:      bookService,
        authorService:    authorService,
        publisherService: publisherService,
    }
}

// Helper function to build BookResponse from a book model
func buildBookResponse(book *models.Book) domains.BookResponse {
    var deletedAt *string
    if book.DeletedAt.Valid {
        dt := book.DeletedAt.Time.Format(time.RFC3339)
        deletedAt = &dt
    }

    return domains.BookResponse{
        ID:          book.ID,
        Title:       book.Title,
        Summary:     book.Summary,
        AuthorID:    book.AuthorID,
        Author:      domains.BookAuthorResponse{
            Name: book.Author.Name,
        },
        PublisherID: book.PublisherID,
        Publisher:   domains.BookPublisherResponse{
            Name: book.Publisher.Name,
        },
        Stock:     book.Stock,
        MaxStock:  book.MaxStock,
        CreatedAt: book.CreatedAt.Format(time.RFC3339),
        UpdatedAt: book.UpdatedAt.Format(time.RFC3339),
        DeletedAt: deletedAt,
    }
}

// CreateBook with AuthorID and PublisherID validation
func (c *BookController) CreateBook(ctx echo.Context) error {
    book := new(models.Book)
    if err := ctx.Bind(book); err != nil {
        response := domains.NewErrorResponse("400", "Invalid input", "Binding error")
        return ctx.JSON(http.StatusBadRequest, response)
    }

    // Validate Author ID
    author, err := c.authorService.GetAuthorByID(book.AuthorID)
    if err != nil {
        response := domains.NewErrorResponse("400", "Invalid author ID", "Author not found")
        return ctx.JSON(http.StatusBadRequest, response)
    }
    book.Author = *author

    // Validate Publisher ID
    publisher, err := c.publisherService.GetPublisherByID(book.PublisherID)
    if err != nil {
        response := domains.NewErrorResponse("400", "Invalid publisher ID", "Publisher not found")
        return ctx.JSON(http.StatusBadRequest, response)
    }
    book.Publisher = *publisher

    // Create Book
    if err := c.bookService.CreateBook(book); err != nil {
        var code, message string
        if err.Error() == "stock cannot exceed max_stock" {
            code = "400"
            message = "Stock cannot exceed max_stock"
        } else {
            code = "500"
            message = "Failed to create book"
        }
        response := domains.NewErrorResponse(code, message, err.Error())
        return ctx.JSON(http.StatusBadRequest, response)
    }

    // Build and send success response
    data := buildBookResponse(book)
    response := domains.NewSuccessResponseWithData("200", "Book created successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// GetBookByID retrieves a book by ID
func (c *BookController) GetBookByID(ctx echo.Context) error {
    id, _ := strconv.Atoi(ctx.Param("id"))
    book, err := c.bookService.GetBookByID(id)
    if err != nil {
        response := domains.NewErrorResponse("404", "Book not found", "No book with specified ID")
        return ctx.JSON(http.StatusNotFound, response)
    }

    // Build and send success response
    data := buildBookResponse(book)
    response := domains.NewSuccessResponseWithData("200", "Book retrieved successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// GetAllBooks retrieves all books
func (c *BookController) GetAllBooks(ctx echo.Context) error {
    books, err := c.bookService.GetAllBooks()
    if err != nil {
        response := domains.NewErrorResponse("500", "Failed to retrieve books", err.Error())
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    // Prepare response data
    bookResponses := make([]domains.BookResponse, len(books))
    for i, book := range books {
        bookResponses[i] = buildBookResponse(book)
    }
    response := domains.NewSuccessResponseWithData("200", "Books retrieved successfully", bookResponses)
    return ctx.JSON(http.StatusOK, response)
}

// UpdateBook updates a book with optional AuthorID and PublisherID validation
func (c *BookController) UpdateBook(ctx echo.Context) error {
    id, _ := strconv.Atoi(ctx.Param("id"))
    book, err := c.bookService.GetBookByID(id)
    if err != nil {
        response := domains.NewErrorResponse("404", "Book not found", "No book with specified ID")
        return ctx.JSON(http.StatusNotFound, response)
    }

    // Temporary struct to hold the incoming update data
    var updateData struct {
        Title       *string `json:"title"`
        AuthorID    *int    `json:"author_id"`
        PublisherID *int    `json:"publisher_id"`
        Summary     *string `json:"summary"`
        Stock       *int    `json:"stock"`
        MaxStock    *int    `json:"max_stock"`
    }

    // Bind the incoming data
    if err := ctx.Bind(&updateData); err != nil {
        response := domains.NewErrorResponse("400", "Invalid input", "Binding error")
        return ctx.JSON(http.StatusBadRequest, response)
    }

    // Update fields if provided
    if updateData.Title != nil {
        book.Title = *updateData.Title
    }
    if updateData.Summary != nil {
        book.Summary = *updateData.Summary
    }
    if updateData.Stock != nil {
        book.Stock = *updateData.Stock
    }
    if updateData.MaxStock != nil {
        book.MaxStock = *updateData.MaxStock
    }

    // Validate and update AuthorID if provided
    if updateData.AuthorID != nil {
        author, err := c.authorService.GetAuthorByID(*updateData.AuthorID)
        if err != nil {
            response := domains.NewErrorResponse("400", "Invalid author ID", "Author not found")
            return ctx.JSON(http.StatusBadRequest, response)
        }
        book.AuthorID = *updateData.AuthorID
        book.Author = *author
    }

    // Validate and update PublisherID if provided
    if updateData.PublisherID != nil {
        publisher, err := c.publisherService.GetPublisherByID(*updateData.PublisherID)
        if err != nil {
            response := domains.NewErrorResponse("400", "Invalid publisher ID", "Publisher not found")
            return ctx.JSON(http.StatusBadRequest, response)
        }
        book.PublisherID = *updateData.PublisherID
        book.Publisher = *publisher
    }

    // Update the book
    if err := c.bookService.UpdateBook(book); err != nil {
        var code, message string
        if err.Error() == "stock cannot exceed max_stock" {
            code = "400"
            message = "Stock cannot exceed max_stock"
        } else {
            code = "500"
            message = "Failed to update book"
        }
        response := domains.NewErrorResponse(code, message, err.Error())
        return ctx.JSON(http.StatusBadRequest, response)
    }

    // Build and send success response
    data := buildBookResponse(book)
    response := domains.NewSuccessResponseWithData("200", "Book updated successfully", data)
    return ctx.JSON(http.StatusOK, response)
}

// DeleteBook deletes a book by ID
func (c *BookController) DeleteBook(ctx echo.Context) error {
    id, err := strconv.Atoi(ctx.Param("id"))
    response := domains.BaseResponse{
        Parameter: "id",
    }

    if err != nil {
        response.Code = strconv.Itoa(http.StatusBadRequest)
        response.Message = "Invalid book ID"
        response.Error = "ID conversion error"
        response.FormatError()
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := c.bookService.DeleteBook(id); err != nil {
        status := http.StatusInternalServerError
        response.Message = "Failed to delete book"
        if err.Error() == "book not found" {
            status = http.StatusNotFound
            response.Message = "Book not found"
        }
        response.Code = strconv.Itoa(status)
        response.Error = err.Error()
        response.FormatError()
        return ctx.JSON(status, response)
    }

    response.Code = strconv.Itoa(http.StatusOK)
    response.Message = "Book deleted successfully"
    response.Data = map[string]interface{}{
        "id": id,
    }
    response.Error = ""
    response.FormatError()
    return ctx.JSON(http.StatusOK, response)
}
