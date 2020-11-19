package AccessIdOverrideDao

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/datastore"

	nf_datastore "NavitasFitness/datastore"
)

type DefaultAccessIdOverrideDao struct{}

var defaultAccessIdOverrideDao = DefaultAccessIdOverrideDao{}

func GetInstance() *DefaultAccessIdOverrideDao {
	return &defaultAccessIdOverrideDao
}

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
	accessIdOverrideCollectionParentKey = datastore.NameKey(ACCESS_ID_OVERRIDE_KIND, ACCESS_ID_OVERRIDE_PARENT_STRING_ID, nil)
	AccessIdNotFoundError               = errors.New("AccessId does not exist in datastore")
)

func (dao *DefaultAccessIdOverrideDao) GetAllAccessIdOverrides(ctx context.Context) ([]*AccessIdOverride, error) {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	accessIdList := make([]*AccessIdOverride, 0, 1)
	query := datastore.NewQuery(ACCESS_ID_OVERRIDE_KIND).
		Ancestor(accessIdOverrideCollectionParentKey)
	keys, err := dsClient.GetAll(ctx, query, &accessIdList)

	keysLen := len(keys)
	for i, id := range accessIdList {
		if keysLen > i && id != nil {
			id.Key = keys[i]
		}
	}

	return accessIdList, err
}

func (dao *DefaultAccessIdOverrideDao) CreateOrUpdateAccessIdOverride(ctx context.Context, accessIdOverride *AccessIdOverride) error {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return err
	}

	key := datastore.IncompleteKey(ACCESS_ID_OVERRIDE_KIND, accessIdOverrideCollectionParentKey)
	newKey, err := dsClient.Put(ctx, key, accessIdOverride)

	accessIdOverride.Key = newKey
	return err
}

func (dao *DefaultAccessIdOverrideDao) DeleteAccessIdOverride(ctx context.Context, accessId string) error {
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return err
	}

	accessIdList := make([]AccessIdOverride, 1)
	query := datastore.NewQuery(ACCESS_ID_OVERRIDE_KIND).
		Ancestor(accessIdOverrideCollectionParentKey).
		Filter("AccessId=", accessId).
		Limit(1)
	keys, err := dsClient.GetAll(ctx, query, &accessIdList)

	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return AccessIdNotFoundError
	}

	return dsClient.Delete(ctx, keys[0])
}
