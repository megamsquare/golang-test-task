package main

import (
	"encoding/json"
	"net/http"
	"twitch_chat_analysis/model"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "worked")
	})

	r.POST("/message", func(ctx *gin.Context) {
		var body model.MessageBody
		err := ctx.ShouldBindJSON(&body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "message sent"})
	})

	r.Run(":8080")
}
