// cmd/main.go
package main

import (
    "fmt"
    "log"
    "auth-user-api/controllers"
    "auth-user-api/repository"
    "auth-user-api/services"
    "auth-user-api/models"
    "auth-user-api/utils"
    "auth-user-api/middleware"

    "github.com/labstack/echo/v4"
    echoMiddleware "github.com/labstack/echo/v4/middleware"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // Konfigurasi Database
    dsn := "host=localhost user=postgres password=arnoarno dbname=api-auth port=5432 sslmode=disable TimeZone=Asia/Jakarta"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Jalankan Migrasi
    err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";").Error
    if err != nil {
        log.Fatalf("Failed to create extension: %v", err)
    }

    err = db.AutoMigrate(&models.User{}, &models.Book{}, &models.Author{}, &models.Publisher{}, &models.LoanRequest{}, &models.LoanRecord{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    // Inisialisasi Repository, Service, dan Controller
    userRepo := repository.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    userController := controllers.NewUserController(userService)

    bookRepo := repository.NewBookRepository(db)
    authorRepo := repository.NewAuthorRepository(db)
    publisherRepo := repository.NewPublisherRepository(db)

    bookService := services.NewBookService(bookRepo)
    authorService := services.NewAuthorService(authorRepo)
    publisherService := services.NewPublisherService(publisherRepo)

    bookController := controllers.NewBookController(bookService, authorService, publisherService)
    authorController := controllers.NewAuthorController(authorService)
    publisherController := controllers.NewPublisherController(publisherService)

    // Inisialisasi Loan Repository, Service, dan Controller
    loanRepo := repository.NewLoanRepository(db)
    loanService := services.NewLoanService(loanRepo) // LoanService needs access to Book and User repositories
    loanController := controllers.NewLoanController(loanService)

    // Inisialisasi Echo
    e := echo.New()

    // Middleware
    e.Use(echoMiddleware.Logger())
    e.Use(echoMiddleware.Recover())

    // Validator
    e.Validator = utils.NewValidator()

    // JWT Middleware
    jwtMiddleware := middleware.NewJWTMiddleware(userService)

    // Routes
    // User Routes
    e.POST("/register", userController.RegisterUser)
    e.POST("/login", userController.LoginUser)

    // Protected User Routes
    e.POST("/register", userController.RegisterUser)
    e.POST("/login", userController.LoginUser)
    e.GET("/users", userController.GetAllUsers)
    e.PUT("/update/:id", userController.UpdateUser)
    e.DELETE("/delete", userController.DeleteUser)
    
    // Book Routes
    e.POST("/books", bookController.CreateBook)
    e.GET("/books/:id", bookController.GetBookByID)
    e.GET("/books", bookController.GetAllBooks)
    e.PUT("/books/:id", bookController.UpdateBook)
    e.DELETE("/books/:id", bookController.DeleteBook)

    // Author Routes
    e.POST("/authors", authorController.CreateAuthor)
    e.GET("/authors/:id", authorController.GetAuthorByID)
    e.GET("/authors", authorController.GetAllAuthors)
    e.PUT("/authors/:id", authorController.UpdateAuthor)
    e.DELETE("/authors/:id", authorController.DeleteAuthor)

    // Publisher Routes
    e.POST("/publishers", publisherController.CreatePublisher)
    e.GET("/publishers/:id", publisherController.GetPublisherByID)
    e.GET("/publishers", publisherController.GetAllPublishers)
    e.PUT("/publishers/:id", publisherController.UpdatePublisher)
    e.DELETE("/publishers/:id", publisherController.DeletePublisher)

    // Loan Routes
    loanGroup := e.Group("/loans", jwtMiddleware.JWTMiddleware)
    loanGroup.POST("/request", loanController.CreateLoanRequest)
    loanGroup.PUT("/cancel/:id", loanController.CancelLoanRequest)           
    loanGroup.PUT("/approve/:id", loanController.ApproveLoanRequest)       
    loanGroup.PUT("/return/:id", loanController.ReturnBook)                
    e.GET("/loan-requests", loanController.GetAllLoanRequests)
    e.GET("/loan-records", loanController.GetAllLoanRecords)
    e.GET("/loans/search/:username", loanController.SearchLoansByUsername)

    // Protected Hello Route Example
    e.GET("/protected/hello", userController.HelloProtected, jwtMiddleware.JWTMiddleware)

    // Start Server
    port := "8080"
    fmt.Printf("Server running on port %s\n", port)
    if err := e.Start(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
