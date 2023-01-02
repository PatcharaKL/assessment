package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *Handler) GetExpenseByIdHandler(c echo.Context) error {
	id := c.Param("id")

	row := h.DB.QueryRow(getExpenseSQL, id)

	e := Expenses{}

	if err := row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}

func (h *Handler) GetExpensesHandler(c echo.Context) error {
	stmt, err := h.DB.Prepare(getExpensesSQL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all expenses statement:" + err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query expenses: " + err.Error()})
	}

	expenses := []Expenses{}
	for rows.Next() {
		e := Expenses{}
		err = rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
		}
		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)
}
