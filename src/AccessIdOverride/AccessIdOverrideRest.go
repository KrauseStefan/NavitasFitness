package AccessIdOverride

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	AccessIdOverrideDao "github.com/KrauseStefan/NavitasFitness/AccessIdOverride/dao"
	UserDao "github.com/KrauseStefan/NavitasFitness/User/Dao"
	UserService "github.com/KrauseStefan/NavitasFitness/User/Service"
)

const accessIdKey = "accessId"

var accessIdOverrideDao = AccessIdOverrideDao.GetInstance()

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
	ctx := r.Context()

	return accessIdOverrideDao.GetAllAccessIdOverrides(ctx)
}

func createOrUpdateAccessIdOverrideHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := r.Context()

	accessId := &AccessIdOverrideDao.AccessIdOverride{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(accessId); err != nil {
		return nil, err
	}

	err := accessIdOverrideDao.CreateOrUpdateAccessIdOverride(ctx, accessId)

	return nil, err
}

func deleteAccessIdOverrideHandler(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := r.Context()

	accessId := mux.Vars(r)[accessIdKey]

	err := accessIdOverrideDao.DeleteAccessIdOverride(ctx, accessId)

	return nil, err
}
