package upper

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

const eventTable = "events"

type upperClient struct {
	connectionURL *postgresql.ConnectionURL
	Session       db.Session
}
