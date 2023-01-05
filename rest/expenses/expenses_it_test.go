//go:build integration
// +build integration

package expenses

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func init() {
	setupServer()
}

func TestCreateExpenseIn(t *testing.T) {
	var e Expenses
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, 79.00, e.Amount)
}

func TestGetExpensesIn(t *testing.T) {
	var e []Expenses

	res := request(http.MethodGet, uri("expenses"), strings.NewReader(""))
	err := res.Decode(&e)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, e[0].ID)
	assert.Equal(t, "strawberry smoothie", e[0].Title)
	assert.Equal(t, 79.00, e[0].Amount)
}

func TestGetExpenseByIDIn(t *testing.T) {
	var e Expenses

	res := request(http.MethodGet, uri("expenses/1"), strings.NewReader(""))
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, 79.00, e.Amount)
}

func TestUpdateExpenseIn(t *testing.T) {
	var e Expenses

	body := bytes.NewBufferString(`{
		"id": 1,
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)

	res := request(http.MethodPut, uri("expenses/1"), body)
	err := res.Decode(&e)

	if assert.Nil(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, 1, e.ID)
		assert.Equal(t, "apple smoothie", e.Title)
		assert.Equal(t, 89.00, e.Amount)
	}
}

func uri(path ...string) string {
	host := "http://localhost:80"
	if path == nil {
		return host
	}
	url := append([]string{host}, path...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

const serverPort = 80

func setupServer() {
	eh := echo.New()

	go func(e *echo.Echo) {
		db := initTestDatabase()
		testsEndpoint(e, NewApplication(db))
	}(eh)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
}

func testsEndpoint(e *echo.Echo, h *Handler) {
	e.POST("/expenses", h.CreateExpensesHandler)
	e.GET("/expenses", h.GetExpensesHandler)
	e.GET("/expenses/:id", h.GetExpenseByIdHandler)
	e.PUT("/expenses/:id", h.UpdateExpensesHandler)
	e.Start(fmt.Sprintf(":%d", serverPort))
}

func initTestDatabase() *sql.DB {
	db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
