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

// try to reuse an existing gmail service if active
var ctx = context.Background()
var gmailsrv *gmail.Service = nil

// CheckAndForward checks for unread emails matching query and forward them to group
func CheckAndForward(ourmail string, chatids []int64, b *tb.Bot) {
	if gmailsrv == nil {
		gmailsrv = gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)
		// full gmail scope needed for read + mark as read
		// know a better combination of scopes? please open an issue :)
	}

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

	// build a request to mark a msg as read
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

		// mark as read on server
		_, err := inboxer.MarkAs(gmailsrv, msg, req)
		if err != nil {
			panic(err)
		}
	}
}
