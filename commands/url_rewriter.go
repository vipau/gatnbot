package commands

import (
	"github.com/vipau/gatnbot/settings"
	tb "gopkg.in/telebot.v3"
	"net/url"
	"strings"
)

// handleURLRewriting checks if a message is a URL and rewrites it if needed
func handleURLRewriting(c tb.Context, b *tb.Bot, configmap settings.Settings) error {
	if !settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) &&
		!settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
		return nil
	}

	// Detect if message is a link
	msg := strings.TrimSpace(c.Message().Text)
	u, err := url.Parse(msg)
	if err != nil {
		return nil
	}

	// send link with the telegram preview and with markdown
	opts := &tb.SendOptions{DisableWebPagePreview: false, ParseMode: ""}

	// try for instagram
	if u.Hostname() == "instagram.com" || u.Hostname() == "www.instagram.com" {
		// if URL is a post or reel
		if returnFragments(u.Path)[0] == "p" || returnFragments(u.Path)[0] == "reel" {
			u.Host = "ddinstagram.com"
			b.Delete(c.Message())
			q := u.Query()
			if q.Has("igshid") {
				q.Del("igshid")
				u.RawQuery = q.Encode()
			}
			_, err = b.Send(c.Chat(), "From: "+findPrintableName(c.Sender())+" who did not use ddinstagram and/or remove the 'igshid' tracking tag... wtf\n\n"+u.String(), opts)
			checkSendErr(err, b, c, false)
		}
	}

	// try for twitter
	if u.Hostname() == "twitter.com" || u.Hostname() == "www.twitter.com" || u.Hostname() == "x.com" || u.Hostname() == "www.x.com" {
		// if URL is not profile (more than 1 path fragment)
		if len(returnFragments(u.Path)) > 1 {
			u.Host = "fixupx.com"
			q := u.Query()
			if q.Has("s") {
				q.Del("s")
				u.RawQuery = q.Encode()
			}
			if q.Has("t") {
				q.Del("t")
				u.RawQuery = q.Encode()
			}
			b.Delete(c.Message())
			b.Send(c.Chat(), "From: "+findPrintableName(c.Sender())+" who did not use fixupx... wtf\n\n"+u.String(), opts)
			checkSendErr(err, b, c, false)
		} else {
			// if it's a profile, just remove 's' and 't' trackers
			send := false
			q := u.Query()
			if q.Has("s") {
				q.Del("s")
				u.RawQuery = q.Encode()
				send = true
			}
			if q.Has("t") {
				q.Del("t")
				u.RawQuery = q.Encode()
				send = true
			}
			if send {
				b.Delete(c.Message())
				b.Send(c.Chat(), "From: "+findPrintableName(c.Sender())+" who did not remove s/t trackers from the link... wtf\n\n"+u.String(), opts)
				checkSendErr(err, b, c, false)
			}

		}
	}

	// try for youtube
	if u.Hostname() == "youtube.com" || u.Hostname() == "www.youtube.com" || u.Hostname() == "youtu.be" {
		var send = false
		var sendstring = ""
		q := u.Query()
		if q.Has("si") {
			q.Del("si")
			sendstring += "Removed 'si' tracking tag\n"
			send = true
		}
		if returnFragments(u.Path)[0] == "shorts" {
			q.Add("v", returnFragments(u.Path)[1])
			u.Path = "/watch"
			send = true
			sendstring += "Fixed short preview\n"
		}
		if send {
			b.Delete(c.Message())
			u.RawQuery = q.Encode()
			b.Send(c.Chat(), "From: "+findPrintableName(c.Sender())+"\n"+sendstring+u.String(), opts)
			checkSendErr(err, b, c, false)
		}
	}

	// try for spotify
	if u.Hostname() == "open.spotify.com" {
		q := u.Query()
		if q.Has("si") {
			b.Delete(c.Message())
			q.Del("si")
			u.RawQuery = q.Encode()
			b.Send(c.Chat(), "From: "+findPrintableName(c.Sender())+" who did not remove the 'si' tracking tag... wtf\n\n"+u.String(), opts)
			checkSendErr(err, b, c, false)

		}
	}

	return nil
}
