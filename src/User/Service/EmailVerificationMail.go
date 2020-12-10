package UserService

import (
	"fmt"
	"net/url"

	"golang.org/x/net/context"

	"github.com/KrauseStefan/NavitasFitness/NavitasFitness/mail"
	UserDao "github.com/KrauseStefan/NavitasFitness/User/Dao"
	log "github.com/KrauseStefan/NavitasFitness/logger"
)

const confirmMessage = `
Thank you for creating an account!
Please confirm your email address by clicking on the link below:

%s
`

func createConfirmationURL(key string) string {
	form := url.Values{}
	form.Set("code", key)
	return "https://navitas-fitness-aarhus.appspot.com/rest/user/verify?" + form.Encode()
}

func SendConfirmationMail(ctx context.Context, user *UserDao.UserDTO) error {
	confirmationUrl := createConfirmationURL(user.Key.Encode())

	log.Infof(ctx, "Password Reset URL: "+confirmationUrl)
	msg := &mail.Message{
		To:       user,
		Subject:  "Confirm your registration",
		Body:     fmt.Sprintf(confirmMessage, confirmationUrl),
		CustomID: "ConfirmYourRegistration",
	}

	if err := mail.Send(ctx, msg); err != nil {
		return err
	}

	return nil
}
