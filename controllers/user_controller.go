// controllers/user_controller.go

package controllers

import (
    "net/http"
    "time"
    "github.com/golang-jwt/jwt/v4"
    "auth-user-api/services"
    "auth-user-api/domains"
    "github.com/labstack/echo/v4"
)

type UserController struct {
    service services.UserService
}

func NewUserController(service services.UserService) *UserController {
    return &UserController{service}
}

// Register User godoc
func (c *UserController) RegisterUser(ctx echo.Context) error {
    type RegisterRequest struct {
        Username  string `json:"username" validate:"required"`
        Email     string `json:"email" validate:"required,email"`
        Password1 string `json:"password_1" validate:"required"`
        Password2 string `json:"password_2" validate:"required"`
        Role      int    `json:"role" validate:"required,oneof=1 2"` // 1 for admin, 2 for member
    }

    var req RegisterRequest
    if err := ctx.Bind(&req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Failed processing input, try again. Error: " + err.Error(),
            Error:   "Binding error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    // Perform additional validation if necessary
    if req.Role != 1 && req.Role != 2 {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Role must be 1 (admin) or 2 (member)",
            Error:   "Validation error",
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := ctx.Validate(req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Validation error. Field: " + err.Error(),
            Error:   "Validation error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    // Register the user with the role
    if err := c.service.Register(req.Username, req.Email, req.Password1, req.Password2, req.Role); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Registration failed. Error: " + err.Error(),
            Error:   "Service error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    userResponse := domains.RegisterResponse{
        Username: req.Username,
        Email:    req.Email,
        Role:     "", // Temporary empty, will assign below
    }
    
    if req.Role == 1 {
        userResponse.Role = "admin"
    } else if req.Role == 2 {
        userResponse.Role = "member"
    }
    
    response := domains.BaseResponse{
        Code:      "200",
        Message:   "User successfully registered",
        Data:      userResponse,
        Parameter: "username",
    }
    return ctx.JSON(http.StatusOK, response)
}

// GetAllUsers retrieves all users with roles in string format (admin/member)
func (c *UserController) GetAllUsers(ctx echo.Context) error {
    users, err := c.service.GetAllUsers()
    if err != nil {
        response := domains.BaseResponse{
            Code:    "500",
            Message: "Failed to retrieve users. Error: " + err.Error(),
            Error:   "Service error: " + err.Error(),
        }
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    // Map to UserResponse with role as a string
    var userResponses []domains.UserResponse
    for _, user := range users {
        role := "member" // default to member
        if user.Role == 1 {
            role = "admin"
        }
        
        userResponse := domains.UserResponse{
            UserID:   user.ID,
            Username: user.Username,
            Email:    user.Email,
            Role:     role,
        }
        userResponses = append(userResponses, userResponse)
    }

    response := domains.BaseResponse{
        Code:    "200",
        Message: "Users retrieved successfully",
        Data:    userResponses,
    }
    return ctx.JSON(http.StatusOK, response)
}

// Update User godoc
func (c *UserController) UpdateUser(ctx echo.Context) error {
    type UpdateRequest struct {
        Username  string `json:"username" validate:"required"`
        Email     string `json:"email"`
        Password1 string `json:"password_1"`
        Password2 string `json:"password_2"`
    }

    userID := ctx.Param("id")
    if userID == "" {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "User ID is required. Field: id",
            Error:   "Validation error",
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    existingUser, err := c.service.GetUserByID(userID)
    if err != nil {
        response := domains.BaseResponse{
            Code:    "404",
            Message: "User not found. UserID: " + userID,
            Error:   "User retrieval error: " + err.Error(),
        }
        return ctx.JSON(http.StatusNotFound, response)
    }

    var req UpdateRequest
    if err := ctx.Bind(&req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Failed processing input. Error: " + err.Error(),
            Error:   "Binding error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    err = c.service.Update(userID, req.Username, req.Email, req.Password1, req.Password2)
    if err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Failed to update user. Error: " + err.Error(),
            Error:   "Service error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    userResponse := domains.UserResponse{
        UserID:   existingUser.ID,
        Username: req.Username,
        Email:    req.Email,
        Password: existingUser.Password,
    }

    response := domains.BaseResponse{
        Code:      "200",
        Message:   "User successfully updated. UserID: " + userID,
        Data:      userResponse,
        Parameter: "username", 
    }    
    return ctx.JSON(http.StatusOK, response)
}

// Delete User godoc
func (c *UserController) DeleteUser(ctx echo.Context) error {
    type DeleteRequest struct {
        UserID string `json:"user_id" validate:"required"`
    }

    var req DeleteRequest
    if err := ctx.Bind(&req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Delete failed. Error: " + err.Error(),
            Error:   "Binding error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := ctx.Validate(req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Delete failed. Validation error: " + err.Error(),
            Error:   "Validation error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    user, err := c.service.GetUserByID(req.UserID)
    if err != nil {
        response := domains.BaseResponse{
            Code:    "404",
            Message: "User not found. UserID: " + req.UserID,
            Error:   "User retrieval error: " + err.Error(),
        }
        return ctx.JSON(http.StatusNotFound, response)
    }

    if user.DeletedAt.Valid {
        response := domains.BaseResponse{
            Code:    "404",
            Message: "User not found. UserID: " + req.UserID,
        }
        return ctx.JSON(http.StatusNotFound, response)
    }

    if err := c.service.Delete(req.UserID); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Failed to delete user. UserID: " + req.UserID + ", Error: " + err.Error(),
            Error:   "Service error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    response := domains.BaseResponse{
        Code:      "200",
        Message:   "User deleted successfully. UserID: " + req.UserID,
        Data:      domains.DeleteResponse{UserID: req.UserID},
        Parameter: "user_id", 
    }    
    return ctx.JSON(http.StatusOK, response)
}

var jwtKey = []byte("my_secret_key")  // Pastikan menggunakan secret key yang sama

type JWTClaims struct {
    Username string `json:"username"`
    Role     int    `json:"role"`
    jwt.RegisteredClaims
}

// Login User
func (c *UserController) LoginUser(ctx echo.Context) error {
    type LoginRequest struct {
        Username string `json:"username" validate:"required"`
        Password string `json:"password" validate:"required"`
    }

    var req LoginRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.BaseResponse{
            Code:    "400",
            Message: "Invalid input",
            Error:   err.Error(),
        })
    }

    if err := ctx.Validate(req); err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.BaseResponse{
            Code:    "400",
            Message: "Validation error",
            Error:   err.Error(),
        })
    }

    user, err := c.service.Authenticate(req.Username, req.Password)
    if err != nil {
        return ctx.JSON(http.StatusUnauthorized, domains.BaseResponse{
            Code:    "401",
            Message: "Invalid username or password",
            Error:   "AuthenticationError",
        })
    }
    
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &JWTClaims{
        Username: user.Username,
        Role:     user.Role, // Ambil role dari user yang berhasil diotentikasi
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.BaseResponse{
            Code:    "500",
            Message: "Failed to generate token",
            Error:   err.Error(),
        })
    }
    
    return ctx.JSON(http.StatusOK, domains.BaseResponse{
        Code:    "200",
        Message: "Successful login",
        Data: map[string]interface{}{
            "token": tokenString,
        },
    })
}    

// Route yang diproteksi
func (c *UserController) HelloProtected(ctx echo.Context) error {
    username := ctx.Get("username").(string)
    role := ctx.Get("role").(int) // Ambil role dari context

    // Izinkan akses untuk role 1 (admin) dan role 2 (member)
    if role == 1 {
        return ctx.JSON(http.StatusOK, map[string]string{
            "Message": "Hello, admin! You have accessed a protected route!",
            "User":    username,
        })
    } else if role == 2 {
        return ctx.JSON(http.StatusOK, map[string]string{
            "Message": "Hello, member! You have accessed a protected route!",
            "User":    username,
        })
    }

    // Jika role tidak 1 atau 2, tolak akses
    return ctx.JSON(http.StatusForbidden, map[string]string{
        "Message": "Access denied",
    })
}
