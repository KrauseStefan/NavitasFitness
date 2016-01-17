package AuthService

import (
	"github.com/gorilla/mux"

	"net/http"

	"src/Services/Common"
)

type LoginCredentials struct {

}

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/auth"

	router.
		Methods("POST").
		Path(path + "/login").
		Name("loginUser").
		HandlerFunc(Common.ParseFormDataWrap(doLogin))
}

func doLogin(w http.ResponseWriter, r *http.Request, getCred Common.FormDataDecoderFn) {

	loginCredentials := new(LoginCredentials)

	if err := getCred(loginCredentials); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	s := GetSecureCookieInst()

	value := "someRandomKey" //TODO

	cookieName := "Session-Id"

	if encoded, err := s.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/rest/",
		}
		http.SetCookie(w, cookie)
	}

}
