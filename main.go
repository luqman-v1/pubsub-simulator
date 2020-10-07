package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	ps "pubsub-simulator/service/pubsub"

	"github.com/tidwall/gjson"

	"cloud.google.com/go/pubsub"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	router := gin.Default()
	router.LoadHTMLGlob("templates/index.html")
	router.Static("/asset", "./templates/asset")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.POST("/ajax/send", func(c *gin.Context) {
		message := c.PostForm("message")
		topicName := c.PostForm("topic_name")
		credential := c.PostForm("base64_credential")

		data := gjson.Get(message, "data").Value()
		bytes, _ := json.Marshal(data)
		attributes := gjson.Get(message, "attributes").Map()
		atr := make(map[string]string)
		for k, v := range attributes {
			atr[k] = v.String()
		}
		p := ps.Config{
			GcpCredential: credential,
		}
		payload := pubsub.Message{
			Data:       bytes,
			Attributes: atr,
		}
		err := p.Publish(ctx, payload, topicName)
		httpResponse := 200
		messageResponse := "Sucess Publish"
		if err != nil {
			log.Println("Failed Publish", err)
			messageResponse = err.Error()
		}
		c.JSON(httpResponse, gin.H{
			"message": messageResponse,
		})
	})
	_ = router.Run(":8080")
}
