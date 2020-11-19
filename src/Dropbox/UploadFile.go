package Dropbox

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
)

const uploadUrl = "/2/files/upload"

type UploadResp struct {
	Name           string `json:"name"`
	PathLower      string `json:"path_lower"`
	PathDisplay    string `json:"path_display"`
	Id             string `json:"id"`
	ClientModified string `json:"client_modified"`
	ServerModified string `json:"server_modified"`
	Rev            string `json:"rev"`
	Size           int    `json:"size"`
	ContentHash    string `json:"content_hash"`
}

type UploadModeDto struct {
	Tag string `json:".tag"`
}

type UploadDto struct {
	Path string        `json:"path"`
	Mode UploadModeDto `json:"mode"`
}

func UploadDoc(ctx context.Context, accessToken string, filePath string, file io.Reader) (*UploadResp, error) {

	req, err := http.NewRequest("POST", baseUrl+uploadUrl, file)
	if err != nil {
		return nil, err
	}

	uploadDto := UploadDto{
		Path: filePath,
		Mode: UploadModeDto{
			Tag: "overwrite",
		},
	}

	js, err := json.Marshal(&uploadDto)
	if err != nil {
		return nil, err
	}

	if len(accessToken) <= 0 {
		return nil, errors.New("Access id has not been assigned, unable to upload to dropbox")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", string(js))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(all))
	}

	uploadResp := UploadResp{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&uploadResp)
	if err != nil {
		return nil, err
	}

	return &uploadResp, nil
}
