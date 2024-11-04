// middleware/jwt_middleware.go
package middleware

import (
    "auth-user-api/controllers"
    "auth-user-api/domains"
    "auth-user-api/services"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v4"
    "github.com/labstack/echo/v4"
)

type JWTMiddlewareConfig struct {
    UserService services.UserService // Inject UserService
}

func NewJWTMiddleware(userService services.UserService) *JWTMiddlewareConfig {
    return &JWTMiddlewareConfig{UserService: userService}
}

// JWTMiddleware verifies the token and sets role and username in context.
func (mw *JWTMiddlewareConfig) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(ctx echo.Context) error {
        tokenString := ctx.Request().Header.Get("Authorization")

        if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
            return ctx.JSON(http.StatusUnauthorized, domains.BaseResponse{
                Code:    "401",
                Message: "Missing or invalid Authorization header",
                Error:   "Authorization header error",
            })
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        claims := &controllers.JWTClaims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte("my_secret_key"), nil
        })

        if err != nil || !token.Valid {
            return ctx.JSON(http.StatusUnauthorized, domains.BaseResponse{
                Code:    "401",
                Message: "Invalid token",
                Error:   "Token validation error",
            })
        }

        user, err := mw.UserService.GetUserByUsername(claims.Username)
        if err != nil || user == nil {
            return ctx.JSON(http.StatusUnauthorized, domains.BaseResponse{
                Code:    "401",
                Message: "Invalid token - user not found",
                Error:   "User not found or deleted",
            })
        }

        // Set username and role in context
        ctx.Set("username", claims.Username)
        ctx.Set("role", claims.Role)

        return next(ctx)
    }
}

// LoanRequestSummary provides a summary structure for each loan request in the list response
type LoanRequestSummary struct {
    ID           uint   `json:"id"`
    BookID       int    `json:"book_id"`
    BorrowerName string `json:"borrower_name"`
    Status       string `json:"status"`
    RequestDate  string `json:"request_date"`
}

// LoanRecordListResponse represents the structure for retrieving all loan records
type LoanRecordListResponse struct {
    TotalRecords int                `json:"total_records"`
    LoanRecords  []LoanRecordDetail `json:"loan_records"`
}

// LoanRecordDetail provides detailed information for each loan record
type LoanRecordDetail struct {
    ID           uint   `json:"id"`
    BookID       int    `json:"book_id"`
    BorrowerName string `json:"borrower_name"`
    Status       string `json:"status"`
    BorrowDate   string `json:"borrow_date"`
    ReturnDate   *string `json:"return_date,omitempty"` // Nullable if not returned
}