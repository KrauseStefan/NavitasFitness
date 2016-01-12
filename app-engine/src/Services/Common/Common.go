package Common
import (
	"net/http"
	"encoding/json"
	"appengine/datastore"
	"appengine"
	"strconv"
)

func WriteJSON(w http.ResponseWriter, data interface{}) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	w.Write(js)

	return js, err
}

type CollectionParentKeyGetter func(c appengine.Context) *datastore.Key
type Int64KeyGetter func(c appengine.Context, id string) *datastore.Key

func CollectionParentKeyGetFnGenerator(kind string, parentStringId string, keyId int64) CollectionParentKeyGetter {
	return func (c appengine.Context) *datastore.Key {
		return datastore.NewKey(c, kind, parentStringId, 0, nil)
	}
}

func IntIDToKeyInt64(kind string, parentKeyGetter CollectionParentKeyGetter) Int64KeyGetter {
	return func (c appengine.Context, id string) *datastore.Key {
		intId, _ := strconv.ParseInt(id, 10, 64)
		return datastore.NewKey(c, kind, "", intId, parentKeyGetter(c))
	}
}
