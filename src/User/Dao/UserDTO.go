package UserDao

import (
	"DAOHelper"
	"crypto/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"gopkg.in/validator.v2"

	"cloud.google.com/go/datastore"

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
	accessIdValidator.EnsureUpdatedIds(ctx)
	isValid, err := accessIdValidator.ValidateAccessId(ctx, []byte(user.AccessId))
	if err != nil {
		return err
	}

	if !isValid {
		return Invalid_accessId
	}

	err = validator.Validate(user)
	if err != nil {
		err, isErrorMap := err.(validator.ErrorMap)

		if isErrorMap {
			for k, errs := range err {
				if len(errs) > 0 {
					return &DAOHelper.ConstraintError{Field: k, Type: DAOHelper.Invalid, Message: errs.Error()}
				}
			}
		} else {
			return err
		}
	}

	return nil
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
	_, err := rand.Read(user.PasswordSalt)
	if err != nil {
		return err
	}

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

func (user *UserDTO) IsEquivalent(other *UserDTO) bool {
	return user.AccessId == other.AccessId || user.Email == other.Email
}

type UserList []*UserDTO

func (users UserList) Filter(filterFn func(*UserDTO) bool) UserList {
	filteredUsers := make(UserList, 0, len(users))
	for _, user := range users {
		if filterFn(user) {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}
