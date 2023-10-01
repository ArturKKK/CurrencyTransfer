package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ArturKKK/CurrencyTransfer/internal/db"
	"github.com/ArturKKK/CurrencyTransfer/pkg/logging"
	"github.com/ArturKKK/CurrencyTransfer/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	config *Config
	client *redis.Client
	logger *logging.Logger
}

func GetClient(config *Config, logger *logging.Logger) (*RedisCache, error) {
	connInfo := fmt.Sprintf("%s:%s", config.Host, config.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     connInfo,
		Password: config.Password,
		DB:       config.DB,
	})

	logger.Infof("connStr for redis: %s", connInfo)

	var err error
	utils.DoWithTries(func() error {
		_, err = client.Ping(context.TODO()).Result()
		if err != nil {
			logger.Errorf("failed ping redis: %v", err)
			return err
		}
		logger.Errorf("redis pinged: %v", err)
		err = nil
		return nil
	}, config.MaxAttempts, 5*time.Second)

	if err != nil {
		logger.Errorf("failed to connect to redis: %v", err)
		return nil, err
	}

	return &RedisCache{
		config: config,
		client: client,
		logger: logger,
	}, nil
}

func (r *RedisCache) SetOne(key string, vuniteRate float64) error {
	marshalVunite, err := json.Marshal(vuniteRate)
	if err != nil {
		r.logger.Errorf("failed to marshal: %v", err)
		return err
	}

	r.client.Set(context.TODO(), key, marshalVunite, time.Duration(r.config.Expires)*time.Second)
	return nil
}

func (r *RedisCache) GetOne(key string) (float64, error) {
	val, err := r.client.Get(context.TODO(), key).Result()
	if err != nil {
		return 0, err
	}

	var curr float64
	err = json.Unmarshal([]byte(val), &curr)
	if err != nil {
		r.logger.Errorf("failed to unmarshall: %v", err)
		return 0, err
	}

	return curr, nil
}

func (r *RedisCache) Set(currencies []db.Currency) error {
	marshalCurrencies, err := json.Marshal(currencies)
	if err != nil {
		r.logger.Errorf("failed to marshal: %v", err)
		return err
	}

	r.client.Set(context.TODO(), "All", marshalCurrencies, time.Duration(r.config.Expires)*time.Second)
	return nil
}

func (r *RedisCache) Get() ([]db.Currency, error) {
	vals, err := r.client.Get(context.TODO(), "All").Result()
	if err != nil {
		return nil, err
	}

	var currs []db.Currency
	err = json.Unmarshal([]byte(vals), &currs)
	if err != nil {
		r.logger.Errorf("failed to unmarshal: %v", err)
		return nil, err
	}

	return currs, nil
}
