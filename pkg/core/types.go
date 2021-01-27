package core

import (
	"errors"
	"net/url"

	"github.com/tuuturu/paged/pkg/core/models"

	"github.com/gin-gonic/gin"
)

type DatabaseOptions struct {
	URI          string
	DatabaseName string
	Username     string
	Password     string
}

type Config struct {
	DiscoveryURL url.URL
	Port         string

	Database *DatabaseOptions
}

type StorageClient interface {
	Open() error
	Close() error

	AddEvent(event *models.Event) error
	GetEvent(id string) (*models.Event, error)
	GetEvents() ([]*models.Event, error)
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
