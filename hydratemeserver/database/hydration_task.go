package database

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// HydrationTask represents a notification configuration for a specific slack user
type HydrationTask struct {
	ID                    *datastore.Key `datastore:"__key__"`
	DateCreated           time.Time
	SlackUserID           string
	SlackWorkspaceID      string
	AlertFrequencyMinutes int
	LastHydration         time.Time
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

	return parseTasksFromQueryResult(it)
}

const hydrationSpacing = 5.0 * time.Minute

// GetOverdueHydrationTasks gets all hydration tasks that are overdue
func GetOverdueHydrationTasks() ([]HydrationTask, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, datastore.DetectProjectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	query := datastore.NewQuery(kindHydrationTask).
		Filter("LastHydration <", time.Now().Add(-hydrationSpacing))
	it := client.Run(ctx, query)

	return parseTasksFromQueryResult(it)
}

// SetWasHydrated sets the task with id as LastHydrated right now
func SetWasHydrated(taskKey *datastore.Key) error {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, datastore.DetectProjectID)
	if err != nil {
		return err
	}
	defer client.Close()

	tx, err := client.NewTransaction(ctx)
	if err != nil {
		return err
	}
	var task HydrationTask
	if err := tx.Get(taskKey, &task); err != nil {
		return err
	}
	task.LastHydration = time.Now()
	if _, err := tx.Put(taskKey, &task); err != nil {
		return err
	}
	if _, err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func parseTasksFromQueryResult(it *datastore.Iterator) ([]HydrationTask, error) {
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
