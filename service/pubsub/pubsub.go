package pubsub

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/tidwall/gjson"
	"google.golang.org/api/option"
)

type Config struct {
	GcpCredential string
}

func (c *Config) NewClient() (*pubsub.Client, error) {
	// open config
	b, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(c.GcpCredential)
	if err != nil {
		log.Printf("failed load credential config : %v", err)
		return nil, err
	}
	projectId := gjson.Get(string(b), "project_id").String()
	// new client
	ctx := context.Background()
	psClient, err := pubsub.NewClient(ctx, projectId, option.WithCredentialsJSON(b))
	if err != nil {
		log.Printf("failed to connect Pubsub Client : %v", err)
		return nil, err
	}
	log.Println("service Pub/Sub started !!!")
	return psClient, nil
}

//publish topic with custom message
func (c *Config) Publish(ctx context.Context, message pubsub.Message, topicName string) error {
	client, err := c.NewClient()
	if err != nil {
		return err
	}
	topic := client.Topic(topicName)

	defer func() {
		topic.Stop()
		_ = client.Close()
	}()

	pubRes := topic.Publish(ctx, &message)
	bytes, err := json.Marshal(message)
	fmt.Println("Pubsub message : " + string(bytes))
	if _, err := pubRes.Get(ctx); err != nil {
		return err
	}
	return nil
}
