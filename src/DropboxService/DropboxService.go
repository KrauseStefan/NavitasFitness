package DropboxService

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"AccessIdValidator"
	"AppEngineHelper"
	"ConfigurationReader"
	"Dropbox"
	"Export/csv"
	UserDao "User/Dao"
	UserService "User/Service"
)

const (
	authenticateUrl = "https://www.dropbox.com/oauth2/authorize"

	redirectUriBase     = "https://navitas-fitness-aarhus.appspot.com"
	redirectUriBaseTest = "http://localhost:8080"
	path                = "/rest/dropbox"
	tokenCallback       = "/tokenCallback"
)

func IntegrateRoutes(router *mux.Router) {

	router.
		Methods("GET").
		Path(path + "/authenticate").
		Name("Authenticate with dropbox redirect").
		HandlerFunc(asAdminIfAlreadyConfigured(authorizeWithDropboxHandler))

	router.
		Methods("GET").
		Path(path + tokenCallback).
		Name("Authenticate with dropbox callback").
		HandlerFunc(asAdminIfAlreadyConfigured(authorizationCallbackHandler))
}

func asAdminIfAlreadyConfigured(f func(http.ResponseWriter, *http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokens, err := Dropbox.GetAccessTokens(ctx)
		if err != nil || len(tokens) > 0 {
			UserService.AsAdmin(func(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
				return f(w, r)
			})(w, r)
		} else {
			AppEngineHelper.HandlerW(f)(w, r)
		}
	}
}

func getRedirectUri(r *http.Request) string {
	if strings.Contains(r.Host, "localhost") {
		return redirectUriBaseTest + path + tokenCallback
	} else {
		return redirectUriBase + path + tokenCallback
	}
}

func authorizeWithDropboxHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	conf, err := ConfigurationReader.GetConfiguration()
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"response_type": "code", // token or code
		"client_id":     conf.ClientKey,
		"redirect_uri":  getRedirectUri(r),
		//"state": fmt.Sprint("%i", rand.Int63()), // Up to 500 bytes of arbitrary data that will be passed back to your redirect URI (CSRF protection)
	}

	paramStr := AppEngineHelper.CreateQueryParamString(params)

	http.Redirect(w, r, authenticateUrl+"?"+paramStr, http.StatusTemporaryRedirect)
	return nil, nil
}

func authorizationCallbackHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	r.ParseForm()

	code := r.Form["code"][0]
	token, err := Dropbox.RetrieveAccessToken(ctx, code, getRedirectUri(r))
	if err != nil {
		return nil, err
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	// Below would be best as a real callback but will cause cyclic dependencies
	// For now this will do
	AccessIdValidator.PushMissingSampleData(ctx, token)
	err = csv.CreateAndUploadFile(ctx, nil)

	return nil, err
}
