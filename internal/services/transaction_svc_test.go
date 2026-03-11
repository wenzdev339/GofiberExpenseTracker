package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gofiber-expense-tracker/internal/models"
)

// mockTransactionRepo implements TransactionRepo for testing
type mockTransactionRepo struct {
	transactions map[int64]*models.Transaction
	nextID       int64
}

func newMockTransactionRepo() *mockTransactionRepo {
	return &mockTransactionRepo{
		transactions: make(map[int64]*models.Transaction),
		nextID:       1,
	}
}

func (m *mockTransactionRepo) Create(tx *models.Transaction) error {
	tx.ID = m.nextID
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	m.transactions[tx.ID] = tx
	m.nextID++
	return nil
}

func (m *mockTransactionRepo) GetByID(id int64) (*models.Transaction, error) {
	tx, ok := m.transactions[id]
	if !ok {
		return nil, nil
	}
	return tx, nil
}

func (m *mockTransactionRepo) List(filter models.TransactionFilter) ([]models.Transaction, int64, error) {
	var result []models.Transaction
	for _, tx := range m.transactions {
		result = append(result, *tx)
	}
	return result, int64(len(result)), nil
}

func (m *mockTransactionRepo) Update(tx *models.Transaction) error {
	if _, ok := m.transactions[tx.ID]; !ok {
		return sql.ErrNoRows
	}
	tx.UpdatedAt = time.Now()
	m.transactions[tx.ID] = tx
	return nil
}

func (m *mockTransactionRepo) Delete(id int64) error {
	if _, ok := m.transactions[id]; !ok {
		return sql.ErrNoRows
	}
	delete(m.transactions, id)
	return nil
}

func (m *mockTransactionRepo) GetSummary(fromDate, toDate string) (*models.TransactionSummary, error) {
	summary := &models.TransactionSummary{}
	for _, tx := range m.transactions {
		if tx.Type == models.Income {
			summary.TotalIncome += tx.Amount
		} else {
			summary.TotalExpense += tx.Amount
		}
	}
	summary.Balance = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}

func TestTransactionService_Create(t *testing.T) {
	svc := NewTransactionService(newMockTransactionRepo())

	t.Run("success income", func(t *testing.T) {
		tx, err := svc.Create(models.CreateTransactionRequest{
			Type:     models.Income,
			Amount:   1000,
			Category: "salary",
			Date:     "2026-03-11",
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), tx.ID)
		assert.Equal(t, models.Income, tx.Type)
		assert.Equal(t, 1000.0, tx.Amount)
	})

	t.Run("invalid amount", func(t *testing.T) {
		_, err := svc.Create(models.CreateTransactionRequest{
			Type:     models.Expense,
			Amount:   -100,
			Category: "food",
			Date:     "2026-03-11",
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := svc.Create(models.CreateTransactionRequest{
			Type:     "invalid",
			Amount:   100,
			Category: "food",
			Date:     "2026-03-11",
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("invalid date format", func(t *testing.T) {
		_, err := svc.Create(models.CreateTransactionRequest{
			Type:     models.Expense,
			Amount:   100,
			Category: "food",
			Date:     "11-03-2026",
		})
		assert.Error(t, err)
	})
}

func TestTransactionService_GetByID(t *testing.T) {
	repo := newMockTransactionRepo()
	svc := NewTransactionService(repo)

	// Seed one transaction
	_, err := svc.Create(models.CreateTransactionRequest{
		Type: models.Expense, Amount: 500, Category: "food", Date: "2026-03-11",
	})
	assert.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		tx, err := svc.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), tx.ID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetByID(999)
		assert.ErrorIs(t, err, ErrTransactionNotFound)
	})
}

func TestTransactionService_Update(t *testing.T) {
	repo := newMockTransactionRepo()
	svc := NewTransactionService(repo)

	_, err := svc.Create(models.CreateTransactionRequest{
		Type: models.Expense, Amount: 500, Category: "food", Date: "2026-03-11",
	})
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		tx, err := svc.Update(1, models.UpdateTransactionRequest{Amount: 750})
		assert.NoError(t, err)
		assert.Equal(t, 750.0, tx.Amount)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.Update(999, models.UpdateTransactionRequest{Amount: 100})
		assert.ErrorIs(t, err, ErrTransactionNotFound)
	})
}

func TestTransactionService_Delete(t *testing.T) {
	repo := newMockTransactionRepo()
	svc := NewTransactionService(repo)

	_, err := svc.Create(models.CreateTransactionRequest{
		Type: models.Expense, Amount: 200, Category: "transport", Date: "2026-03-11",
	})
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := svc.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("already deleted", func(t *testing.T) {
		err := svc.Delete(1)
		assert.Error(t, err)
	})
}

func TestTransactionService_List_Pagination(t *testing.T) {
	repo := newMockTransactionRepo()
	svc := NewTransactionService(repo)

	for i := 0; i < 5; i++ {
		_, err := svc.Create(models.CreateTransactionRequest{
			Type: models.Expense, Amount: float64((i + 1) * 100), Category: "food", Date: "2026-03-11",
		})
		assert.NoError(t, err)
	}

	result, err := svc.List(models.TransactionFilter{Page: 1, Limit: 3})
	assert.NoError(t, err)
	assert.Equal(t, int64(5), result.TotalCount)
	assert.Equal(t, 2, result.TotalPages)
}

func TestTransactionService_GetSummary(t *testing.T) {
	repo := newMockTransactionRepo()
	svc := NewTransactionService(repo)

	_, err := svc.Create(models.CreateTransactionRequest{
		Type: models.Income, Amount: 5000, Category: "salary", Date: "2026-03-01",
	})
	assert.NoError(t, err)
	_, err = svc.Create(models.CreateTransactionRequest{
		Type: models.Expense, Amount: 1500, Category: "rent", Date: "2026-03-05",
	})
	assert.NoError(t, err)

	summary, err := svc.GetSummary("", "")
	assert.NoError(t, err)
	assert.Equal(t, 5000.0, summary.TotalIncome)
	assert.Equal(t, 1500.0, summary.TotalExpense)
	assert.Equal(t, 3500.0, summary.Balance)
}
