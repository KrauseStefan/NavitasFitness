package Dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"AppEngineHelper"
	"SystemSettingDAO"
	"appengine"
	"appengine/urlfetch"
)

const (
	baseUrl  = "https://content.dropboxapi.com"
	tokenUrl = "https://api.dropboxapi.com/oauth2/token"

	accessTokenSystemSettingKey = "accessToken"
)

var (
	clientKey    = ""
	clientSecret = ""

	accessToken = ""
)

type TokenRspDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AccountId   string `json:"account_id"`
	Uid         string `json:"uid"`
}

func GetAccessToken(ctx appengine.Context) (string, error) {
	var (
		err              error = nil
		accessTokenValue string
	)

	if accessToken == "" {
		_, accessTokenValue, err = SystemSettingDAO.GetSetting(ctx, accessTokenSystemSettingKey)
		accessToken = accessTokenValue
	}

	return accessToken, err
}

func RetrieveAccessToken(ctx appengine.Context, code string, redirectUri string) error {

	params := map[string]string{
		"code":         code,                 // String The code acquired by directing users to /oauth2/authorize?response_type=code.
		"grant_type":   "authorization_code", //The grant type, which must be authorization_code.
		"redirect_uri": redirectUri,
	}

	paramStr := AppEngineHelper.CreateQueryParamString(params)

	client := urlfetch.Client(ctx)

	req, err := http.NewRequest("POST", tokenUrl+"?"+paramStr, &bytes.Buffer{})
	if err != nil {
		return err
	}

	req.SetBasicAuth(clientKey, clientSecret)

	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(all))
	}

	tokenRspDTO := TokenRspDTO{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tokenRspDTO)
	if err != nil {
		return err
	}

	accessToken = tokenRspDTO.AccessToken

	return SystemSettingDAO.PersistSetting(ctx, "accessToken", accessToken)
}
