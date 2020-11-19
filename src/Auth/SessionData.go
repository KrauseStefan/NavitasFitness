package Auth

import (
	"cloud.google.com/go/datastore"
)

type SessionData struct {
	Uuid    string
	UserKey *datastore.Key
}

func (sd *SessionData) HasLoginInfo() bool {
	return sd.UserKey != nil && sd.Uuid != ""
}
