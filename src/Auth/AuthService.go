package Auth

import (
	"github.com/gorilla/mux"

	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KrauseStefan/NavitasFitness/AppEngineHelper"
	UserDao "github.com/KrauseStefan/NavitasFitness/User/Dao"
	log "github.com/KrauseStefan/NavitasFitness/logger"

	"cloud.google.com/go/datastore"
)

var userDAO = UserDao.GetInstance()

const sessionCookieName = "Session-Key"

type UserLogin struct {
	Password string `json:"password"`
	AccessId string `json:"accessId"`
}

func (ul UserLogin) hasValues() bool {
	return len(ul.AccessId) > 0 && len(ul.Password) > 0
}

var (
	invalidLoginError = errors.New("Invalid login information, both password and email must be provided")
)

func generateUUID() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

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

func UpdateSessionDataUserKey(r *http.Request, w http.ResponseWriter, userKey *datastore.Key) error {
	sessionData, err := GetSessionData(r)
	if err != nil {
		return err
	}

	sessionData.UserKey = userKey

	return setSessionCookie(w, &sessionData)
}

func setSessionCookie(w http.ResponseWriter, sessionData *SessionData) error {
	const MaxAgeDeleteNow = -1
	const MaxAgeDefault = 0
	var (
		encoded string
		maxAge  int = MaxAgeDefault
	)
	s, err := GetSecureCookieInst()
	if err != nil {
		return err
	}

	if sessionData != nil && sessionData.UserKey != nil {
		encoded, err = s.Encode(sessionCookieName, sessionData)
		if err != nil {
			return err
		}
	} else {
		maxAge = MaxAgeDeleteNow
	}

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    encoded,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return nil
}

func doLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionData, err := GetSessionData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := UserDao.GetInstance().GetByKey(ctx, sessionData.UserKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionData.Uuid = ""
	user.CurrentSessionUUID = ""

	if err := setSessionCookie(w, &sessionData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = userDAO.SaveUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
}

func doLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	loginRequestUser := new(UserLogin)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(loginRequestUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !loginRequestUser.hasValues() {
		http.Error(w, invalidLoginError.Error(), http.StatusBadRequest)
		return
	}

	user, err := userDAO.GetByAccessId(ctx, loginRequestUser.AccessId)
	if err != nil && err != UserDao.UserNotFoundError {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil || err == UserDao.UserNotFoundError {
		log.Warningf(ctx, "Failed to login, %s does not exist in DB", loginRequestUser.AccessId)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.Verified {
		log.Warningf(ctx, "Failed to login, %s email is not verified", loginRequestUser.AccessId)
		http.Error(w, "Email is not verified", http.StatusForbidden)
		return
	}

	if err := user.VerifyPassword(loginRequestUser.Password); err != nil {
		log.Errorf(ctx, "Failed to login, %s Invalid password", loginRequestUser.AccessId)
		log.Errorf(ctx, err.Error())
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	uuid, err := generateUUID()
	if err != nil {
		http.Error(w, "Error Generating UUID", http.StatusInternalServerError)
		return
	}

	sessionData := &SessionData{
		Uuid:    uuid,
		UserKey: user.Key,
	}

	err = setSessionCookie(w, sessionData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = userDAO.SetSessionUUID(ctx, user, sessionData.Uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Errorf(ctx, "login: %+v", user)

	if _, err := AppEngineHelper.WriteJSON(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetSessionData(r *http.Request) (SessionData, error) {
	sessionData := &SessionData{}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return *sessionData, nil
		}
		return *sessionData, err
	}

	if cookie.Value == "" {
		return *sessionData, nil
	}

	s, err := GetSecureCookieInst()
	if err != nil {
		return *sessionData, err
	}

	if s.Decode(sessionCookieName, cookie.Value, &sessionData); err != nil {
		ctx := r.Context()
		log.Errorf(ctx, "Coockie decode error: "+err.Error())
		return *sessionData, nil
	}

	return *sessionData, nil
}
