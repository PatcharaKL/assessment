package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpensesHandler(c echo.Context) error {
	id := c.Param("id")
	e := Expenses{}
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := h.DB.Prepare("UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare update expense statement:" + err.Error()})
	}

	if _, err := stmt.Exec(id, e.Title, e.Amount, e.Note, pq.Array(e.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't update expense data:" + err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}
