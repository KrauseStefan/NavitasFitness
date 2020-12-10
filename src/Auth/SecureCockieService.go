package Auth

import (
	"github.com/gorilla/securecookie"

	"github.com/KrauseStefan/NavitasFitness/ConfigurationReader"
)

func GetSecureCookieInst() (*securecookie.SecureCookie, error) {
	var conf, err = ConfigurationReader.GetConfiguration()

	if err != nil {
		return nil, err
	}

	//The hashKey is required, used to authenticate the cookie value using HMAC
	// It is recommended to use a key with 32 or 64 bytes.
	var hashKey, _ = conf.GetAuthCookieSecret()

	//The blockKey is optional, used to encrypt the cookie value -- set it to nil to not use encryption.
	//If set, the length must correspond to the block size of the encryption algorithm.
	//For AES, used by default, valid lengths are 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	// var blockKey = []byte("a-lot-secret")
	var s = securecookie.New(hashKey, nil)

	return s, nil
}
