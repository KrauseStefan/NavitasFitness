package UserService

import (
	"fmt"
	"net/url"

	"golang.org/x/net/context"

	"github.com/KrauseStefan/NavitasFitness/NavitasFitness/mail"
	UserDao "github.com/KrauseStefan/NavitasFitness/User/Dao"
	log "github.com/KrauseStefan/NavitasFitness/logger"
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

func SendPasswordResetMail(ctx context.Context, user *UserDao.UserDTO, secret string) error {
	form := url.Values{}
	form.Add(KeyFieldName, user.Key.Encode())
	form.Add(secretFieldName, secret)

	passwordResetUrl := createPasswordResetURL(&form)

	log.Infof(ctx, "Password Reset URL: "+passwordResetUrl)
	msg := &mail.Message{
		To:       user,
		Subject:  "Password Reset Request",
		Body:     fmt.Sprintf(passwordResetMessageTpl, passwordResetUrl),
		CustomID: "PasswordResetRequest",
	}

	return mail.Send(ctx, msg)
}
