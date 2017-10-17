package Auth

import (
	"google.golang.org/appengine/datastore"
)

type SessionData struct {
	Uuid    string
	UserKey *datastore.Key
}

func (sd *SessionData) HasLoginInfo() bool {
	return sd.UserKey != nil && sd.Uuid != ""
}
