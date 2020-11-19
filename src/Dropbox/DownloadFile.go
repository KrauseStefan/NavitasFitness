package Dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"DAOHelper"
)

const serviceUrl = "https://content.dropboxapi.com/2/files/download"

type FileDownloadRequestDto struct {
	Path string `json:"path"`
}

type FileDownloadResponseDto struct {
	Name                     string
	Id                       string
	ClientModified           string
	ServerModified           string // "2015-05-12T15:50:38Z"
	Rev                      string
	Size                     int
	PathLower                string
	PathDisplay              string
	SharingInfo              SharingInfoDto
	PropertyGroups           []PropertyGroupDto
	HasExplicitSharedMembers bool
	ContentHash              string
}

type SharingInfoDto struct {
	ReadOnly             bool
	ParentSharedFolderId string
	ModifiedBy           string
}

type PropertyGroupDto struct {
	TemplateId string
	Fields     []FieldDto
}

type FieldDto struct {
	Name  string
	Value string
}

func DownloadFile(ctx context.Context, accessToken string, fileUrl string) (io.ReadCloser, *FileDownloadResponseDto, error) {
	if accessToken == "" || fileUrl == "" {
		return nil, nil, errors.New("Invalid arguments for downloading file")
	}

	req, err := http.NewRequest("POST", serviceUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	js, err := json.Marshal(&FileDownloadRequestDto{Path: fileUrl})
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Dropbox-API-Arg", string(js))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		jsonError, _ := ioutil.ReadAll(resp.Body)

		errorDtoMap := make(map[string]interface{})
		if err := json.Unmarshal([]byte(jsonError), &errorDtoMap); err != nil {
			return nil, nil, err
		}

		err := &DAOHelper.DefaultHttpError{StatusCode: resp.StatusCode}
		if val, ok := errorDtoMap["error_summary"].(string); ok {
			err.InnerError = errors.New(string(val))
		} else {
			err.InnerError = errors.New(string(jsonError))
		}

		return nil, nil, err
	}

	respJson := resp.Header.Get("dropbox-api-result")
	fileDownloadResponseDto := FileDownloadResponseDto{}
	decoder := json.NewDecoder(bytes.NewBufferString(respJson))
	err = decoder.Decode(&fileDownloadResponseDto)
	if err != nil {
		return nil, nil, err
	}

	return resp.Body, &fileDownloadResponseDto, nil
}
