package UserDao

import (
	"crypto/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"

	"appengine"
	"appengine/datastore"

	"AccessIdValidator"
)

type UserDTO struct {
	Name               string         `json:"name",datastore:",noindex",validate:"min=2"`
	Email              string         `json:"email",validate:"email"`
	AccessId           string         `json:"accessId"`
	Password           string         `json:"password,omitempty",datastore:",noindex",validate:"min=2"`
	PasswordHash       []byte         `json:"-",datastore:",noindex"`
	PasswordSalt       []byte         `json:"-",datastore:",noindex"`
	Key                *datastore.Key `json:"-",datastore:"-"`
	CreatedDate        time.Time      `json:"-"`
	CurrentSessionUUID string         `json:"-"`
	IsAdmin            bool           `json:"-"`
	Verified           bool           `json:"-"`
}

func (user *UserDTO) ValidateUser(ctx appengine.Context) error {
	isValid, err := AccessIdValidator.ValidateAccessId(ctx, []byte(user.AccessId))
	if err != nil {
		return err
	}
	if !isValid {
		return Invalid_accessId
	}

	return validator.Validate(user)
}

func (user *UserDTO) hasKey() bool {
	return user.Key != nil
}

func (user *UserDTO) getPasswordWithSalt(password []byte) []byte {
	return append(user.PasswordSalt, password...)
}

func (user *UserDTO) UpdatePasswordHash(password []byte) error {
	if password == nil && user.Password != "" {
		password = []byte(user.Password)
	}
	if password == nil {
		return passwordCanNotBeEmpty
	}
	// https://crackstation.net/hashing-security.htm
	user.PasswordSalt = make([]byte, PW_SALT_BYTES)
	rand.Read(user.PasswordSalt)

	passwordHash, err := bcrypt.GenerateFromPassword(user.getPasswordWithSalt(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	user.Password = ""

	return nil
}

func (user *UserDTO) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, user.getPasswordWithSalt([]byte(password)))
}
