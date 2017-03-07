package Dropbox

import (
	"appengine"
	"appengine/urlfetch"
	"bytes"
	"io/ioutil"
	"time"
)

const downloadLink = "https://www.dropbox.com/s/x61e0hter4kr9ry/Kort%20navitas.csv?_download_id=462948430621974902922706508949901694802577845535546831947836167394&dl=1"

var downloadedIds *[][]byte = nil

var lastDownload time.Time

func downloadValidAccessIds(ctx appengine.Context) error {

	client := urlfetch.Client(ctx)

	rsp, err := client.Get(downloadLink)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	ids := bytes.Split(data, []byte{'\n'})

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

func ValidateAccessId(ctx appengine.Context, accessId []byte) (bool, error) {

	if err := ensureUpdatedIds(ctx); err != nil {
		return false, err
	}

	for _, id := range *downloadedIds {
		if bytes.Equal(id, accessId) {
			return true, nil
		}
	}

	ctx.Infof("%s", *downloadedIds)

	return false, nil
}
