package consumer

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Clients reused between function invocations.
var (
	firestoreClient *firestore.Client
)

// PubSubMessage is the payload of a Pub/Sub event.
// See https://cloud.google.com/functions/docs/calling/pubsub.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Post struct {
	Id           string    `json:"id"`
	Message      string    `json:"message"`
	CreationDate time.Time `json:"creationDate"`
}

// initializeClients creates translateClient and firestoreClient if they haven't been created yet.
func initializeClients() error {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return fmt.Errorf("GOOGLE_CLOUD_PROJECT must be set")
	}

	if firestoreClient == nil {
		// Pre-declare err to avoid shadowing firestoreClient.
		var err error
		// Use context.Background() so the client can be reused.
		firestoreClient, err = firestore.NewClient(context.Background(), projectID)
		if err != nil {
			return fmt.Errorf("firestore.NewClient: %v", err)
		}
	}
	return nil
}

func SavePost(ctx context.Context, m PubSubMessage) error {
	initializeClients()

	post := Post{}
	if err := json.Unmarshal(m.Data, &post); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	_, err := firestoreClient.Collection("posts").Doc(post.Id).Set(ctx, post)
	if err != nil {
		return fmt.Errorf("failed adding post: %v", err)
	}

	return nil
}
