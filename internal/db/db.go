package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ArturKKK/CurrencyTransfer/pkg/logging"
	"github.com/ArturKKK/CurrencyTransfer/pkg/utils"
	_ "github.com/lib/pq"
)

type Database struct {
	client *sql.DB
	logger *logging.Logger
}

func NewDatabase(config *Config, logger *logging.Logger) (*Database, error) {
	connInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, config.SslMode,
	)
	logger.Infof("connStr for postgres: %s", connInfo)

	var client *sql.DB
	var err error

	utils.DoWithTries(func() error {
		client, err = sql.Open("postgres", connInfo)
		if err != nil {
			logger.Errorf("fatal open postgres: %v", err)
			return err
		}
		logger.Info("postgres opened")
		err = nil
		return nil
	}, config.MaxAttempts, 5*time.Second)

	if err != nil {
		logger.Error("failed to connect to postgres")
		return nil, err
	}

	utils.DoWithTries(func() error {
		if err := client.Ping(); err != nil {
			logger.Errorf("fatal ping postgres: %v", err)
			return err
		}
		logger.Info("postgres pinged")
		err = nil
		return nil
	}, config.MaxAttempts, 5*time.Second)

	if err != nil {
		logger.Error("failed to ping postgres")
		return nil, err
	}

	return &Database{
		client: client,
		logger: logger,
	}, nil
}

func (db *Database) Init(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, initRequest)
	return err
}

// Drop drops all tables in database.
func (db *Database) Drop(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, dropRequest)
	return err
}

// Clean cleans all tables in database.
func (db *Database) Clean(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, cleanRequest)
	return err
}

func (db *Database) Save(ctx context.Context, charCode string, value float64) error {
	_, err := db.client.ExecContext(ctx, saveRequest, charCode, value)
	return err
}

func (db *Database) SelectByCharCode(ctx context.Context, charCode string) (*Currency, error) {
	row := db.client.QueryRowContext(ctx, selectByCharCode, charCode)

	var curr Currency
	err := row.Scan(&curr.CharCode, &curr.VunitRate)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrCharCodeNotFound
	case err != nil:
		return nil, err
	}

	return &curr, nil
}

func (db *Database) SelectByValue(ctx context.Context, value float64) (*Currency, error) {
	row := db.client.QueryRowContext(ctx, selectByValue, value)

	var curr Currency
	err := row.Scan(&curr.CharCode, &curr.VunitRate)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrValueNotFound
	case err != nil:
		return nil, err
	}

	return &curr, nil
}

func (db *Database) GetCurrencies(ctx context.Context) ([]Currency, error) {
	rows, err := db.client.QueryContext(ctx, selectCurrencyArr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var currencies []Currency
	for rows.Next() {
		var currency Currency
		if err = rows.Scan(&currency.CharCode, &currency.VunitRate); err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return currencies, nil
}

func (db *Database) GetCurrencyRate(ctx context.Context, charcode string) (float64, error) {
	row := db.client.QueryRowContext(ctx, selectByCharCode, charcode)

	var currency Currency
	err := row.Scan(&currency.CharCode, &currency.VunitRate)

	switch {
	case err == sql.ErrNoRows:
		return 0, ErrCharCodeNotFound
	case err != nil:
		return 0, err
	}

	return currency.VunitRate, nil
}
