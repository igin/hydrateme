package database

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// HydrationTask represents a notification configuration for a specific slack user
type HydrationTask struct {
	DateCreated           time.Time
	SlackUserID           string
	SlackWorkspaceID      string
	AlertFrequencyMinutes int
	StartTime             time.Time
	EndTime               time.Time
}

const kindHydrationTask = "HydrationTask"

// CreateHydrationTask creates a new hydration task in the database
func CreateHydrationTask(userID string) (*datastore.Key, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, datastore.DetectProjectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	task := &HydrationTask{
		SlackUserID: userID,
	}
	key := datastore.IncompleteKey(kindHydrationTask, nil)
	return client.Put(ctx, key, task)
}

// GetHydrationTasksOfUser gets all hydration tasks of a user
func GetHydrationTasksOfUser(userID string) ([]HydrationTask, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, datastore.DetectProjectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	query := datastore.NewQuery(kindHydrationTask).
		Filter("SlackUserID =", userID)
	it := client.Run(ctx, query)

	tasks := make([]HydrationTask, 0)
	for {
		var task HydrationTask
		_, err := it.Next(&task)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
