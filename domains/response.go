// domains/response.go
package domains

// BaseResponse is the general structure for all API responses
type BaseResponse struct {
    Code      string      `json:"code"`                 // HTTP response code
    Message   string      `json:"message"`              // Response message
    Data      interface{} `json:"data,omitempty"`       // Data payload (optional)
    Error     string      `json:"error,omitempty"`      // Error details (optional)
    Parameter string      `json:"parameter,omitempty"`  // Related parameter (optional)
}

// FormatError sets Error to "null" if it's an empty string
func (r *BaseResponse) FormatError() {
    if r.Error == "" {
        r.Error = "null"
    }
}

// TokenResponse represents a response with a token
type TokenResponse struct {
    Token string `json:"token"`  // JWT token string
}

// UserResponse represents the user details in the response
type UserResponse struct {
    UserID   string `json:"user_id"`      // Unique user ID
    Username string `json:"username"`     // User's username
    Email    string `json:"email"`        // User's email
    Password string `json:"-"`            // Password is omitted in the response
    Role     string `json:"role"`         // User role (admin or member)
}

// DeleteResponse represents the response after a user is deleted
type DeleteResponse struct {
    UserID string `json:"user_id"`  // ID of the deleted user
}

// RegisterResponse represents the response after user registration
type RegisterResponse struct {
    Username  string `json:"username"`   // Registered username
    Email     string `json:"email"`      // Registered email
    Password1 string `json:"-"`          // Password is hidden in the response
    Password2 string `json:"-"`          // Password is hidden in the response
    Role      string `json:"role"`       // User role (admin or member)
}

// ErrorResponse is used to format error messages with extra details
type ErrorResponse struct {
    Code      string            `json:"code"`                 // HTTP response code
    Message   string            `json:"message"`              // Error message
    Errors    map[string]string `json:"errors,omitempty"`     // Map of field errors (optional)
    Parameter string            `json:"parameter,omitempty"`  // Related parameter (optional)
}

type BookResponse struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	Summary     string            `json:"summary"`
	AuthorID    int               `json:"author_id"`
	Author      BookAuthorResponse    `json:"Author"`
	PublisherID int               `json:"publisher_id"`
	Publisher   BookPublisherResponse `json:"Publisher"`
	Stock       int               `json:"stock"`
	MaxStock    int               `json:"max_stock"`
	CreatedAt   string            `json:"CreatedAt"`
	UpdatedAt   string            `json:"UpdatedAt"`
	DeletedAt   *string           `json:"DeletedAt,omitempty"`
}

type BookAuthorResponse struct {
    Name string `json:"name"`
}

type BookPublisherResponse struct {
    Name string `json:"name"`
}

// AuthorResponse is the struct for author data response.
type AuthorResponse struct {
    ID        int     `json:"id"`
    Name      string  `json:"name"`
    CreatedAt string  `json:"CreatedAt"`
    UpdatedAt string  `json:"UpdatedAt"`
    DeletedAt *string `json:"DeletedAt,omitempty"`
}

// PublisherResponse struct for publisher data response
type PublisherResponse struct {
    ID        int        `json:"id"`
    Name      string     `json:"name"`
    CreatedAt string  `json:"createdAt"`
    UpdatedAt string  `json:"updatedAt"`
    DeletedAt *string `json:"deletedAt,omitempty"`
}


// Helper functions to create response

func NewErrorResponse(code, message, err string) BaseResponse {
    response := BaseResponse{
        Code:    code,
        Message: message,
        Error:   err,
    }
    response.FormatError()
    return response
}

func NewSuccessResponseWithData(code, message string, data interface{}) BaseResponse {
    return BaseResponse{
        Code:    code,
        Message: message,
        Data:    data,
    }
}

func NewSuccessResponse(code, message string) BaseResponse {
    return BaseResponse{
        Code:    code,
        Message: message,
    }
}

// LoanRequestResponse represents the structure for loan request responses
type LoanRequestResponse struct {
    ID           uint   `json:"id"`
    BookID       int    `json:"book_id"`
    UserID       string `json:"user_id"`
    BorrowerName string `json:"borrower_name"`
    Status       string `json:"status"`
    RequestDate  string `json:"request_date"`
}

type LoanApprovalResponse struct {
    ID        uint   `json:"id"`
    BookID    int    `json:"book_id"`
    UserID    string `json:"user_id"`
    LoanDate  string `json:"loan_date"`
    DueDate   string `json:"due_date"`
    Returned  bool   `json:"returned"`
}

// LoanRejectionResponse represents the response when a loan request is rejected
type LoanRejectionResponse struct {
    ID     uint   `json:"id"`
    Status string `json:"status"`
    Reason string `json:"reason"`
}

// LoanReturnResponse represents the response when a loan is returned
type LoanReturnResponse struct {
    ID           uint   `json:"id"`
    Status       string `json:"status"`
    BorrowerName string `json:"borrower_name"`
    ReturnDate   string `json:"return_date"`
    LateFee      int    `json:"late_fee"`
}

// LoanSearchResponse represents the structure for searching loans by username
type LoanSearchResponse struct {
    Username string      `json:"username"`
    Loans    interface{} `json:"loans"`
}

type LoanRecordResponse struct {
    ID           uint   `json:"id"`
    BookID       int    `json:"book_id"`
    UserID       string `json:"user_id"`
    BorrowerName string `json:"borrower_name"`
    LoanDate     string `json:"loan_date"`
    DueDate      string `json:"due_date"`
    Returned     bool   `json:"returned"`
    ReturnDate   *string `json:"return_date,omitempty"`
}

// LoanRequestDetails is the structure for detailed loan request response
type LoanRequestDetails struct {
    ID           uint   `json:"id"`
    BookID       int    `json:"book_id"`
    UserID       string `json:"user_id"`
    BorrowerName string `json:"borrower_name"`
    Status       string `json:"status"`
    RequestDate  string `json:"request_date"`
}

// LoanRecordDetails represents a detailed response for loan records
type LoanRecordDetails struct {
    ID           uint    `json:"id"`
    BookID       int     `json:"book_id"`
    UserID       string  `json:"user_id"`
    BorrowerName string  `json:"borrower_name"`
    LoanDate     string  `json:"loan_date"`
    DueDate      string  `json:"due_date"`
    Returned     bool    `json:"returned"`
    ReturnDate   *string `json:"return_date,omitempty"`
}

// LoanSearchDetails is the structure for detailed search results by username
type LoanSearchDetails struct {
    ID           uint    `json:"id"`
    BookID       int     `json:"book_id"`
    BorrowerName string  `json:"borrower_name"`
    LoanDate     string  `json:"loan_date"`
    DueDate      string  `json:"due_date"`
    Returned     bool    `json:"returned"`
}

type LoanCancellationResponse struct {
    ID     uint   `json:"id"`
    Status string `json:"status"`
    Reason string `json:"reason"`
}