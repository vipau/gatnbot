package commands

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/vipau/gatnbot/crontasks"
	fakernewsmod "github.com/vipau/gatnbot/fakernews-mod"
	"github.com/vipau/gatnbot/settings"
	tb "gopkg.in/telebot.v3"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
)

// registerSimpleCommands registers all simple command handlers
func registerSimpleCommands(b *tb.Bot, configmap settings.Settings) {
	b.Handle("/myid", func(c tb.Context) error {
		opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
		_, err := b.Send(c.Message().Chat, "Your ID: `"+strconv.FormatInt(c.Sender().ID, 10)+"`", opts)
		return err
	})

	b.Handle("/links", func(c tb.Context) error {
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
			opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
			_, err := b.Send(c.Message().Chat, configmap.Linksmsg, opts)
			checkPrintErr(err)
		}
		return nil
	})

	b.Handle("/turbo", func(c tb.Context) error {
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
			vmin := 4
			vmax := 57
			rando := rand.Intn(vmax-vmin+1) + vmin
			_, err := b.Send(c.Message().Chat, fmt.Sprintf("this chat is now cringe-protected for %d minutes thanks the power of TURBO", rando))
			checkPrintErr(err)
		}
		return nil
	})

	b.Handle("/hackernews", func(c tb.Context) error {
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
			if _, err := os.Stat("model.json"); err == nil {
			} else if os.IsNotExist(err) {
				fakernewsmod.TrainModel()
			} // train the model first if it doesn't exist

			out := fakernewsmod.GenerateNews()
			_, err := b.Send(c.Message().Chat, out)
			checkPrintErr(err)
			return nil
		}
		return nil
	})

	b.Handle("/glados", func(c tb.Context) error {
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
			gladosLine := GetGladosVoiceline()
			a := &tb.Audio{File: tb.FromDisk("glados/" + gladosLine), Title: gladosLine, Performer: "GLaDOS"}
			_, err := b.Send(c.Message().Chat, a)
			if err != nil {
				slog.Error(err.Error())
				slog.Error(gladosLine)
				checkSendErr(errors.Wrap(err, "error playing glados line"), b, c, true)
			}
		}
		return nil
	})

	// manual viernes/sabado invocation
	b.Handle("/viernes", func(c tb.Context) error {
		_, err := b.Send(c.Message().Chat, crontasks.Viernes)
		return err
	})
	b.Handle("/sabado", func(c tb.Context) error {
		_, err := b.Send(c.Message().Chat, crontasks.Sabado)
		return err
	})

	b.Handle("/coin", func(c tb.Context) error {
		// Note: rand.Seed is deprecated in Go 1.20+ and no longer needed
		outcome := rand.Intn(2)
		output := ""
		if outcome == 0 {
			output = "Heads"
		} else {
			output = "Tails"
		}
		b.Send(c.Message().Chat, output)
		return nil
	})
}
