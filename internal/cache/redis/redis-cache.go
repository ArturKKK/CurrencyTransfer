package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ArturKKK/CurrencyTransfer/internal/db"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	config *Config
	client *redis.Client
}

func GetClient(config *Config) *RedisCache {
	//Standart port: 6379
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalf("redis: %v", err)
	}

	return &RedisCache{
		config: config,
		client: client,
	}
}

func (r *RedisCache) SetOne(key string, vuniteRate float64) error {
	marshalVunite, err := json.Marshal(vuniteRate)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return curr, nil
}

func (r *RedisCache) Set(currencies []db.Currency) error {
	marshalCurrencies, err := json.Marshal(currencies)
	if err != nil {
		log.Fatal(err)
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
		return nil, err
	}

	return currs, nil
}
