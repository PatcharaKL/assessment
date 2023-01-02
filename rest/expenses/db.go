package expenses

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	createExpensesTableSQL = `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	createExpenseSQL = "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id;"
	getExpensesSQL   = "SELECT * FROM expenses"
	getExpenseSQL    = "SELECT * FROM expenses WHERE id = $1"
	updateExpenseSQL = "UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1"
)

func InitDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_STR"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	_, err = db.Exec(createExpensesTableSQL)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	return db
}
