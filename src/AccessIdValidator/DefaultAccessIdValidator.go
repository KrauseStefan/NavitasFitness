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
	lastDownload time.Time
)

var instance = DefaultAccessIdValidator{}

func GetInstance() AccessIdValidator {
	return &instance
}

func (v *DefaultAccessIdValidator) downloadValidAccessIds(ctx context.Context, dropboxAccessToken string) ([][]byte, error) {
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

func (v *DefaultAccessIdValidator) EnsureUpdatedIds(ctx context.Context) error {
	var err error
	primaryIds, err = v.updateTokenCache(ctx, Dropbox.PrimaryAccessTokenSystemSettingKey, primaryIds)
	if err != nil {
		return err
	}

	if len(primaryIds) <= 0 {
		log.Warningf(ctx, "No valid access Ids, list is empty!")
	}

	return nil
}

func (v *DefaultAccessIdValidator) updateTokenCache(ctx context.Context, settingKey string, currentCache [][]byte) ([][]byte, error) {
	cacheExpirationTime := lastDownload.Add(4 * time.Hour)
	if len(currentCache) > 0 && cacheExpirationTime.Before(time.Now()) {
		return currentCache, nil
	}
	if len(currentCache) > 0 {
		log.Infof(ctx, "Cache is empty")
	} else {
		cacheExpirationTimeStr := cacheExpirationTime.Format("02-01-06 15:04:05")
		nowStr := time.Now().Format("02-01-06 15:04:05")
		log.Infof(ctx, "Cache expired cacheExpirationTime is %s, current time is: %s", cacheExpirationTimeStr, nowStr)
		lastDownload = time.Now()
	}

	token, err := Dropbox.GetAccessToken(ctx, settingKey)
	if err != nil {
		return nil, err
	}
	if token == "" {
		log.Warningf(ctx, "No valid Dropbox access token found, please configure dropbox integration")
		return nil, nil
	}

	ids, err := v.downloadValidAccessIds(ctx, token)
	if err != nil {
		log.Warningf(ctx, "Unable to download valid accessIds, old Ids will be used, error: %s", err.Error())
		return currentCache, err
	}

	dateStr := lastDownload.Format("02-01-06 15:04:05")
	log.Infof(ctx, "lastDownload: %s", dateStr)

	return ids, nil
}

func (v *DefaultAccessIdValidator) ValidateAccessId(ctx context.Context, accessId []byte) (bool, error) {
	for _, validId := range primaryIds {
		if bytes.Equal(validId, accessId) {
			return true, nil
		}
	}

	log.Infof(ctx, "accessId not valid - str length: %v, str: '%q', hex: %X", len(accessId), accessId, accessId)

	return false, nil
}
