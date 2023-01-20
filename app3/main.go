package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"twitch_chat_analysis/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const prefix = "message"

func main() {
	r := gin.Default()

	r.GET("/message/list", func(ctx *gin.Context) {

		redisClient := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		// defer redisClient.Close()

		messages := make([]model.MessageBody, 0)

		keys, err := redisClient.Keys(fmt.Sprintf("%s:*", prefix)).Result()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, k := range keys {
			bodyString, err := redisClient.Get(k).Result()
			if err != nil {
				log.Printf("could not get key %v", err)
				continue
			}

			var b model.MessageBody

			err = json.Unmarshal([]byte(bodyString), &b)
			if err != nil {
				log.Printf("could not unmarshal body %v", err)
				continue
			}

			messages = append(messages, b)

		}

		ctx.JSON(http.StatusOK, gin.H{"messages": messages})
	})

	r.Run(":8082")
}
