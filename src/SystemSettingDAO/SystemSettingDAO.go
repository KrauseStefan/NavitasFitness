package SystemSettingDAO

import (
	"golang.org/x/net/context"

	"cloud.google.com/go/datastore"

	nf_datastore "github.com/KrauseStefan/NavitasFitness/NavitasFitness/datastore"
	log "github.com/KrauseStefan/NavitasFitness/logger"
)

const (
	SETTING_KIND             = "Setting"
	SETTING_PARENT_STRING_ID = "default_setting"
)

var settingCollectionParentKey = datastore.NameKey(SETTING_KIND, SETTING_PARENT_STRING_ID, nil)

type systemSetting struct {
	Key   string
	Value string
}

// Setting cannot be retrieved in a consistent way. Be aware
func PersistSetting(ctx context.Context, key string, value string) error {
	dsKey, currentValue, err := GetSetting(ctx, key)
	if err != nil {
		return err
	}

	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return err
	}

	if value == "" {
		if dsKey != nil {
			dsClient.Delete(ctx, dsKey)
		}
		return nil
	}

	if value == currentValue {
		return nil
	}

	if dsKey == nil {
		dsKey = datastore.IncompleteKey(SETTING_KIND, settingCollectionParentKey)
	}

	_, err = dsClient.Put(ctx, dsKey, &systemSetting{Key: key, Value: value})
	if err != nil {
		return err
	}

	return nil
}

// Returns eventually consistent results
func GetSetting(ctx context.Context, key string) (*datastore.Key, string, error) {
	log.Debugf(ctx, "GetSetting with key: %s", key)

	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, "", err
	}

	settings := make([]systemSetting, 0, 1)
	query := datastore.NewQuery(SETTING_KIND).
		Ancestor(settingCollectionParentKey).
		Filter("Key = ", key).
		Limit(1)

	keys, err := dsClient.GetAll(ctx, query, &settings)
	if err != nil {
		return nil, "", err
	}

	if len(keys) > 0 {
		return keys[0], settings[0].Value, nil
	}

	return nil, "", nil
}
