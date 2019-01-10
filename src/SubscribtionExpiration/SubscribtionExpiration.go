package subscriptionExpiration

import (
	"IPN/Transaction"
	"User/Dao"
	"User/Service"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"net/http"
)

// var transactionDao = TransactionDao.GetInstance()
var userDao = UserDao.GetInstance()

func IntegrateRoutes(router *mux.Router) {
	path := "/rest/subscriptionExpiration"

	router.
		Methods("GET").
		Path(path + "/dryRun").
		Name("subscriptionExpiration dryRun").
		HandlerFunc(UserService.AsAdmin(dryRun))

	router.
		Methods("GET").
		Path(path + "/send").
		Name("subscriptionExpiration send").
		HandlerFunc(UserService.AsAdmin(send))

}

func dryRun(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)
	log.Errorf(ctx, "I WAS here -1")

	txnKeys, err := TransactionDao.GetTransactionsAboutToExpire(ctx)
	if err != nil {
		return nil, err
	}

	userKeys := make([]*datastore.Key, 0, len(txnKeys))

	for _, key := range txnKeys {
		userKeys = append(userKeys, key.Parent())
	}

	users, err := userDao.GetByKeys(ctx, userKeys)
	// return userKeys, err
	return users, err
}

func send(w http.ResponseWriter, r *http.Request, user *UserDao.UserDTO) (interface{}, error) {
	ctx := appengine.NewContext(r)
	log.Errorf(ctx, "I WAS here -2")
	return nil, nil
}

// func dryRunImpl() {
// 	transactionDao := TransactionDao.GetInstance()

// 	transactionDao.
// }
