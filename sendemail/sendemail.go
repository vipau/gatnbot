package sendemail

import (
	"context"
	"fmt"
	"gitlab.com/hartsfield/gmailAPI"
	"gitlab.com/hartsfield/inboxer"
	gmail "google.golang.org/api/gmail/v1"
	tb "gopkg.in/tucnak/telebot.v2"
	"html"
)

var gmailsrv *gmail.Service = nil

func CheckAndForward(ourmail string, chatids []int64, b *tb.Bot) {
	ctx := context.Background()
	gmailsrv = gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)

	// die if the email alias is empty
	if ourmail == "" {
		panic("Email recipient is empty! Safety warning!")
	}

	// build a query so that only if the mail is to our alias AND unread it gets returned by Google
	msgs, err := inboxer.Query(gmailsrv, fmt.Sprintf("to:%s is:unread", ourmail))
	if err != nil {
		fmt.Println(err)
		return // exit function if the query fails
	}

	// build a request to mark it as read (so we don't process it again)
	req := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}

	// our sending options
	opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "HTML"}

	// Range over the messages
	for _, msg := range msgs {
		md := inboxer.GetPartialMetadata(msg)
		message := "<b><u>Ao c'è posta per Gattini</u></b>" + string('\n') + string('\n')
		message += fmt.Sprintf("<i><u>Probabilmente arriva da:</u></i>\n%s\n\n", html.EscapeString(md.From))
		message += fmt.Sprintf("<i><u>Il titolo dice:</u></i>\n%s", html.EscapeString(md.Subject))

		// send to every group in array
		for i := range chatids {
			b.Send(&tb.Chat{ID: chatids[i]}, message, opts)
		}
		_, err := inboxer.MarkAs(gmailsrv, msg, req)
		if err != nil {
			panic(err)
		}
	}
}
