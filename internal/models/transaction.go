package models

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          int64           `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateTransactionRequest struct {
	Type        TransactionType `json:"type" validate:"required,oneof=income expense"`
	Amount      float64         `json:"amount" validate:"required,gt=0"`
	Category    string          `json:"category" validate:"required,min=1,max=100"`
	Description string          `json:"description" validate:"max=500"`
	Date        string          `json:"date" validate:"required"`
}

type UpdateTransactionRequest struct {
	Type        TransactionType `json:"type" validate:"omitempty,oneof=income expense"`
	Amount      float64         `json:"amount" validate:"omitempty,gt=0"`
	Category    string          `json:"category" validate:"omitempty,min=1,max=100"`
	Description string          `json:"description" validate:"omitempty,max=500"`
	Date        string          `json:"date"`
}

type TransactionFilter struct {
	Type     TransactionType `query:"type"`
	Category string          `query:"category"`
	FromDate string          `query:"from_date"`
	ToDate   string          `query:"to_date"`
	Page     int             `query:"page"`
	Limit    int             `query:"limit"`
}

type TransactionSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}
