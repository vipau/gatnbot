package crontasks

import (
	"github.com/go-co-op/gocron"
	fakernewsmod "github.com/vipau/gatnbot/fakernews-mod"
	"github.com/vipau/gatnbot/settings"
	tb "gopkg.in/telebot.v3"
	"log/slog"
	"time"
)

func checkErr(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}

var Viernes = &tb.Video{File: tb.FromDisk("jvazquez/viernes.mp4")}
var Sabado = &tb.Video{File: tb.FromDisk("jvazquez/sabado.mp4")}

func sendToAllChats(message interface{}, config settings.Settings, b *tb.Bot) {
	for _, i := range config.Chatid {
		// get group instance from ID
		group := tb.ChatID(i)
		_, err := b.Send(group, message)
		checkErr(err)
	}
}

func StartCronProcesses(config settings.Settings, b *tb.Bot) {
	// make a scheduler
	tmz, _ := time.LoadLocation(config.Timezone)
	s := gocron.NewScheduler(tmz)

	// important
	_, err := s.Every(1).Day().At("9:00").Do(func() { sendToAllChats("Gattini(tm) reminder: Have you drank water and stretched?", config, b) })
	checkErr(err)
	_, err = s.Every(1).Day().At("13:00").Do(func() { sendToAllChats("Gattini(tm) reminder: Have you drank water and stretched?", config, b) })
	checkErr(err)
	_, err = s.Every(1).Day().At("16:30").Do(func() { sendToAllChats("Gattini(tm) reminder: Have you drank water and stretched?", config, b) })
	checkErr(err)

	// its friday then
	_, err = s.Every(1).Friday().At("08:55").Do(func() { sendToAllChats("https://www.youtube.com/watch?v=1AnG04qnLqI", config, b) })
	checkErr(err)
	_, err = s.Every(1).Friday().At("11:00").Do(func() { sendToAllChats(Viernes, config, b) })
	checkErr(err)

	// SABADOOOOOO
	_, err = s.Every(1).Saturday().At("10:00").Do(func() { sendToAllChats(Sabado, config, b) })
	checkErr(err)

	// misc shotpost
	_, err = s.Every(1).Day().At("13:12").Do(func() { sendToAllChats("A.C.A.B.", config, b) })
	checkErr(err)

	// reload top 500 hacker news articles for the markov chain at midnight
	_, err = s.Every(1).Day().At("00:00").Do(func() { fakernewsmod.TrainModel() })
	checkErr(err)

	// start scheduler asynchronously
	slog.Info("Starting asynchronous scheduler...")
	s.StartAsync()

}
