package AccessIdValidator

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"time"

	"appengine"

	"Dropbox"
)

const (
	downloadLink = "https://www.dropbox.com/s/x61e0hter4kr9ry/Kort%20navitas.csv?_download_id=462948430621974902922706508949901694802577845535546831947836167394&dl=1"
)

var (
	bomPrefix               = []byte{0xef, 0xbb, 0xbf}
	downloadedIds *[][]byte = nil
	lastDownload  time.Time
)

func downloadValidAccessIds(ctx appengine.Context) error {

	resp, err := Dropbox.DownloadFile(ctx, downloadLink)
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
func ValidateAccessId(ctx appengine.Context, accessId []byte) (bool, error) {

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
