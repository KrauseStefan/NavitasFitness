package AuthService

import (
	"github.com/gorilla/mux"

	"net/http"

	"src/User/Dao"
	"appengine"
	"errors"
	"encoding/json"
	"src/Common"
	"crypto/rand"
	"encoding/hex"
)

//const UserLoggedInSessionKey = "UserLoggedIn"
//const AdminLoggedInSessionKey = "AdminLoggedIn"

const sessionCookieName = "Session-Key"

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

func generateUUID() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	u[8] = (u[8] | 0x80) & 0xBF // what does this do?
	u[6] = (u[6] | 0x40) & 0x4F // what does this do?

	return hex.EncodeToString(u), nil
}

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/auth"

	router.
		Methods("POST").
		Path(path + "/login").
		Name("loginUser").
		HandlerFunc(doLogin)

	router.
		Methods("POST").
		Path(path + "/logout").
		Name("logoutUser").
		HandlerFunc(doLogout)

}

func setSessionCookie(w http.ResponseWriter, uuid string) error {
	var (
		err error
		encoded string
	)
	s := GetSecureCookieInst()

	if(uuid != ""){
		encoded, err = s.Encode(sessionCookieName, uuid)
		if err != nil {
			return err
		}
	}

	cookie := &http.Cookie{
		Name:  sessionCookieName,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	return nil
}

func doLogout(w http.ResponseWriter, r *http.Request) {
	err := setSessionCookie(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func doLogin(w http.ResponseWriter, r *http.Request) {

	var (
		user 	*UserDao.UserDTO
		uuid 	string
		err		error
	)

	ctx := appengine.NewContext(r)

	loginRequestUser := new(UserLogin)

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(loginRequestUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if(!loginRequestUser.hasValues()) {
		http.Error(w, invalidLoginError.Error(), http.StatusBadRequest)
		return
	}

	user, err = UserDao.GetUserByEmail(ctx, loginRequestUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if(user == nil || user.Password != loginRequestUser.Password){
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	uuid, err = generateUUID()
	if err != nil {
		http.Error(w, "Error Generating UUID", http.StatusUnauthorized)
		return
	}

	err = setSessionCookie(w, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = UserDao.SetSessionUUID(ctx, user, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if _, err := Common.WriteJSON(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetSessionUUID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return "", err
	}

	if cookie.Value == "" {
		return "", nil
	}

	uuid := ""

	if err := GetSecureCookieInst().Decode(sessionCookieName, cookie.Value, &uuid); err != nil {
		return "", err
	}

	return uuid, nil
}