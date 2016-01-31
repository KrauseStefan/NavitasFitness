package AuthService

import (
	"github.com/gorilla/mux"

	"net/http"

	"src/Services/User"
	"appengine"
	"errors"
	"encoding/json"
)

const UserLoggedInSessionKey = "UserLoggedIn"
const AdminLoggedInSessionKey = "AdminLoggedIn"


type UserLogin struct {
	Password	string `json:"password"`
	Email			string `json:"email"`
}

func (ul UserLogin) hasValues() bool {
	return len(ul.Email) > 0 && len(ul.Password) > 0
}

var (
	invalidLoginError	= errors.New("Invalid login information, both password and email must be provided")
)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/auth"

	router.
		Methods("POST").
		Path(path + "/login").
		Name("loginUser").
		HandlerFunc(doLogin)
}

func doLogin(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	loginRequestUser := new(UserLogin)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(loginRequestUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if(!loginRequestUser.hasValues()) {
		http.Error(w, invalidLoginError.Error(), http.StatusBadRequest)
		return
	}

	user, err := UserService.GetUserByEmail(ctx, loginRequestUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if(user == nil || user.Password != loginRequestUser.Password){
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	s := GetSecureCookieInst()

	value := UserLoggedInSessionKey

	cookieName := "Session-Key"

	if encoded, err := s.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/rest/",
		}
		http.SetCookie(w, cookie)
	}

}

func isLoggedIn() {

}
