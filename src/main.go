package main

import (
	"errors"
	"fmt"
	"log"
	"path"
	"strings"

	"math/rand"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"

	"github.com/KrauseStefan/NavitasFitness/AccessIdOverride"
	"github.com/KrauseStefan/NavitasFitness/Auth"
	"github.com/KrauseStefan/NavitasFitness/DropboxService"
	"github.com/KrauseStefan/NavitasFitness/Export/csv"
	"github.com/KrauseStefan/NavitasFitness/IPN"
	MainPageService "github.com/KrauseStefan/NavitasFitness/MainPage"
	"github.com/KrauseStefan/NavitasFitness/NavitasFitness/spaHandler"
	subscriptionExpiration "github.com/KrauseStefan/NavitasFitness/SubscribtionExpiration"
	UserRest "github.com/KrauseStefan/NavitasFitness/User"
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

func init() {
	log.Printf("Init")
	// goPaths := strings.Split(os.Getenv("GOPATH"), ":")
	webappFolder := "./webapp"
	// for _, goPath := range goPaths {
	// 	folder := path.Join(goPath, "/src/webapp")
	// 	_, err := os.Stat(folder)
	// 	if err == nil {
	// 		webappFolder = folder
	// 	}

	// 	fmt.Println("goPath:", goPath)
	// }
	// if webappFolder == "./webapp" {
	folder := path.Join(os.Args[0], "webapp")
	_, err := os.Stat(folder)
	if err == nil {
		webappFolder = folder
	}
	// }

	fmt.Println("webappFolder:", webappFolder)

	r := mux.NewRouter().StrictSlash(true)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.RequestURI, "/rest/") {
				log.Println(r.RequestURI)
			}
			next.ServeHTTP(w, r)
		})
	})

	subscriptionExpiration.IntegrateRoutes(r)
	MainPageService.IntegrateRoutes(r)
	UserRest.IntegrateRoutes(r)
	Auth.IntegrateRoutes(r)
	IPN.IntegrateRoutes(r)
	csv.IntegrateRoutes(r)
	DropboxService.IntegrateRoutes(r)
	AccessIdOverride.IntegrateRoutes(r)

	spa := spaHandler.SpaHandler{StaticPath: webappFolder, IndexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	http.Handle("/", r)

	validator.SetValidationFunc("email", validateEmail)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
