package MainPageService

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"AppEngineHelper"
	UserDao "User/Dao"
	UserService "User/Service"
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

var mainPageCollectionParentKey = AppEngineHelper.CollectionParentKeyGetFnGenerator(MAIN_PAGE_KIND, MAIN_PAGE_PARENT_STRING_ID, 0)
var mainPageIntIDToKeyInt64 = AppEngineHelper.IntIDToKeyInt64(MAIN_PAGE_KIND, mainPageCollectionParentKey)

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
	ctx := appengine.NewContext(r)
	q := datastore.NewQuery(MAIN_PAGE_KIND).Ancestor(mainPageCollectionParentKey(ctx))

	frontPageEntries := make([]MainPageEntry, 0, 1)

	keys, err := q.GetAll(ctx, &frontPageEntries)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		frontPageEntries[i].Key = strconv.FormatInt(key.IntID(), 10)
	}

	if len(keys) == 0 {
		frontPageEntries = append(frontPageEntries, MainPageEntry{})
	}

	return frontPageEntries[0], err
}

func updateMainPage(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)
	var mainPage MainPageEntry
	var key *datastore.Key

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mainPage)
	if err != nil {
		return nil, err
	}

	if mainPage.hasId() {
		key = mainPageIntIDToKeyInt64(ctx, mainPage.Key)
	} else {
		key = datastore.NewIncompleteKey(ctx, MAIN_PAGE_KIND, mainPageCollectionParentKey(ctx))
	}
	mainPage.Date = time.Now()
	mainPage.LastEditedBy = user.Email

	key, err = datastore.Put(ctx, key, &mainPage)
	return mainPage, err
}
