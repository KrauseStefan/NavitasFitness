package UserService

import (
	"User/Dao"
	"appengine"
	"appengine/mail"
	"fmt"
	"net/url"
)

const KeyFieldName = "passwordResetKey"
const secretFieldName = "passwordResetSecret"

const passwordResetMessageTpl = `
Thank you for creating an account!
Please confirm your email address by clicking on the link below:

%s
`

func createPasswordResetURL(form *url.Values) string {
	return "https://navitas-fitness-aarhus.appspot.com/main-page/?" + form.Encode()
}

func SendPasswordResetMail(ctx appengine.Context, user *UserDao.UserDTO, secret string) error {
	form := url.Values{}
	form.Add(KeyFieldName, user.Key.Encode())
	form.Add(secretFieldName, secret)

	passwordResetUrl := createPasswordResetURL(&form)

	ctx.Infof("Password Reset URL: " + passwordResetUrl)
	msg := &mail.Message{
		Sender:  "noreply - Navitass Fitness <navitas-fitness-aarhus@appspot.gserviceaccount.com>",
		To:      []string{user.Email},
		Subject: "Confirm your registration",
		Body:    fmt.Sprintf(passwordResetMessageTpl, passwordResetUrl),
	}

	if err := mail.Send(ctx, msg); err != nil {
		return err
	}

	return nil
}
