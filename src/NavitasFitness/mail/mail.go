package mail

import (
	UserDao "User/Dao"
	"context"
	log "logger"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type Message struct {
	To       *UserDao.UserDTO
	Bcc      []string
	Subject  string
	Body     string
	HTMLBody string
	CustomID string
}

func Send(ctx context.Context, msg *Message) error {

	from := &mailjet.RecipientV31{
		Email: "kontakt@navitasfitness.dk",
		Name:  "noreply - navitasfitness",
	}

	to := mailjet.RecipientV31{
		Email: msg.To.Email,
		Name:  msg.To.Name,
	}

	bccs := make([]mailjet.RecipientV31, len(msg.Bcc))
	for _, reciverBcc := range msg.Bcc {
		toBcc := mailjet.RecipientV31{
			Email: reciverBcc,
			Name:  reciverBcc,
		}
		bccs = append(bccs, toBcc)
	}
	bccsConverted := mailjet.RecipientsV31(bccs)

	mailjetClient := mailjet.NewMailjetClient("62ac9e506457bebbcf1abbace2563ce2", "88a72124cd46d2553b575fdab89d5c11")
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: from,
			To: &mailjet.RecipientsV31{
				to,
			},
			Bcc:      &bccsConverted,
			Subject:  msg.Subject,
			TextPart: msg.Body,
			HTMLPart: msg.HTMLBody,
			CustomID: msg.CustomID,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)

	log.Debugf(ctx, "Data: %+v\n", res)

	return err
}
