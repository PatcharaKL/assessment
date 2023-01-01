package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PatcharaKL/assessment/rest/expenses"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func main() {
	h := expenses.NewApplication(expenses.InitDB())
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "Patchara" && password == "Password" {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", healthHandler)
	e.GET("/expenses/:id", h.GetExpensesByIdHandler)
	e.POST("/expenses", h.CreateExpensesHandler)

	go func() {
		if err := e.Start(":2565"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down server")
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Print("Server shuted down")
}
