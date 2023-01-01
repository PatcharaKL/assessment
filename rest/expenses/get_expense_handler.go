package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpensesByIdHandler(c echo.Context) error {
	id := c.Param("id")

	stmt, err := h.DB.Prepare("SELECT * FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	row := stmt.QueryRow(id)

	e := Expenses{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
