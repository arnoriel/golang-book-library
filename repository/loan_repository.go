// repository/loan_repository.go
package repository

import (
    "auth-user-api/models"
    "gorm.io/gorm"
    "github.com/google/uuid"
)

type LoanRepository struct {
    DB *gorm.DB
}

func NewLoanRepository(db *gorm.DB) *LoanRepository {
    return &LoanRepository{DB: db}
}

func (r *LoanRepository) CreateLoanRequest(req *models.LoanRequest) error {
    return r.DB.Create(req).Error
}

func (r *LoanRepository) UpdateLoanRequest(req *models.LoanRequest) error {
    return r.DB.Save(req).Error
}

func (r *LoanRepository) GetLoanRequestByID(id uint) (*models.LoanRequest, error) {
    var req models.LoanRequest
    if err := r.DB.First(&req, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &req, nil
}

func (r *LoanRepository) CreateLoanRecord(record *models.LoanRecord) error {
    return r.DB.Create(record).Error
}

func (r *LoanRepository) UpdateLoanRecord(record *models.LoanRecord) error {
    return r.DB.Save(record).Error
}

func (r *LoanRepository) GetLoanRecordByID(id uint) (*models.LoanRecord, error) {
    var record models.LoanRecord
    if err := r.DB.First(&record, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &record, nil
}

func (r *LoanRepository) GetBookByID(id int) (*models.Book, error) {
    var book models.Book
    return &book, r.DB.First(&book, "id = ?", id).Error
}

func (r *LoanRepository) UpdateBookStock(bookID int, change int) error {
    return r.DB.Model(&models.Book{}).
        Where("id = ?", bookID).
        Update("stock", gorm.Expr("stock + ?", change)).Error
}

// GetLoanRecordWithUserInfo fetches loan record with the corresponding user's username.
func (r *LoanRepository) GetLoanRecordWithUserInfo(loanID uint) (map[string]interface{}, error) {
    var result map[string]interface{}

    query := `
        SELECT lr.id, lr.book_id, lr.loan_date, lr.due_date, lr.returned, 
               u.username AS borrower_name
        FROM loan_records lr
        JOIN users u ON lr.user_id = u.id
        WHERE lr.id = ?;
    `

    if err := r.DB.Raw(query, loanID).Scan(&result).Error; err != nil {
        return nil, err
    }
    return result, nil
}

// GetUsernameByUserID fetches the username associated with a given user UUID.
func (r *LoanRepository) GetUsernameByUserID(userID uuid.UUID) (string, error) {
    var username string
    query := `SELECT username FROM users WHERE id = ?`
    if err := r.DB.Raw(query, userID).Scan(&username).Error; err != nil {
        return "", err
    }
    return username, nil
}

// GetAllLoanRequests fetches all loan requests along with the borrower's username.
func (r *LoanRepository) GetAllLoanRequests() ([]map[string]interface{}, error) {
    var results []map[string]interface{}

    query := `
        SELECT lr.id, lr.book_id, lr.request_time, lr.status, lr.reject_reason,
               u.username AS borrower_name
        FROM loan_requests lr
        JOIN users u ON lr.user_id = u.id;
    `

    if err := r.DB.Raw(query).Scan(&results).Error; err != nil {
        return nil, err
    }
    return results, nil
}

// GetAllLoanRecords retrieves all loan records with the corresponding user's username and return_date.
func (r *LoanRepository) GetAllLoanRecords() ([]map[string]interface{}, error) {
    var results []map[string]interface{}

    query := `
        SELECT lr.id, lr.book_id, lr.loan_date, lr.due_date, lr.returned, 
               lr.return_date, u.username AS borrower_name
        FROM loan_records lr
        JOIN users u ON lr.user_id = u.id;
    `

    if err := r.DB.Raw(query).Scan(&results).Error; err != nil {
        return nil, err
    }
    return results, nil
}

// GetLoansByUsername retrieves all loans associated with the given username.
func (r *LoanRepository) GetLoansByUsername(username string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}

    query := `
        SELECT lr.id, lr.book_id, lr.loan_date, lr.due_date, lr.returned,
               lr.return_date, u.username AS borrower_name
        FROM loan_records lr
        JOIN users u ON lr.user_id = u.id
        WHERE u.username = ?;
    `

    if err := r.DB.Raw(query, username).Scan(&results).Error; err != nil {
        return nil, err
    }

    return results, nil
}

// GetActiveLoanByUsername retrieves an active loan record for the given username.
func (r *LoanRepository) GetActiveLoanByUsername(username string) (*models.LoanRecord, error) {
    var loan models.LoanRecord

    query := `
        SELECT lr.*
        FROM loan_records lr
        JOIN users u ON lr.user_id = u.id
        WHERE u.username = ? AND lr.returned = false
        LIMIT 1;
    `

    if err := r.DB.Raw(query, username).Scan(&loan).Error; err != nil {
        return nil, err
    }
    return &loan, nil
}

// Fungsi untuk mencari LoanRequest berdasarkan ID
func (r *LoanRepository) FindByID(requestID uint) (*models.LoanRequest, error) {
    var loanRequest models.LoanRequest
    if err := r.DB.First(&loanRequest, requestID).Error; err != nil {
        return nil, err
    }
    return &loanRequest, nil
}