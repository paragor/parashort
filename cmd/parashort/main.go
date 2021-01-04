package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/paragor/parashort/pkg/adapters/app_config"
	"github.com/paragor/parashort/pkg/adapters/http_web"
	"github.com/paragor/parashort/pkg/adapters/redis_storage"
	"github.com/paragor/parashort/pkg/domain/parashort"
	"github.com/paragor/parashort/pkg/domain/storage"
)

func main() {
	help, err := app_config.ShowConfigHelp()
	if err != nil {
		panic(err)
	}
	fmt.Println(help)
	config, err := app_config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	engine, err := createStorage(config, errChan)
	if err != nil {
		panic(err)
	}
	shortApp := parashort.NewParashortApp(config.CalcAppTimeout(), engine)
	webApp := http_web.NewWebServer(shortApp, config.TemplateDir, config.AssetsDir)

	go func() {
		errChan <- webApp.Run(ctx)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("SIGTERM...")
		cancel()
	}()

	err = <-errChan
	cancel()
	if err != nil {
		log.Println(err)
	}
	log.Fatal("exit...")
}

func createStorage(config *app_config.AppConfig, errChan chan error) (storage.StorageEngine, error) {
	redisClient := createRedisClient(config)
	redisStorage := redis_storage.NewRedisStorage(
		redisClient,
	)
	go func() {
		timer := time.NewTimer(time.Second * 10)
		for {
			<-timer.C
			timeout, _ := context.WithTimeout(context.Background(), config.CalcAppTimeout())
			err := redisStorage.Ping(timeout)
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	return redisStorage, redisStorage.Ping(context.Background())
}

func createRedisClient(config *app_config.AppConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         config.RedisAddr,
		Username:     "",
		Password:     "",
		DB:           0,
		DialTimeout:  config.CalcAppTimeout(),
		ReadTimeout:  config.CalcAppTimeout(),
		WriteTimeout: config.CalcAppTimeout(),
	})
}
