// controllers/loan_controller.go
package controllers

import (
    "auth-user-api/models"
    "auth-user-api/services"
    "auth-user-api/domains"
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/google/uuid"
    "strconv"
    "errors"
    "time"
)

type LoanController struct {
    Service *services.LoanService
}

func NewLoanController(service *services.LoanService) *LoanController {
    return &LoanController{Service: service}
}

// Create Loan Request
func (lc *LoanController) CreateLoanRequest(ctx echo.Context) error {
    var req models.LoanRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.NewErrorResponse("400", "Failed to parse request", err.Error()))
    }

    userID, err := uuid.Parse(req.UserID.String())
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.NewErrorResponse("400", "Invalid user UUID", err.Error()))
    }
    req.UserID = userID

    if err := lc.Service.CreateLoanRequest(&req); err != nil {
        // Tangani error stok habis sebagai 404
        if errors.Is(err, services.ErrBookOutOfStock) {
            return ctx.JSON(http.StatusNotFound, domains.NewErrorResponse("404", "Book out of stock", err.Error()))
        }
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to create loan request", err.Error()))
    }

    username, err := lc.Service.Repo.GetUsernameByUserID(req.UserID)
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to retrieve borrower username", err.Error()))
    }

    loanResponse := domains.LoanRequestResponse{
        ID:           req.ID,
        BookID:       req.BookID,
        UserID:       req.UserID.String(),
        BorrowerName: username,
        Status:       "PENDING",
        RequestDate:  time.Now().Format(time.RFC3339),
    }

    response := domains.NewSuccessResponseWithData("200", "Loan request created successfully", loanResponse)
    return ctx.JSON(http.StatusOK, response)
}

// ApproveLoanRequest handles loan request approval.
func (lc *LoanController) ApproveLoanRequest(ctx echo.Context) error {
    // Check if the user is an admin
    role := ctx.Get("role").(int)
    if role != 1 {
        return ctx.JSON(http.StatusForbidden, domains.NewErrorResponse("403", "Access denied", "Only admins can approve loan requests"))
    }

    requestID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.NewErrorResponse("400", "Invalid request ID", err.Error()))
    }

    var body struct {
        Approve bool   `json:"approve"`
        Reason  string `json:"reason"`
    }
    if err := ctx.Bind(&body); err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.NewErrorResponse("400", "Failed to parse request body", err.Error()))
    }

    if body.Approve {
        loan, err := lc.Service.ApproveLoanRequest(uint(requestID))
        if err != nil {
            return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to approve loan request", err.Error()))
        }

        loanInfo := domains.LoanApprovalResponse{
            ID:        loan.ID,
            BookID:    loan.BookID,
            UserID:    loan.UserID.String(),
            LoanDate:  loan.LoanDate.Format(time.RFC3339),
            DueDate:   loan.DueDate.Format(time.RFC3339),
            Returned:  false,
        }

        response := domains.NewSuccessResponseWithData("200", "Loan request approved", loanInfo)
        return ctx.JSON(http.StatusOK, response)
    } else {
        reason := body.Reason
        if reason == "" {
            reason = "No specific reason provided"
        }

        if err := lc.Service.RejectLoanRequest(uint(requestID), reason); err != nil {
            return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to reject loan request", err.Error()))
        }

        rejectionData := domains.LoanRejectionResponse{
            ID:     uint(requestID),
            Status: "REJECTED",
            Reason: reason,
        }
        response := domains.NewSuccessResponseWithData("200", "Loan request rejected", rejectionData)
        return ctx.JSON(http.StatusOK, response)
    }
}

// ReturnBook with Late Fee Handling
func (lc *LoanController) ReturnBook(ctx echo.Context) error {
    // Check if the user is an admin
    role := ctx.Get("role").(int)
    if role != 1 {
        return ctx.JSON(http.StatusForbidden, domains.NewErrorResponse("403", "Access denied", "Only admins can return books"))
    }

    loanID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.NewErrorResponse("400", "Invalid loan ID", err.Error()))
    }

    loan, lateFee, err := lc.Service.ReturnBook(uint(loanID))
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to return book", err.Error()))
    }

    username, err := lc.Service.Repo.GetUsernameByUserID(loan.UserID)
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to retrieve borrower username", err.Error()))
    }

    loanRecord := domains.LoanReturnResponse{
        ID:           loan.ID,
        Status:       "RETURNED",
        BorrowerName: username,
        ReturnDate:   loan.ReturnDate.Format(time.RFC3339),
        LateFee:      lateFee,
    }

    response := domains.NewSuccessResponseWithData("200", "Book returned successfully", loanRecord)
    return ctx.JSON(http.StatusOK, response)
}

// GetAllLoanRequests retrieves all loan requests with borrower names.
func (lc *LoanController) GetAllLoanRequests(ctx echo.Context) error {
    loanRequests, err := lc.Service.GetAllLoanRequests()
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to fetch loan requests", err.Error()))
    }

    response := domains.NewSuccessResponseWithData("200", "Loan requests retrieved successfully", loanRequests)
    return ctx.JSON(http.StatusOK, response)
}

// GetAllLoanRecords retrieves all loan records with borrower and loan details.
func (lc *LoanController) GetAllLoanRecords(ctx echo.Context) error {
    loanRecords, err := lc.Service.GetAllLoanRecords()
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to retrieve loan records", err.Error()))
    }

    response := domains.NewSuccessResponseWithData("200", "Loan records retrieved successfully", loanRecords)
    return ctx.JSON(http.StatusOK, response)
}

// SearchLoansByUsername fetches all loans associated with the given username.
func (lc *LoanController) SearchLoansByUsername(ctx echo.Context) error {
    username := ctx.Param("username")

    loans, err := lc.Service.SearchLoansByUsername(username)
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to fetch loans for the specified username", err.Error()))
    }

    responseData := domains.LoanSearchResponse{
        Username: username,
        Loans:    loans,
    }

    response := domains.NewSuccessResponseWithData("200", "Loans retrieved successfully", responseData)
    return ctx.JSON(http.StatusOK, response)
}

// CancelLoanRequest handles loan request cancellation with a reason
func (lc *LoanController) CancelLoanRequest(ctx echo.Context) error {
    // Validasi token melalui middleware JWT yang sudah ada
    username := ctx.Get("username").(string)

    requestID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.NewErrorResponse("400", "Invalid request ID", err.Error()))
    }

    var body struct {
        Reason string `json:"reason"`
    }
    if err := ctx.Bind(&body); err != nil {
        return ctx.JSON(http.StatusBadRequest, domains.NewErrorResponse("400", "Failed to parse request body", err.Error()))
    }

    reason := body.Reason
    if reason == "" {
        reason = "No specific reason provided"
    }

    // Ambil data loan request berdasarkan ID dan verifikasi apakah username sesuai
    loanRequest, err := lc.Service.Repo.GetLoanRequestByID(uint(requestID))
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to retrieve loan request", err.Error()))
    }

    borrowerUsername, err := lc.Service.Repo.GetUsernameByUserID(loanRequest.UserID)
    if err != nil || borrowerUsername != username {
        return ctx.JSON(http.StatusUnauthorized, domains.NewErrorResponse("401", "Unauthorized to cancel this loan request", "User does not match loan borrower"))
    }

    if err := lc.Service.CancelLoanRequest(uint(requestID), reason); err != nil {
        return ctx.JSON(http.StatusInternalServerError, domains.NewErrorResponse("500", "Failed to cancel loan request", err.Error()))
    }

    cancellationData := domains.LoanCancellationResponse{
        ID:     uint(requestID),
        Status: "CANCELLED",
        Reason: reason,
    }
    response := domains.NewSuccessResponseWithData("200", "Loan request cancelled", cancellationData)
    return ctx.JSON(http.StatusOK, response)
}
