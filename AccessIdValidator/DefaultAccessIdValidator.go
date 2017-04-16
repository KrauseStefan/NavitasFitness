package AccessIdValidator

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"time"

	"appengine"

	"Dropbox"
	"SystemSettingDAO"
)

type DefaultAccessIdValidator struct{}

const (
	fitnessAccessIdsPathSettingKey = "fitnessAccessIdsPath"
	defaultFitnessAccessIdsPath    = "/AccessIds/AccessIds.csv"
)

var (
	bomPrefix               = []byte{0xef, 0xbb, 0xbf}
	downloadedIds *[][]byte = nil
	lastDownload  time.Time
)

var instance = DefaultAccessIdValidator{}

func GetInstance() AccessIdValidator {
	return &instance
}

func getPath(ctx appengine.Context) string {
	_, value, err := SystemSettingDAO.GetSetting(ctx, fitnessAccessIdsPathSettingKey)
	if err != nil {
		ctx.Errorf(err.Error())
	}

	if value == "" {
		value = defaultFitnessAccessIdsPath
		if err := SystemSettingDAO.PersistSetting(ctx, fitnessAccessIdsPathSettingKey, value); err != nil {
			ctx.Errorf(err.Error())
		}
	}

	return value
}

func downloadValidAccessIds(ctx appengine.Context) error {

	resp, _, err := Dropbox.DownloadFile(ctx, getPath(ctx))
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	//BOM does not make sense for UTF-8, should be safe to strip
	dataWithoutBom := bytes.TrimPrefix(data, bomPrefix)

	ids := bytes.Split(dataWithoutBom, []byte{'\n'})

	idsNoSpace := make([][]byte, len(ids), len(ids))

	for i, id := range ids {
		idsNoSpace[i] = bytes.TrimSpace(id)
	}

	downloadedIds = &idsNoSpace
	lastDownload = time.Now()
	return nil
}

func ensureUpdatedIds(ctx appengine.Context) error {
	if downloadedIds == nil || lastDownload.Add(12*time.Hour).Before(time.Now()) {
		ctx.Infof("Downloading ids")
		err := downloadValidAccessIds(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

//Go strings are UTF 8 without bom converting it to byte should be safe
func (v *DefaultAccessIdValidator) ValidateAccessId(ctx appengine.Context, accessId []byte) (bool, error) {

	if err := ensureUpdatedIds(ctx); err != nil {
		return false, err
	}

	for _, id := range *downloadedIds {
		if bytes.Equal(id, accessId) {
			return true, nil
		}
	}

	ctx.Infof("length %s - %d", hex.Dump(accessId), len(accessId))

	return false, nil
}
