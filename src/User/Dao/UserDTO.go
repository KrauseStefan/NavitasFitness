package UserDao

import (
	"crypto/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"AccessIdValidator"
)

var accessIdValidator = AccessIdValidator.GetInstance()

type UserDTO struct {
	Name                string         `json:"name" datastore:",noindex"`
	Email               string         `json:"email" validate:"email"`
	AccessId            string         `json:"accessId"`
	Password            string         `json:"password,omitempty" datastore:"-"`
	PasswordHash        []byte         `json:"-" datastore:",noindex"`
	PasswordSalt        []byte         `json:"-" datastore:",noindex"`
	Key                 *datastore.Key `json:"-" datastore:"-"`
	CurrentSessionUUID  string         `json:"-"`
	IsAdmin             bool           `json:"-" datastore:",noindex"`
	Verified            bool           `json:"-" datastore:",noindex"`
	PasswordResetTime   time.Time      `json:"-" datastore:",noindex"`
	PasswordResetSecret string         `json:"-" datastore:",noindex"`
}

func (user *UserDTO) ValidateUser(ctx context.Context) error {
	isValid, err := accessIdValidator.ValidateAccessIdPrimary(ctx, []byte(user.AccessId))

	if !isValid && err == nil {
		isValid, err = accessIdValidator.ValidateAccessIdSecondary(ctx, []byte(user.AccessId))
	}

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

func (user *UserDTO) UpdatePasswordHash(password string) error {
	if password == "" {
		return passwordCanNotBeEmpty
	}
	// https://crackstation.net/hashing-security.htm
	user.PasswordSalt = make([]byte, PW_SALT_BYTES)
	rand.Read(user.PasswordSalt)

	passwordHash, err := bcrypt.GenerateFromPassword(user.getPasswordWithSalt([]byte(password)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash

	return nil
}

func (user *UserDTO) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, user.getPasswordWithSalt([]byte(password)))
}
