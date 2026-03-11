package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"gofiber-expense-tracker/internal/models"
	"gofiber-expense-tracker/internal/services"
)

type TransactionHandler struct {
	svc *services.TransactionService
}

func NewTransactionHandler(svc *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// Create godoc
// @Summary Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body models.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} models.Transaction
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var req models.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	tx, err := h.svc.Create(req)
	if err != nil {
		if err == services.ErrInvalidInput {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tx)
}

// GetByID godoc
// @Summary Get transaction by ID
// @Tags transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid transaction ID",
		})
	}

	tx, err := h.svc.GetByID(id)
	if err != nil {
		if err == services.ErrTransactionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(tx)
}

// List godoc
// @Summary List transactions with filters
// @Tags transactions
// @Produce json
// @Param type query string false "Transaction type (income/expense)"
// @Param category query string false "Category filter"
// @Param from_date query string false "From date (YYYY-MM-DD)"
// @Param to_date query string false "To date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} models.PaginatedResponse
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) List(c *fiber.Ctx) error {
	var filter models.TransactionFilter
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid query parameters",
		})
	}

	result, err := h.svc.List(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(result)
}

// Update godoc
// @Summary Update a transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body models.UpdateTransactionRequest true "Updated data"
// @Success 200 {object} models.Transaction
// @Router /api/v1/transactions/{id} [put]
func (h *TransactionHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid transaction ID",
		})
	}

	var req models.UpdateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	tx, err := h.svc.Update(id, req)
	if err != nil {
		if err == services.ErrTransactionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(tx)
}

// Delete godoc
// @Summary Delete a transaction
// @Tags transactions
// @Param id path int true "Transaction ID"
// @Success 204
// @Router /api/v1/transactions/{id} [delete]
func (h *TransactionHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid transaction ID",
		})
	}

	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetSummary godoc
// @Summary Get financial summary
// @Tags transactions
// @Produce json
// @Param from_date query string false "From date (YYYY-MM-DD)"
// @Param to_date query string false "To date (YYYY-MM-DD)"
// @Success 200 {object} models.TransactionSummary
// @Router /api/v1/transactions/summary [get]
func (h *TransactionHandler) GetSummary(c *fiber.Ctx) error {
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	summary, err := h.svc.GetSummary(fromDate, toDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(summary)
}
