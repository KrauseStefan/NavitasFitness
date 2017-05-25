package DropboxService

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"appengine"

	"AppEngineHelper"
	"Dropbox"
	"User/Dao"
	"User/Service"
)

const (
	authenticateUrl = "https://www.dropbox.com/oauth2/authorize"

	clientId = "v34s5hrxzkjw8ie"

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

func asAdminIfAlreadyConfigured(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		value, err := Dropbox.GetAccessToken(ctx)
		if err != nil || value != "" {
			UserService.AsAdmin(func(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) {
				f(w, r)
			})(w, r)
		} else {
			f(w, r)
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

func authorizeWithDropboxHandler(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{
		"response_type": "code", // token or code
		"client_id":     clientId,
		"redirect_uri":  getRedirectUri(r),
		//"state": fmt.Sprint("%i", rand.Int63()), // Up to 500 bytes of arbitrary data that will be passed back to your redirect URI (CSRF protection)
	}

	paramStr := AppEngineHelper.CreateQueryParamString(params)

	http.Redirect(w, r, authenticateUrl+"?"+paramStr, http.StatusTemporaryRedirect)
}

func authorizationCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	r.ParseForm()

	code := r.Form["code"][0]
	if err := Dropbox.RetrieveAccessToken(ctx, code, getRedirectUri(r)); err != nil {
		ctx.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
