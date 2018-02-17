package AccessIdValidator

import (
	"DAOHelper"
	"Dropbox"
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"net/http"
)

func PushMissingSampleData(ctx context.Context, accessToken string) {
	accessIdPath := GetAccessIdPath(ctx)
	if err := uploadFileIfMissing(ctx, accessToken, accessIdPath); err != nil {
		log.Errorf(ctx, err.Error())
	}

	accessListPath := GetAccessListPath(ctx)
	if err := uploadFileIfMissing(ctx, accessToken, accessListPath); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func uploadFileIfMissing(ctx context.Context, accessToken string, path string) error {
	exists, err := CheckFileExistence(ctx, accessToken, path)
	if err != nil || exists {
		return err
	}
	log.Infof(ctx, "%s is missing uploading empty file", path)

	_, err = Dropbox.UploadDoc(ctx, accessToken, path, &bytes.Buffer{})
	return err
}

func CheckFileExistence(ctx context.Context, accessToken string, path string) (bool, error) {
	_, _, err := Dropbox.DownloadFile(ctx, accessToken, path)
	if err != nil {
		if httpErr, ok := err.(*DAOHelper.DefaultHttpError); ok && httpErr.StatusCode == http.StatusConflict {
			log.Debugf(ctx, "Checking for file error %s - (%d) %s", path, httpErr.StatusCode, httpErr.Error())
			return false, nil
		}

		return false, err
	}

	log.Infof(ctx, "File exists in dropbox: %s", path)
	return true, nil
}
