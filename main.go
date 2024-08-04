package main

import (
	"log"
	"message-system/app/cache/redis"
	"message-system/app/client"
	"message-system/app/handler"
	"strconv"

	"message-system/app/repository/mongodb"
	"message-system/app/service"
	"message-system/config"
)

func main() {
	conf, err := config.New(".env", "local")
	if err != nil {
		log.Fatal(err)
	}
	conf.Print()

	repo, err := mongodb.NewMessageRepository(conf.MongoDBConfig)
	if err != nil {
		log.Fatal(err)
	}

	cache := redis.NewCache(conf.RedisConfig)
	webHookClient := client.NewWebhookClient(conf.WebhookConfig)
	service := service.NewService(cache, webHookClient, repo)
	messageHandler := handler.NewHandler(service)
	echoServer := handler.RegisterRoutes(messageHandler)

	log.Fatal(echoServer.Start(":" + strconv.Itoa(conf.Port)))
}
