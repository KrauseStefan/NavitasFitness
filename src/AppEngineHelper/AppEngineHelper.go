package AppEngineHelper

import (
	"DAOHelper"
	"encoding/json"
	"google.golang.org/appengine"
	"net/http"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type FormDataDecoderFn func(interface{}) error

type httpHandlerWithData func(http.ResponseWriter, *http.Request, FormDataDecoderFn)

type HttpHandler func(http.ResponseWriter, *http.Request) (interface{}, error)

func HandlerW(f HttpHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)

		dto, err := f(w, r)
		if err == nil && dto != nil {
			if dtoStr, ok := dto.(string); ok {
				w.Write([]byte(dtoStr))
			} else if dtoBytes, ok := dto.([]byte); ok {
				w.Write(dtoBytes)
			} else {
				_, err = WriteJSON(w, dto)
			}
		}

		if err != nil {
			DAOHelper.ReportError(ctx, w, err)
		}
	}
}

func WriteJSON(w http.ResponseWriter, data interface{}) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	if err == nil {
		w.Write(js)
	}

	return js, err
}

type CollectionParentKeyGetter func(c context.Context) *datastore.Key
type Int64KeyGetter func(c context.Context, id string) *datastore.Key

func CollectionParentKeyGetFnGenerator(kind string, parentStringId string, keyId int64) CollectionParentKeyGetter {
	return func(ctx context.Context) *datastore.Key {
		return datastore.NewKey(ctx, kind, parentStringId, 0, nil)
	}
}

func IntIDToKeyInt64(kind string, parentKeyGetter CollectionParentKeyGetter) Int64KeyGetter {
	return func(ctx context.Context, id string) *datastore.Key {
		intId, _ := strconv.ParseInt(id, 10, 64)
		return datastore.NewKey(ctx, kind, "", intId, parentKeyGetter(ctx))
	}
}

func CreateQueryParamString(params map[string]string) string {
	var paramStr string
	for k, v := range params {
		if len(paramStr) != 0 {
			paramStr = paramStr + "&"
		}
		paramStr = paramStr + k + "=" + v
	}
	return paramStr
}

func StringIdsToDsKeys(ids []string) ([]*datastore.Key, error) {
	idKeys := make([]*datastore.Key, len(ids))

	for i, id := range ids {
		key, err := datastore.DecodeKey(id)
		if err != nil {
			return nil, err
		}
		idKeys[i] = key
	}
	return idKeys, nil
}
