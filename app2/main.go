package main

import (
	"encoding/json"
	"fmt"
	"log"
	"twitch_chat_analysis/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

const prefix = "message"

func main() {

	r := gin.Default()

	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq, %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel, %v", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to publish a message %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to consume a message %v", err)
	}

	go func() {
		for msg := range msgs {
			bodyBytes := msg.Body

			var body model.MessageBody
			err := json.Unmarshal(bodyBytes, &body)
			if err != nil {
				log.Printf("could not unmarshal body %v", err)
				continue
			}

			redisClient := redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			})

			// defer redisClient.Close()

			key := fmt.Sprintf("%s:%s", prefix, body.Message)

			err = redisClient.Set(key, string(bodyBytes), 0).Err()
			if err != nil {
				log.Printf("could not set key %v", err)
				continue
			}

			log.Printf("Successfully saved message in redis \n")

		}
	}()

	// ctx.JSON(http.StatusOK, gin.H{"message": "successful"})

	r.Run(":8081")
}
