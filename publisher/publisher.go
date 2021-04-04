package publisher

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const topicName = "posts"

type Post struct {
	Id           string    `json:"id"`
	Message      string    `json:"message"`
	CreationDate time.Time `json:"creationDate"`
}

func Publish(w http.ResponseWriter, r *http.Request) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("GOOGLE_CLOUD_PROJECT must be set")
	}

	ctx := context.Background()
	pubsubClient, e := pubsub.NewClient(ctx, projectID)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}

	decoder := json.NewDecoder(r.Body)

	var post *Post
	err := decoder.Decode(&post)
	if err != nil {
		panic(err)
	}

	pubsubTopic := pubsubClient.Topic(topicName)

	msg, err := json.Marshal(post)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
		http.Error(w, "Error sending post", http.StatusInternalServerError)
		return
	}

	log.Printf("Post: %s", msg)

	result := pubsubTopic.Publish(ctx, &pubsub.Message{Data: msg})
	id, err := result.Get(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}
