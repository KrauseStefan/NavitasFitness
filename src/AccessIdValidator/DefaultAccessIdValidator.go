package AccessIdValidator

import (
	"bytes"
	"io/ioutil"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"

	"Dropbox"
	"SystemSettingDAO"
)

type DefaultAccessIdValidator struct{}

const (
	fitnessAccessIdsPathSettingKey = "fitnessAccessIdsPath"
	defaultFitnessAccessIdsPath    = "/AccessIds/AccessIds.csv"
)

var (
	bomPrefix             = []byte{0xef, 0xbb, 0xbf}
	primaryIds   [][]byte = nil
	secondaryIds [][]byte = nil
	lastDownload time.Time
)

var instance = DefaultAccessIdValidator{}

func GetInstance() AccessIdValidator {
	return &instance
}

func getPath(ctx context.Context) string {
	_, value, err := SystemSettingDAO.GetSetting(ctx, fitnessAccessIdsPathSettingKey)
	if err != nil {
		log.Errorf(ctx, err.Error())
	}

	if value == "" {
		value = defaultFitnessAccessIdsPath
		if err := SystemSettingDAO.PersistSetting(ctx, fitnessAccessIdsPathSettingKey, value); err != nil {
			log.Errorf(ctx, err.Error())
		}
	}

	return value
}

func downloadValidAccessIds(ctx context.Context, dropboxAccessToken string) ([][]byte, error) {
	log.Infof(ctx, "Downloading AccessIds: %s", dropboxAccessToken)
	resp, _, err := Dropbox.DownloadFile(ctx, dropboxAccessToken, getPath(ctx))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	//BOM does not make sense for UTF-8, should be safe to strip
	dataWithoutBom := bytes.TrimPrefix(data, bomPrefix)

	ids := bytes.Split(dataWithoutBom, []byte{'\n'})

	idsNoSpace := make([][]byte, len(ids), len(ids))

	for i, id := range ids {
		idsNoSpace[i] = bytes.TrimSpace(id)
	}

	return idsNoSpace, nil
}

func ensureUpdatedIds(ctx context.Context) error {
	var err error
	primaryIds, err = updateTokenCache(ctx, Dropbox.PrimaryAccessTokenSystemSettingKey, primaryIds)
	if err != nil {
		return err
	}

	secondaryIds, err = updateTokenCache(ctx, Dropbox.SecondaryAccessTokenSystemSettingKey, secondaryIds)
	if err != nil {
		return err
	}

	log.Infof(ctx, "length primary: %d, secondary: %d", len(primaryIds), len(secondaryIds))
	lastDownload = time.Now()

	return nil
}

func updateTokenCache(ctx context.Context, settingKey string, currentCache [][]byte) ([][]byte, error) {
	token, err := Dropbox.GetAccessToken(ctx, settingKey)
	if token == "" {
		return nil, nil
	}

	if err != nil || !(len(currentCache) <= 0 || lastDownload.Add(4*time.Hour).Before(time.Now())) {
		return currentCache, err
	}

	ids, err := downloadValidAccessIds(ctx, token)
	if err != nil {
		return currentCache, err
	}

	return ids, nil
}

func validateAccessId(ctx context.Context, accessId []byte, validIdList *[][]byte) (bool, error) {
	if err := ensureUpdatedIds(ctx); err != nil {
		return false, err
	}

	if len(*validIdList) <= 0 {
		log.Warningf(ctx, "No valid access Ids, list is empty!")
		return false, nil
	}

	for _, validId := range *validIdList {
		if bytes.Equal(validId, accessId) {
			return true, nil
		}
	}

	log.Infof(ctx, "Id not validated")
	log.Infof(ctx, "length %v - hex: %X", len(accessId), accessId)

	return false, nil
}

func (v *DefaultAccessIdValidator) ValidateAccessIdPrimary(ctx context.Context, accessId []byte) (bool, error) {
	log.Infof(ctx, "primaryIds: %d", len(primaryIds))
	return validateAccessId(ctx, accessId, &primaryIds)
}

func (v *DefaultAccessIdValidator) ValidateAccessIdSecondary(ctx context.Context, accessId []byte) (bool, error) {
	return validateAccessId(ctx, accessId, &secondaryIds)
}
