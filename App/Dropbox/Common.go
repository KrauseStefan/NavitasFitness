package Dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"AppEngineHelper"
	"ConfigurationReader"
	"SystemSettingDAO"
	"appengine"
	"appengine/urlfetch"
)

var (
	accessToken_intenal = ""
)

const (
	baseUrl  = "https://content.dropboxapi.com"
	tokenUrl = "https://api.dropboxapi.com/oauth2/token"

	accessTokenSystemSettingKey = "accessToken"
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

	if accessToken_intenal == "" {
		_, accessTokenValue, err = SystemSettingDAO.GetSetting(ctx, accessTokenSystemSettingKey)
		accessToken_intenal = accessTokenValue
	}

	return accessToken_intenal, err
}

func RetrieveAccessToken(ctx appengine.Context, code string, redirectUri string) error {
	conf, err := ConfigurationReader.GetConfiguration()
	if err != nil {
		return err
	}

	params := map[string]string{
		"code":         code,
		"grant_type":   "authorization_code", //The grant type, which must be authorization_code.
		"redirect_uri": redirectUri,
	}

	paramStr := AppEngineHelper.CreateQueryParamString(params)

	client := urlfetch.Client(ctx)

	req, err := http.NewRequest("POST", tokenUrl+"?"+paramStr, &bytes.Buffer{})
	if err != nil {
		return err
	}

	req.SetBasicAuth(conf.ClientKey, conf.ClientSecret)

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

	accessToken_intenal = tokenRspDTO.AccessToken

	return SystemSettingDAO.PersistSetting(ctx, "accessToken", accessToken_intenal)
}
