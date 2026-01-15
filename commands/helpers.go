package commands

import (
	"fmt"
	tb "gopkg.in/telebot.v3"
	"log/slog"
	"strings"
)

func checkPrintErr(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}

func checkSendErr(err error, b *tb.Bot, c tb.Context, isReply bool, outer ...string) {
	if err != nil {
		opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
		errmsg := "gatnbot: lol an error occurred, check it out yo...\n\n```error\n" + err.Error() + "```"
		if len(outer) != 0 {
			errmsg += "\n\n" + outer[0]
		}
		fmt.Println(errmsg)
		if isReply {
			_, err = b.Reply(c.Message(), errmsg, opts)
			checkPrintErr(err)
		} else {
			_, err = b.Send(c.Chat(), errmsg, opts)
			checkPrintErr(err)
		}
	}
}

func findPrintableName(u *tb.User) string {
	if u.LastName == "" {
		return u.FirstName
	} else {
		return u.FirstName + " " + u.LastName
	}
}

func returnFragments(path string) []string {
	//This cuts off the leading forward slash.
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	//This cuts off the trailing forward slash.
	if strings.HasSuffix(path, "/") {
		cut_off_last_char_len := len(path) - 1
		path = path[:cut_off_last_char_len]
	}
	//We need to isolate the individual components of the path.
	components := strings.Split(path, "/")
	return components
}
