package AccessIdOverride

import (
	"AccessIdOverride/dao"
	"User/Dao"
	"User/Service"
	"encoding/json"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"net/http"
)

const accessIdKey = "accessId"

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/AccessIdOverride"

	router.
		Methods("GET").
		Path(path).
		Name("Get all Access Id Overrides").
		HandlerFunc(UserService.AsUser(getAllAccessIdOverrideHandler))

	router.
		Methods("POST").
		Path(path).
		Name("Create or update Access Id Overrides").
		HandlerFunc(UserService.AsUser(createOrUpdateAccessIdOverrideHandler))

	router.
		Methods("Delete").
		Path(path + "/{" + accessIdKey + "}").
		Name("Delete Access Id Overrides").
		HandlerFunc(UserService.AsUser(deleteAccessIdOverrideHandler))
}

func getAllAccessIdOverrideHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	return AccessIdOverrideDao.GetAllAccessIdOverrides(ctx)
}

func createOrUpdateAccessIdOverrideHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	accessId := &AccessIdOverrideDao.AccessIdOverride{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(accessId); err != nil {
		return nil, err
	}

	err := AccessIdOverrideDao.CreateOrUpdateAccessIdOverride(ctx, accessId)

	return nil, err
}

func deleteAccessIdOverrideHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)

	accessId := mux.Vars(r)[accessIdKey]

	err := AccessIdOverrideDao.DeleteAccessIdOverride(ctx, accessId)

	return nil, err
}
