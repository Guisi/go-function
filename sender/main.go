package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Guisi/go-function/model"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
)

const topicName = "posts"

// An app holds the clients and parsed templates that can be reused between
// requests.
type app struct {
	pubsubClient *pubsub.Client
	pubsubTopic  *pubsub.Topic
}

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("GOOGLE_CLOUD_PROJECT must be set")
	}

	a, err := newApp(projectID)
	if err != nil {
		log.Fatalf("newApp: %v", err)
	}

	for i := 1; i < 101; i++ {
		t := model.Post{
			Id:           fmt.Sprintf("id_test_%d", i),
			Message:      fmt.Sprintf("Message test %d", i),
			CreationDate: time.Now(),
		}

		log.Printf("Sending message: %s", t)

		msg, err := json.Marshal(t)
		if err != nil {
			log.Printf("json.Marshal: %v", err)
			return
		}

		ctx := context.Background()

		res := a.pubsubTopic.Publish(ctx, &pubsub.Message{Data: msg})
		if _, err := res.Get(ctx); err != nil {
			log.Printf("Publish.Get: %v", err)
			return
		}
	}
}

// newApp creates a new app.
func newApp(projectID string) (*app, error) {
	ctx := context.Background()

	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	pubsubTopic := pubsubClient.Topic(topicName)

	return &app{
		pubsubClient: pubsubClient,
		pubsubTopic:  pubsubTopic,
	}, nil
}
