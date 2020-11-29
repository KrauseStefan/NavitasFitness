package Dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"AppEngineHelper"
	"ConfigurationReader"
	"SystemSettingDAO"
)

var (
	accessTokenIntenal = ""
)

const (
	baseUrl  = "https://content.dropboxapi.com"
	tokenUrl = "https://api.dropboxapi.com/oauth2/token"

	PrimaryAccessTokenSystemSettingKey = "PrimaryAccessToken"
)

type TokenRspDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AccountId   string `json:"account_id"`
	Uid         string `json:"uid"`
}

func GetAccessToken(ctx context.Context, key string) (string, error) {
	_, accessToken, err := SystemSettingDAO.GetSetting(ctx, key)
	if err != nil {
		return "", err
	}
	return accessToken, err
}

func appendAccessToken(ctx context.Context, key string, tokens []string) ([]string, error) {
	accessToken, err := GetAccessToken(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(accessToken) > 0 {
		tokens = append(tokens, accessToken)
	}
	return tokens, nil
}

func GetAccessTokens(ctx context.Context) ([]string, error) {
	var (
		tokens       = []string{}
		err    error = nil
	)

	if tokens, err = appendAccessToken(ctx, PrimaryAccessTokenSystemSettingKey, tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func RetrieveAccessToken(ctx context.Context, code string, redirectUri string) (string, error) {
	conf, err := ConfigurationReader.GetConfiguration()
	if err != nil {
		return "", err
	}

	params := map[string]string{
		"code":         code,
		"grant_type":   "authorization_code", //The grant type, which must be authorization_code.
		"redirect_uri": redirectUri,
	}

	paramStr := AppEngineHelper.CreateQueryParamString(params)

	req, err := http.NewRequest("POST", tokenUrl+"?"+paramStr, &bytes.Buffer{})
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(conf.ClientKey, conf.ClientSecret)

	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(all))
	}

	tokenRspDTO := TokenRspDTO{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tokenRspDTO)
	if err != nil {
		return "", err
	}

	accessTokenIntenal = tokenRspDTO.AccessToken

	if err := SystemSettingDAO.PersistSetting(ctx, PrimaryAccessTokenSystemSettingKey, accessTokenIntenal); err != nil {
		return "", err
	}
	return accessTokenIntenal, nil
}
