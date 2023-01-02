package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *Handler) CreateExpensesHandler(c echo.Context) error {
	e := Expenses{}

	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if err := h.DB.QueryRow(createExpenseSQL, e.Title, e.Amount, e.Note, pq.Array(e.Tags)).Scan(&e.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
