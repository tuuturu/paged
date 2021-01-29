package upper

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/tuuturu/paged/pkg/core"
	"github.com/tuuturu/paged/pkg/core/models"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

func (c *upperClient) AddEvent(event *models.Event) error {
	collection := c.Session.Collection(eventTable)

	_, err := collection.Insert(event)
	if err != nil {
		return fmt.Errorf("error inserting event: %w", err)
	}

	return nil
}

func (c *upperClient) GetEvent(id string) (requestedEvent *models.Event, err error) {
	collection := c.Session.Collection(eventTable)

	condition := db.Cond{"id": id}

	results := collection.Find(condition)

	exists, err := results.Exists()
	if err != nil {
		return nil, fmt.Errorf("error fetching event: %w", err)
	}

	if !exists {
		return nil, core.StorageErrorNotFound
	}

	err = results.One(&requestedEvent)
	if err != nil {
		return nil, fmt.Errorf("error finding event: %w", err)
	}

	return requestedEvent, nil
}

func (c *upperClient) GetEvents() (result []*models.Event, err error) {
	collection := c.Session.Collection(eventTable)

	var events []models.Event

	err = collection.Find().All(&events)
	if err != nil {
		return nil, fmt.Errorf("error fetching all events: %w", err)
	}

	for _, event := range events {
		event := event

		result = append(result, &event)
	}

	return result, nil
}

func (c *upperClient) UpdateEvent(update *models.Event) (updateResult *models.Event, err error) {
	collection := c.Session.Collection(eventTable)

	var originalEvent models.Event

	condition := db.Cond{"id": update.Id}

	result := collection.Find(condition)

	exists, err := result.Exists()
	if err != nil {
		return nil, fmt.Errorf("error fetching event: %w", err)
	}

	if !exists {
		return nil, core.StorageErrorNotFound
	}

	err = result.One(&originalEvent)
	if err != nil {
		return nil, fmt.Errorf("error fetching original event: %w", err)
	}

	err = mergo.Merge(&originalEvent, *update, mergo.WithOverride)
	if err != nil {
		return nil, fmt.Errorf("error merging updated with original event: %w", err)
	}

	err = collection.UpdateReturning(&originalEvent)
	if err != nil {
		return nil, fmt.Errorf("error updating story: %w", err)
	}

	return &originalEvent, nil
}

func (c *upperClient) DeleteEvent(id string) (err error) {
	collection := c.Session.Collection(eventTable)

	condition := db.Cond{"id": id}

	result := collection.Find(condition)

	exists, err := result.Exists()
	if err != nil {
		return fmt.Errorf("error fetching event: %w", err)
	}

	if !exists {
		return core.StorageErrorNotFound
	}

	err = result.Delete()
	if err != nil {
		return fmt.Errorf("error deleting story with ID %s: %w", id, err)
	}

	return nil
}

func (c *upperClient) Open() error {
	sess, err := postgresql.Open(c.connectionURL)
	if err != nil {
		return fmt.Errorf("error connecting to Postgres: %w", err)
	}

	c.Session = sess

	return c.setup()
}

func (c *upperClient) Close() error {
	return c.Session.Close()
}

func (c *upperClient) setup() error {
	sql := c.Session.SQL()

	_, err := sql.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id text primary key,
		timestamp text not null,
		title text not null,
		description text not null,
		imageurl text,
		readmoreurl text
	)`, eventTable))
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	return nil
}

func NewUpperClient(uri, dbname, username, password string) core.StorageClient {
	return &upperClient{
		connectionURL: &postgresql.ConnectionURL{
			User:     username,
			Password: password,
			Host:     uri,
			Database: dbname,
		},
	}
}
