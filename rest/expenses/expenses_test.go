//go:build unit
// +build unit

package expenses

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseU(t *testing.T) {
	// Arrange
	e := echo.New()
	body := bytes.NewBufferString(`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": "night market promotion discount 10 bath",
				"tags": ["food", "beverage"]
			}`)
	req := httptest.NewRequest(http.MethodPost, "/expenses", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock rows to return when querying
	expectedID := 1
	expectedRow := sqlmock.NewRows([]string{"id"}).
		AddRow(expectedID)

	// Set up mock to expect a query and return mock rows
	mock.ExpectQuery("INSERT INTO expenses").WithArgs("strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).WillReturnRows(expectedRow)
	h := Handler{db}
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	// Act
	err = h.CreateExpensesHandler(c)

	// Assertions
	fmt.Println(strings.TrimSpace(rec.Body.String()))
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetExpenseByIDU(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/1", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock rows to return when querying
	expectedRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"}))

	// Set up mock to expect a query and return mock rows
	mock.ExpectQuery("SELECT \\* FROM expenses WHERE id = \\$1").WithArgs("1").WillReturnRows(expectedRow)
	h := Handler{db}
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	// Act
	err = h.GetExpenseByIdHandler(c)

	// Assertions
	fmt.Println(strings.TrimSpace(rec.Body.String()))
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdateExpenseU(t *testing.T) {
	// Arrange
	e := echo.New()
	body := bytes.NewBufferString(`{
		"id": 1,
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)
	req := httptest.NewRequest(http.MethodPut, "/expenses/1", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock rows to return when querying
	sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"}))

	// Set up mock to expect a query and return mock rows
	mock.ExpectPrepare("UPDATE expenses SET (.+) WHERE (.+)").ExpectExec().WithArgs("1", "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"})).WillReturnResult(sqlmock.NewResult(0, 0))
	h := Handler{db}
	expected := "{\"id\":1,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}"

	// Act
	err = h.UpdateExpensesHandler(c)

	// Assertions
	fmt.Println(strings.TrimSpace(rec.Body.String()))
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetExpensesU(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses1", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock rows to return when querying
	expectedRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).
		AddRow(2, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"}))

	// Set up mock to expect a query and return mock rows
	mock.ExpectPrepare("SELECT \\* FROM expenses").ExpectQuery().WillReturnRows(expectedRow)
	h := Handler{db}
	expected := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]},{\"id\":2,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}]"

	// Act
	err = h.GetExpensesHandler(c)

	// Assertions
	fmt.Println(strings.TrimSpace(rec.Body.String()))
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
