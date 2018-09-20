package AccessIdValidator

import (
	"bytes"
	"io/ioutil"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"

	"Dropbox"
)

type DefaultAccessIdValidator struct{}

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

func downloadValidAccessIds(ctx context.Context, dropboxAccessToken string) ([][]byte, error) {
	log.Infof(ctx, "Downloading AccessIds: %s", dropboxAccessToken)
	resp, _, err := Dropbox.DownloadFile(ctx, dropboxAccessToken, GetAccessIdPath(ctx))
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

	idsNoSpace := make([][]byte, 0, len(ids))

	for _, id := range ids {
		trimmedId := bytes.TrimSpace(id)
		if len(trimmedId) > 0 {
			idsNoSpace = append(idsNoSpace, trimmedId)
		}
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

	lastDownload = time.Now()

	return nil
}

func updateTokenCache(ctx context.Context, settingKey string, currentCache [][]byte) ([][]byte, error) {
	token, err := Dropbox.GetAccessToken(ctx, settingKey)
	if err != nil {
		return nil, err
	}
	if token == "" {
		log.Warningf(ctx, "No valid Dropbox access token found, please configure dropbox integration")
		return nil, nil
	}

	if len(currentCache) > 0 && lastDownload.Add(4*time.Hour).After(time.Now()) {
		return currentCache, nil
	}

	log.Infof(ctx, "Cache expired downloading accessId list with settingKey: %s", settingKey)

	ids, err := downloadValidAccessIds(ctx, token)
	if err != nil {
		log.Warningf(ctx, "Unable to download valid accessIds, old Ids will be used, error: %s", err.Error())
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

	log.Infof(ctx, "accessId not valid - str length: %v, str: '%q', hex: %X", len(accessId), accessId, accessId)

	return false, nil
}

func (v *DefaultAccessIdValidator) ValidateAccessIdPrimary(ctx context.Context, accessId []byte) (bool, error) {
	return validateAccessId(ctx, accessId, &primaryIds)
}

func (v *DefaultAccessIdValidator) ValidateAccessIdSecondary(ctx context.Context, accessId []byte) (bool, error) {
	return validateAccessId(ctx, accessId, &secondaryIds)
}
