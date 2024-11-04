// services/loan_services.go

package services

import (
    "auth-user-api/models"
    "auth-user-api/repository"
    "errors"
    "time"
)

// Tambahkan variabel error untuk kode 404
var ErrBookOutOfStock = errors.New("book out of stock")

type LoanService struct {
    Repo *repository.LoanRepository
}

func NewLoanService(repo *repository.LoanRepository) *LoanService {
    return &LoanService{Repo: repo}
}

// Cek stok buku sebelum membuat request peminjaman
func (s *LoanService) CreateLoanRequest(req *models.LoanRequest) error {
    // Periksa apakah stok buku ada
    book, err := s.Repo.GetBookByID(int(req.BookID))
    if err != nil {
        return err
    }
    if book.Stock <= 0 {
        return ErrBookOutOfStock
    }
    
    req.RequestTime = time.Now()
    req.Status = "PENDING"
    return s.Repo.CreateLoanRequest(req)
}

func (s *LoanService) ApproveLoanRequest(requestID uint) (*models.LoanRecord, error) {
    req, err := s.Repo.GetLoanRequestByID(requestID)
    if err != nil {
        return nil, err
    }

    if req.Status != "PENDING" {
        return nil, errors.New("request already processed")
    }

    book, err := s.Repo.GetBookByID(req.BookID)
    if err != nil || book.Stock <= 0 {
        return nil, errors.New("book out of stock")
    }

    req.Status = "APPROVED"
    if err := s.Repo.UpdateLoanRequest(req); err != nil {
        return nil, err
    }

    loan := &models.LoanRecord{
        BookID:   req.BookID,
        UserID:   req.UserID,
        LoanDate: time.Now(),
        DueDate:  time.Now().AddDate(0, 0, 3), // Menetapkan tanggal pengembalian otomatis 3 hari dari sekarang
    }

    if err := s.Repo.CreateLoanRecord(loan); err != nil {
        return nil, err
    }

    s.Repo.UpdateBookStock(req.BookID, -1) // Kurangi stok
    return loan, nil
}

// RejectLoanRequest rejects a loan request with a custom reason
func (s *LoanService) RejectLoanRequest(requestID uint, reason string) error {
    req, err := s.Repo.GetLoanRequestByID(requestID)
    if err != nil {
        return err
    }

    if req.Status != "PENDING" {
        return errors.New("request already processed")
    }

    req.Status = "REJECTED"
    req.RejectReason = &reason // Set the custom rejection reason

    return s.Repo.UpdateLoanRequest(req)
}

func (s *LoanService) ReturnBook(loanID uint) (*models.LoanRecord, int, error) {
    loan, err := s.Repo.GetLoanRecordByID(loanID)
    if err != nil {
        return nil, 0, err
    }

    if loan.Returned {
        return nil, 0, errors.New("book already returned")
    }

    loan.Returned = true
    loan.ReturnDate = timePtr(time.Now())

    lateFee := 0
    if time.Now().After(loan.DueDate) {
        daysLate := int(time.Since(loan.DueDate).Hours() / 24)
        lateFee = daysLate * 5000
    }

    s.Repo.UpdateLoanRecord(loan)
    s.Repo.UpdateBookStock(loan.BookID, 1) // Increase stock

    return loan, lateFee, nil
}

func timePtr(t time.Time) *time.Time {
    return &t
}

// GetAllLoanRequests fetches all loan requests with borrower names.
func (s *LoanService) GetAllLoanRequests() ([]map[string]interface{}, error) {
    return s.Repo.GetAllLoanRequests()
}

// GetAllLoanRecords retrieves all loan records.
func (s *LoanService) GetAllLoanRecords() ([]map[string]interface{}, error) {
    return s.Repo.GetAllLoanRecords()
}

// SearchLoansByUsername retrieves all loans associated with a given username.
func (s *LoanService) SearchLoansByUsername(username string) ([]map[string]interface{}, error) {
    return s.Repo.GetLoansByUsername(username)
}

// CancelLoanRequest cancels a loan request with a custom reason
func (s *LoanService) CancelLoanRequest(requestID uint, reason string) error {
    req, err := s.Repo.GetLoanRequestByID(requestID)
    if err != nil {
        return err
    }

    if req.Status != "PENDING" {
        return errors.New("request already processed")
    }

    req.Status = "CANCELLED"
    req.RejectReason = &reason // Set the custom cancellation reason

    return s.Repo.UpdateLoanRequest(req)
}
