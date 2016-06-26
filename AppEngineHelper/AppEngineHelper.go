package AppEngineHelper

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"

	"appengine"
	"appengine/datastore"
)

var formDataDecoder = schema.NewDecoder()

type FormDataDecoderFn func(interface{}) error

type httpHandlerWithData func(http.ResponseWriter, *http.Request, FormDataDecoderFn)

func ParseFormDataWrap(handler httpHandlerWithData) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		parseFormData := func(dst interface{}) error {
			err := r.ParseForm()

			if err != nil {
				return err
			}

			err = formDataDecoder.Decode(dst, r.PostForm)

			if err != nil {
				return err
			}
			return nil
		}

		handler(w, r, parseFormData)
	}
}

func WriteJSON(w http.ResponseWriter, data interface{}) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	w.Write(js)

	return js, err
}

type CollectionParentKeyGetter func(c appengine.Context) *datastore.Key
type Int64KeyGetter func(c appengine.Context, id string) *datastore.Key

func CollectionParentKeyGetFnGenerator(kind string, parentStringId string, keyId int64) CollectionParentKeyGetter {
	return func(ctx appengine.Context) *datastore.Key {
		return datastore.NewKey(ctx, kind, parentStringId, 0, nil)
	}
}

func IntIDToKeyInt64(kind string, parentKeyGetter CollectionParentKeyGetter) Int64KeyGetter {
	return func(ctx appengine.Context, id string) *datastore.Key {
		intId, _ := strconv.ParseInt(id, 10, 64)
		return datastore.NewKey(ctx, kind, "", intId, parentKeyGetter(ctx))
	}
}
