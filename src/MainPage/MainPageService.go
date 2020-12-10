package MainPageService

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"

	"github.com/KrauseStefan/NavitasFitness/AppEngineHelper"
	nf_datastore "github.com/KrauseStefan/NavitasFitness/NavitasFitness/datastore"
	UserDao "github.com/KrauseStefan/NavitasFitness/User/Dao"
	UserService "github.com/KrauseStefan/NavitasFitness/User/Service"
)

type MainPageEntry struct {
	Key          string    `json:"key" datastore:"-"`
	LastEditedBy string    `json:"lastEditedBy"`
	Content      string    `json:"content" datastore:",noindex"`
	Date         time.Time `json:"date"`
}

func (mainPage MainPageEntry) hasId() bool {
	return len(mainPage.Key) > 0
}

const MAIN_PAGE_KIND = "mainPage"
const MAIN_PAGE_PARENT_STRING_ID = "default_main_page_entry"

var mainPageCollectionParentKey = datastore.NameKey(MAIN_PAGE_KIND, MAIN_PAGE_PARENT_STRING_ID, nil)

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/mainPage"

	router.
		Methods("GET").
		Path(path).
		Name("GetMainPage").
		HandlerFunc(AppEngineHelper.HandlerW(getMainPage))

	router.
		Methods("PUT").
		Path(path).
		Name("UpdateMainPage").
		HandlerFunc(UserService.AsAdmin(updateMainPage))
}

func getMainPage(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery(MAIN_PAGE_KIND).
		Ancestor(mainPageCollectionParentKey)

	frontPageEntries := make([]MainPageEntry, 0, 1)

	keys, err := dsClient.GetAll(ctx, query, &frontPageEntries)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		frontPageEntries[i].Key = strconv.FormatInt(key.ID, 10)
	}

	if len(keys) == 0 {
		frontPageEntries = append(frontPageEntries, MainPageEntry{})
	}

	return frontPageEntries[0], err
}

func updateMainPage(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := r.Context()
	var mainPage MainPageEntry
	var key *datastore.Key

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mainPage)
	if err != nil {
		return nil, err
	}

	if mainPage.hasId() {
		intId, _ := strconv.ParseInt(mainPage.Key, 10, 64)
		key = datastore.IDKey(MAIN_PAGE_KIND, intId, mainPageCollectionParentKey)
	} else {
		key = datastore.IncompleteKey(MAIN_PAGE_KIND, mainPageCollectionParentKey)
	}
	mainPage.Date = time.Now()
	mainPage.LastEditedBy = user.Email

	dsClient, err := nf_datastore.GetDsClient()
	if err != nil {
		return nil, err
	}
	key, err = dsClient.Put(ctx, key, &mainPage)
	return mainPage, err
}
