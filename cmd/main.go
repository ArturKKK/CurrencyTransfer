package main

import (
	"github.com/ArturKKK/CurrencyTransfer/internal/config"
	"github.com/ArturKKK/CurrencyTransfer/internal/db"
	"github.com/ArturKKK/CurrencyTransfer/internal/handler"
	"github.com/ArturKKK/CurrencyTransfer/pkg/logging"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	url = "http://www.cbr.ru/scripts/XML_daily.asp"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	postgres, err := db.NewDatabase(&cfg.Postgres)
	if err != nil{
		log.Fatal(err)
	}
	logger.Info("postgres started")

	err = postgres.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("postgres initialized")

	logger.Info("start parsing")
	Parse(url, postgres)

	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				Parse(url, postgres)
				logger.Info("tick-tick")
			}
		}
	}()

	handler := handler.NewHander(postgres)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.Listen.Port),
		Handler: handler.Router(),
	}

	go func() {
		logger.Info("service started")
		_ = srv.ListenAndServe()
	}()

	// wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

	done <- true // stop the goroutine 

}