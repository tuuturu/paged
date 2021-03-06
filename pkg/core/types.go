package core

import (
	"errors"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/tuuturu/pager-event-service/pkg/core/models"

	"github.com/gin-gonic/gin"
)

type DSN struct {
	Scheme       string
	URI          string
	Port         string
	DatabaseName string
	Username     string
	Password     string
}

type Config struct {
	DiscoveryURL *url.URL
	Port         string

	Database *DSN

	LogLevel log.Level

	ClientID     string
	ClientSecret string
}

type GetEventsFilter struct {
	Read *bool
}

// StorageClient defines the interface a storage client should expose
type StorageClient interface {
	// Open initiates the connection to the storage backend
	Open() error
	Close() error

	AddEvent(event *models.Event) error
	GetEvent(id string) (*models.Event, error)
	GetEvents(filter GetEventsFilter) ([]*models.Event, error)
	UpdateEvent(event *models.Event) (*models.Event, error)
	DeleteEvent(id string) error
}

var StorageErrorNotFound = errors.New("not found")

type HandlerFuncGenerator func(StorageClient) gin.HandlerFunc

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc HandlerFuncGenerator
}

// Routes is the list of the generated Route.
type Routes []Route
