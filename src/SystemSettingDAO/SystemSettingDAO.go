package SystemSettingDAO

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"../AppEngineHelper"
)

const (
	SETTING_KIND             = "Setting"
	SETTING_PARENT_STRING_ID = "default_setting"
)

var (
	settingCollectionParentKey = AppEngineHelper.CollectionParentKeyGetFnGenerator(SETTING_KIND, SETTING_PARENT_STRING_ID, 0)
	settingIntIDToKeyInt64     = AppEngineHelper.IntIDToKeyInt64(SETTING_KIND, settingCollectionParentKey)
)

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

	if value == "" {
		if dsKey != nil {
			datastore.Delete(ctx, dsKey)
		}
		return nil
	}

	if value == currentValue {
		return nil
	}

	if dsKey == nil {
		dsKey = datastore.NewIncompleteKey(ctx, SETTING_KIND, settingCollectionParentKey(ctx))
	}

	_, err = datastore.Put(ctx, dsKey, &systemSetting{Key: key, Value: value})
	if err != nil {
		return err
	}

	return nil
}

// Returns eventually consistent results
func GetSetting(ctx context.Context, key string) (*datastore.Key, string, error) {
	settings := make([]systemSetting, 0, 1)
	keys, err := datastore.NewQuery(SETTING_KIND).
		Ancestor(settingCollectionParentKey(ctx)).
		Filter("Key = ", key).
		Limit(1).
		GetAll(ctx, &settings)

	if err != nil {
		return nil, "", err
	}

	if len(keys) > 0 {
		return keys[0], settings[0].Value, nil
	}

	return nil, "", nil
}
