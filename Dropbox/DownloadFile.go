package Dropbox

import (
	"io"

	"appengine"
	"appengine/urlfetch"
)

func DownloadFile(ctx appengine.Context, url string) (io.ReadCloser, error) {
	client := urlfetch.Client(ctx)

	rsp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}
