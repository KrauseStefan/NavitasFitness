package main

import (
	"AccessIdOverride"
	"SubscribtionExpiration"
	"errors"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"

	"Auth"
	"DropboxService"
	"Export/csv"
	"IPN"
	"MainPage"
	"User"
)

const emailRegStr = `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`

var emailReg = regexp.MustCompile(emailRegStr)

func validateEmail(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}
	if !emailReg.MatchString(st.String()) {
		return errors.New("Invalid email")
	}

	return nil
}

// http://blog.golang.org/context
// http://blog.golang.org/go-videos-from-google-io-2012

func init() {
	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter().StrictSlash(true)
	subscriptionExpiration.IntegrateRoutes(router)
	MainPageService.IntegrateRoutes(router)
	UserRest.IntegrateRoutes(router)
	Auth.IntegrateRoutes(router)
	IPN.IntegrateRoutes(router)
	csv.IntegrateRoutes(router)
	DropboxService.IntegrateRoutes(router)
	AccessIdOverride.IntegrateRoutes(router)
	http.Handle("/", router)
	//	http.HandleFunc("/rest/", root)

	validator.SetValidationFunc("email", validateEmail)
}

//func root(w http.ResponseWriter, r *http.Request) {
//	//	w.Header().Set("Location", "static/index.html")
//	//	w.WriteHeader(http.StatusFound)
//	fmt.Fprint(w, "Hello, world!")
//}
