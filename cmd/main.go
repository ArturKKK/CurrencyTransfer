package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	rediscache "github.com/ArturKKK/CurrencyTransfer/internal/cache/redis"
	"github.com/ArturKKK/CurrencyTransfer/internal/config"
	"github.com/ArturKKK/CurrencyTransfer/internal/db"
	"github.com/ArturKKK/CurrencyTransfer/internal/handler"
	"github.com/ArturKKK/CurrencyTransfer/internal/parser"
	"github.com/ArturKKK/CurrencyTransfer/pkg/logging"
)

const (
	url = "http://www.cbr.ru/scripts/XML_daily.asp"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	cache, err := rediscache.GetClient(&cfg.Redis, logger)
	if err != nil {
		logger.Fatalf("failed to connect to redis")
	}
	logger.Info("redis started")

	postgres, err := db.NewDatabase(&cfg.Postgres, logger)
	if err != nil {
		logger.Fatalf("failed to connect to the database: %v", err)
	}
	logger.Info("postgres started")

	err = postgres.Init(ctx)
	if err != nil {
		logger.Fatal("failed to initialize database")
	}
	logger.Info("postgres initialized")

	logger.Info("start parsing")
	err = parser.Parse(url, postgres, logger)
	if err != nil {
		logger.Errorf("failed to parse")
	}
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				logger.Info("start parsing")
				err = parser.Parse(url, postgres, logger)
				if err != nil {
					logger.Errorf("failed to parse")
				} else {
				logger.Info("parsing completed successfully")
				}
			}
		}
	}()

	handler := handler.NewHander(postgres, cache, logger)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Listen.Port),
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
