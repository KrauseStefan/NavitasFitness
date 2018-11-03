package AccessIdOverrideDao

import (
	"AppEngineHelper"
	"context"
	"errors"
	"google.golang.org/appengine/datastore"
	"time"
)

type AccessIdOverride struct {
	AccessId  string         `json:"accessId"`
	StartDate time.Time      `json:"startDate" datastore:",noindex"`
	Key       *datastore.Key `json:"-" datastore:"-"`
}

const (
	ACCESS_ID_OVERRIDE_KIND             = "accessIdOverride"
	ACCESS_ID_OVERRIDE_PARENT_STRING_ID = "default_accessIdOverride"
)

var (
	accessIdOverrideCollectionParentKey = AppEngineHelper.CollectionParentKeyGetFnGenerator(ACCESS_ID_OVERRIDE_KIND, ACCESS_ID_OVERRIDE_PARENT_STRING_ID, 0)
	accessIdOverrideIntIDToKeyInt64     = AppEngineHelper.IntIDToKeyInt64(ACCESS_ID_OVERRIDE_KIND, accessIdOverrideCollectionParentKey)

	AccessIdNotFoundError = errors.New("AccessId does not exist in datastore")
)

func GetAllAccessIdOverrides(ctx context.Context) ([]*AccessIdOverride, error) {
	accessIdList := make([]*AccessIdOverride, 0, 1)
	keys, err := datastore.NewQuery(ACCESS_ID_OVERRIDE_KIND).
		Ancestor(accessIdOverrideCollectionParentKey(ctx)).
		GetAll(ctx, &accessIdList)

	keysLen := len(keys)
	for i, id := range accessIdList {
		if keysLen > i && id != nil {
			id.Key = keys[i]
		}
	}

	return accessIdList, err
}

func CreateOrUpdateAccessIdOverride(ctx context.Context, accessIdOverride *AccessIdOverride) error {
	key := datastore.NewIncompleteKey(ctx, ACCESS_ID_OVERRIDE_KIND, accessIdOverrideCollectionParentKey(ctx))
	newKey, err := datastore.Put(ctx, key, accessIdOverride)

	accessIdOverride.Key = newKey
	return err
}

func DeleteAccessIdOverride(ctx context.Context, accessId string) error {
	accessIdList := make([]AccessIdOverride, 1)
	keys, err := datastore.NewQuery(ACCESS_ID_OVERRIDE_KIND).
		Ancestor(accessIdOverrideCollectionParentKey(ctx)).
		Filter("AccessId=", accessId).
		Limit(1).
		GetAll(ctx, &accessIdList)

	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return AccessIdNotFoundError
	}

	return datastore.Delete(ctx, keys[0])
}
