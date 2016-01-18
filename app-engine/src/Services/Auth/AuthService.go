package AuthService

import (
	"github.com/gorilla/mux"

	"net/http"

	"src/Services/Common"
	"src/Services/User"
	"appengine"
)


func IntegrateRoutes(router *mux.Router) {
	path := "/rest/auth"

	router.
		Methods("POST").
		Path(path + "/login").
		Name("loginUser").
		HandlerFunc(Common.ParseFormDataWrap(doLogin))
}

func doLogin(w http.ResponseWriter, r *http.Request, getCred Common.FormDataDecoderFn) {

	ctx := appengine.NewContext(r)

	loginRequestUser := new(UserService.UserDTO)

	if err := getCred(loginRequestUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err := UserService.GetUserByUserName(ctx, loginRequestUser.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if(user.UserName != loginRequestUser.UserName || user.Password != loginRequestUser.Password){
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
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
