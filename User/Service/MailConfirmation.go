package UserService

import (
	"fmt"
	"net/url"

	"appengine"
	"appengine/mail"

	"User/Dao"
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

func SendConfirmationMail(ctx appengine.Context, user *UserDao.UserDTO) error {
	confirmationUrl := createConfirmationURL(user.Key.Encode())

	ctx.Infof("Confirmatin URL: " + confirmationUrl)
	msg := &mail.Message{
		Sender:  "Navitass Fitness <navitas-fitness-aarhus@appspot.gserviceaccount.com>",
		To:      []string{user.Email},
		Subject: "Confirm your registration",
		Body:    fmt.Sprintf(confirmMessage, confirmationUrl),
	}

	if err := mail.Send(ctx, msg); err != nil {
		return err
	}

	return nil
}
