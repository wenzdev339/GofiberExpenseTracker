package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"gofiber-expense-tracker/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (type, amount, category, description, date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query,
		tx.Type, tx.Amount, tx.Category, tx.Description, tx.Date,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)
}

func (r *TransactionRepository) GetByID(id int64) (*models.Transaction, error) {
	tx := &models.Transaction{}
	query := `
		SELECT id, type, amount, category, description, date, created_at, updated_at
		FROM transactions WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID, &tx.Type, &tx.Amount, &tx.Category,
		&tx.Description, &tx.Date, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *TransactionRepository) List(filter models.TransactionFilter) ([]models.Transaction, int64, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.FromDate != "" {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", argIdx))
		args = append(args, filter.FromDate)
		argIdx++
	}
	if filter.ToDate != "" {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", argIdx))
		args = append(args, filter.ToDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	var totalCount int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions %s", whereClause)
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}
	offset := (filter.Page - 1) * filter.Limit

	query := fmt.Sprintf(`
		SELECT id, type, amount, category, description, date, created_at, updated_at
		FROM transactions %s
		ORDER BY date DESC, id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(
			&tx.ID, &tx.Type, &tx.Amount, &tx.Category,
			&tx.Description, &tx.Date, &tx.CreatedAt, &tx.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, totalCount, nil
}

func (r *TransactionRepository) Update(tx *models.Transaction) error {
	query := `
		UPDATE transactions
		SET type = $1, amount = $2, category = $3, description = $4, date = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`
	return r.db.QueryRow(query,
		tx.Type, tx.Amount, tx.Category, tx.Description, tx.Date, tx.ID,
	).Scan(&tx.UpdatedAt)
}

func (r *TransactionRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM transactions WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TransactionRepository) GetSummary(fromDate, toDate string) (*models.TransactionSummary, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if fromDate != "" {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", argIdx))
		args = append(args, fromDate)
		argIdx++
	}
	if toDate != "" {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", argIdx))
		args = append(args, toDate)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense
		FROM transactions %s
	`, whereClause)

	summary := &models.TransactionSummary{}
	err := r.db.QueryRow(query, args...).Scan(&summary.TotalIncome, &summary.TotalExpense)
	if err != nil {
		return nil, err
	}
	summary.Balance = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}
