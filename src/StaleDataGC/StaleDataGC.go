package StaleDataGC

import (
	"IPN/Transaction"
	"User/Dao"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

var (
	userDao        = UserDao.GetInstance()
	transactionDao = TransactionDao.GetInstance()
)

func FindStaleUsers(ctx context.Context) error {
	keys, users, err := userDao.GetAll(ctx)
	if err != nil {
		return err
	}

	twoDays := 2 * 24 * time.Hour

	threshold := time.Now().Add(twoDays)

	usersToDelete := make([]*datastore.Key, 0, 20)

	for i, user := range users {
		if !user.Verified && user.CreationDate.Before(threshold) {
			usersToDelete = append(usersToDelete, keys[i])
		}
	}

	return userDao.DeleteUsers(ctx, usersToDelete)
}
