package Auth

import (
	"google.golang.org/appengine/datastore"
)

type SessionData struct {
	Uuid    string
	UserKey *datastore.Key
}
