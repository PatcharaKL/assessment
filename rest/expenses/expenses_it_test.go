//go:build integration
// +build integration

package expenses

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// func seedExpenses(t *testing.T) Expenses {
// 	var e Expenses
// 	body := bytes.NewBufferString(`{
// 		"title": "strawberry smoothie",
// 		"amount": 79,
// 		"note": "night market promotion discount 10 bath",
// 		"tags": ["food", "beverage"]
// 	}`)

// 	err := request(http.MethodPost, uri("expenses"), body).Decode(&e)
// 	if err != nil {
// 		t.Fatal("Can't create user: ", err)
// 	}
// 	return e
// }

func TestCreateExpenses(t *testing.T) {
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
	assert.Equal(t, 79, e.Amount)
}

func uri(path ...string) string {
	host := "http://localhost:2565"
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
	req.Header.Add("Authorization", "Basic UGF0Y2hhcmE6UGFzc3dvcmQ=")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

// const serverPort = 80

// func setupServer() {
// 	eh := echo.New()
// 	go func(e *echo.Echo) {
// 		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		h := NewApplication(db)

// 		e.POST("/expenses", h.CreateExpensesHandler)
// 		e.Start(fmt.Sprintf(":%d", serverPort))
// 	}(eh)
// 	for {
// 		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		if conn != nil {
// 			conn.Close()
// 			break
// 		}
// 	}
// }
