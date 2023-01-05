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

func setupTestServer(method, uri string, body *bytes.Buffer) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(method, uri, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return rec, c
}
func TestCreateExpenseU(t *testing.T) {
	successRes := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"
	badRequestRes := "{\"message\":\"code=400, message=Syntax error: offset=115, error=invalid character '}' looking for beginning of object key string, internal=invalid character '}' looking for beginning of object key string\"}"
	InternalServerErrorRes := "{\"message\":\"all expectations were already fulfilled, call to Query 'INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id;' with args [{Name: Ordinal:1 Value:strawberry smoothie} {Name: Ordinal:2 Value:79} {Name: Ordinal:3 Value:night market promotion discount 10 bath} {Name: Ordinal:4 Value:{\\\"food\\\",\\\"beverage\\\"}}] was not expected\"}"

	tests := []struct {
		name         string
		body         *bytes.Buffer
		expectedRes  string
		expectedCode int
	}{
		{
			name: "testSucceed",
			body: bytes.NewBufferString(`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": "night market promotion discount 10 bath",
				"tags": ["food", "beverage"]
			}`),
			expectedRes:  successRes,
			expectedCode: http.StatusCreated,
		},
		{
			name: "testBadRequest",
			body: bytes.NewBufferString(`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": "night market promotion discount 10 bath",
			}`),
			expectedRes:  badRequestRes,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testInternalServerError",
			body: bytes.NewBufferString(`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": "night market promotion discount 10 bath",
				"tags": ["food", "beverage"]
			}`),
			expectedRes:  InternalServerErrorRes,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rec, c := setupTestServer(http.MethodPost, "/expenses", tt.body)

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
			if tt.name != "testInternalServerError" {
				mock.ExpectQuery("INSERT INTO expenses").WithArgs("strawberry smoothie", 79.00, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).WillReturnRows(expectedRow)
			}
			h := Handler{db}

			// Act
			err = h.CreateExpensesHandler(c)

			// Assertions
			fmt.Println(strings.TrimSpace(rec.Body.String()))
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedRes, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}

func TestGetExpenseByIDU(t *testing.T) {
	successRes := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"
	InternalServerErrorRes := "{\"message\":\"all expectations were already fulfilled, call to Query 'SELECT * FROM expenses WHERE id = $1' with args [{Name: Ordinal:1 Value:1}] was not expected\"}"

	tests := []struct {
		name         string
		expectedRes  string
		expectedCode int
	}{
		{
			name:         "testSucceed",
			expectedRes:  successRes,
			expectedCode: http.StatusOK,
		},
		{
			name:         "testInternalServerError",
			expectedRes:  InternalServerErrorRes,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		// Arrange
		rec, c := setupTestServer(http.MethodGet, "/expenses", bytes.NewBufferString(``))
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
		if tt.name != "testInternalServerError" {
			mock.ExpectQuery("SELECT \\* FROM expenses WHERE id = \\$1").WithArgs("1").WillReturnRows(expectedRow)
		}
		h := Handler{db}

		// Act
		err = h.GetExpenseByIdHandler(c)

		// Assertions
		fmt.Println(strings.TrimSpace(rec.Body.String()))
		if assert.NoError(t, err) {
			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Equal(t, tt.expectedRes, strings.TrimSpace(rec.Body.String()))
		}
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
	mock.ExpectPrepare("UPDATE expenses SET (.+) WHERE (.+)").ExpectExec().WithArgs("1", "apple smoothie", 89.00, "no discount", pq.Array([]string{"beverage"})).WillReturnResult(sqlmock.NewResult(0, 0))
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
	successRes := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]},{\"id\":2,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}]"
	prepareStmtErrorRes := "{\"message\":\"can't prepare query all expenses statement:all expectations were already fulfilled, call to Prepare 'SELECT * FROM expenses' query was not expected\"}"
	queryStmtErrorRes := "{\"message\":\"can't query expenses: all expectations were already fulfilled, call to Query 'SELECT * FROM expenses' with args [] was not expected\"}"
	scanErrorRes := "{\"message\":\"can't scan user:sql: Scan error on column index 4, name \\\"tags\\\": pq: unable to parse array; expected '{' at offset 0\"}"

	tests := []struct {
		name         string
		expectedRes  string
		expectedCode int
	}{
		{
			name:         "testSucceed",
			expectedRes:  successRes,
			expectedCode: http.StatusOK,
		},
		{
			name:         "testPrepareStmtError",
			expectedRes:  prepareStmtErrorRes,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "testQueryStmtError",
			expectedRes:  queryStmtErrorRes,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "testScanError",
			expectedRes:  scanErrorRes,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		// Arrange
		rec, c := setupTestServer(http.MethodGet, "/expenses", bytes.NewBufferString(``))

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// Set up mock rows to return when querying
		var expectedRow *sqlmock.Rows
		if tt.name != "testScanError" {
			expectedRow = sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
				AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).
				AddRow(2, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"}))
		} else {
			expectedRow = sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
				AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", "food").
				AddRow(2, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"}))
		}
		if tt.name != "testPrepareStmtError" {
			// Set up mock to expect a query and return mock rows
			expectedPrepare := mock.ExpectPrepare("SELECT \\* FROM expenses")
			if tt.name != "testQueryStmtError" {
				expectedPrepare.ExpectQuery().WillReturnRows(expectedRow)
			}
		}
		h := Handler{db}

		// Act
		err = h.GetExpensesHandler(c)

		// Assertions
		fmt.Println(strings.TrimSpace(rec.Body.String()))
		if assert.NoError(t, err) {
			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Equal(t, tt.expectedRes, strings.TrimSpace(rec.Body.String()))
		}
	}

	// // Arrange
	// rec, c := setupTestServer(http.MethodGet, "/expenses", bytes.NewBufferString(``))

	// db, mock, err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	// }
	// defer db.Close()

	// // Set up mock rows to return when querying
	// expectedRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
	// 	AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).
	// 	AddRow(2, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"}))

	// // Set up mock to expect a query and return mock rows
	// expectedPrepare := mock.ExpectPrepare("SELECT \\* FROM expenses")
	// expectedPrepare.ExpectQuery().WillReturnRows(expectedRow)
	// h := Handler{db}
	// expected := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]},{\"id\":2,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}]"

	// // Act
	// err = h.GetExpensesHandler(c)

	// // Assertions
	// fmt.Println(strings.TrimSpace(rec.Body.String()))
	// if assert.NoError(t, err) {
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	// }
}
