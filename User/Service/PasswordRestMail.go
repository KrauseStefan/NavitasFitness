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
A Password reset request has been made for this mail, if you did not request this please ignore this mail.

If you did request the change then please click the below link to reset your accunt password.

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
		Subject: "Password Reset Request",
		Body:    fmt.Sprintf(passwordResetMessageTpl, passwordResetUrl),
	}

	if err := mail.Send(ctx, msg); err != nil {
		return err
	}

	return nil
}
