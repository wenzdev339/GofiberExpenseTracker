package services

import (
	"errors"
	"math"
	"time"

	"gofiber-expense-tracker/internal/models"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidInput        = errors.New("invalid input")
)

type TransactionRepo interface {
	Create(tx *models.Transaction) error
	GetByID(id int64) (*models.Transaction, error)
	List(filter models.TransactionFilter) ([]models.Transaction, int64, error)
	Update(tx *models.Transaction) error
	Delete(id int64) error
	GetSummary(fromDate, toDate string) (*models.TransactionSummary, error)
}

type TransactionService struct {
	repo TransactionRepo
}

func NewTransactionService(repo TransactionRepo) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Create(req models.CreateTransactionRequest) (*models.Transaction, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidInput
	}
	if req.Type != models.Income && req.Type != models.Expense {
		return nil, ErrInvalidInput
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	tx := &models.Transaction{
		Type:        req.Type,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Date:        date,
	}

	if err := s.repo.Create(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *TransactionService) GetByID(id int64) (*models.Transaction, error) {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrTransactionNotFound
	}
	return tx, nil
}

func (s *TransactionService) List(filter models.TransactionFilter) (*models.PaginatedResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	transactions, totalCount, err := s.repo.List(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.Limit)))

	return &models.PaginatedResponse{
		Data:       transactions,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (s *TransactionService) Update(id int64, req models.UpdateTransactionRequest) (*models.Transaction, error) {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrTransactionNotFound
	}

	if req.Type != "" {
		tx.Type = req.Type
	}
	if req.Amount > 0 {
		tx.Amount = req.Amount
	}
	if req.Category != "" {
		tx.Category = req.Category
	}
	if req.Description != "" {
		tx.Description = req.Description
	}
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, expected YYYY-MM-DD")
		}
		tx.Date = date
	}

	if err := s.repo.Update(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *TransactionService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *TransactionService) GetSummary(fromDate, toDate string) (*models.TransactionSummary, error) {
	return s.repo.GetSummary(fromDate, toDate)
}
