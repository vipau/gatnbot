package crontasks

import (
	"github.com/go-co-op/gocron"
	fakernewsmod "github.com/vipau/gatnbot/fakernews-mod"
	"github.com/vipau/gatnbot/settings"
	tb "gopkg.in/telebot.v3"
	"log/slog"
	"time"
)

var Viernes = &tb.Video{File: tb.FromDisk("jvazquez/viernes.mp4")}
var Sabado = &tb.Video{File: tb.FromDisk("jvazquez/sabado.mp4")}

func sendToAllChats(message interface{}, config settings.Settings, b *tb.Bot) {
	for _, i := range config.Chatid {
		// get group instance from ID
		group := tb.ChatID(i)
		_, err := b.Send(group, message)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func StartCronProcesses(config settings.Settings, b *tb.Bot) {
	// make a scheduler
	tmz, _ := time.LoadLocation(config.Timezone)
	s := gocron.NewScheduler(tmz)

	warning := "Gattini(tm) reminder: Have you drunk water and stretched?"

	// important
	_, err := s.Every(1).Monday().Tuesday().Wednesday().Thursday().Friday().
		At("11:00").At("13:30").At("16:30").Do(func() { sendToAllChats(warning, config, b) })
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.Every(1).Saturday().Sunday().
		At("14:30").Do(func() { sendToAllChats(warning, config, b) })
	if err != nil {
		slog.Error(err.Error())
	}

	// its friday then
	_, err = s.Every(1).Friday().At("08:55").Do(func() { sendToAllChats("https://www.youtube.com/watch?v=1AnG04qnLqI", config, b) })
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.Every(1).Friday().At("11:00").Do(func() { sendToAllChats(Viernes, config, b) })
	if err != nil {
		slog.Error(err.Error())
	}

	// SABADOOOOOO
	_, err = s.Every(1).Saturday().At("10:00").Do(func() { sendToAllChats(Sabado, config, b) })
	if err != nil {
		slog.Error(err.Error())
	}

	// reload top 500 hacker news articles for the markov chain at midnight
	_, err = s.Every(1).Day().At("00:00").Do(func() { fakernewsmod.TrainModel() })
	if err != nil {
		slog.Error(err.Error())
	}

	// start scheduler asynchronously
	slog.Info("Starting asynchronous scheduler...")
	s.StartAsync()

}
