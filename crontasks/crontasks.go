package crontasks

import (
	"github.com/go-co-op/gocron"
	"github.com/prometheus/common/log"
	fakernews_mod "github.com/vipau/gatnbot/fakernews-mod"
	"github.com/vipau/gatnbot/settings"
	tb "gopkg.in/telebot.v3"
	"time"
)

var Viernes = &tb.Video{File: tb.FromDisk("jvazquez/viernes.mp4")}
var Sabado = &tb.Video{File: tb.FromDisk("jvazquez/sabado.mp4")}

func sendToAllChats(message interface{}, config settings.Settings, b *tb.Bot) {
	for _, i := range config.Chatid {
		// get group instance from ID
		group := tb.ChatID(i)
		b.Send(group, message)
	}
}

func StartCronProcesses(config settings.Settings, b *tb.Bot) {
	// make a scheduler
	tmz, _ := time.LoadLocation(config.Timezone)
	s := gocron.NewScheduler(tmz)

	// important
	s.Every(1).Day().At("9:00").Do(func() { sendToAllChats("Gattini(tm) reminder: Have you drank water and stretched?", config, b) })
	s.Every(1).Day().At("13:00").Do(func() { sendToAllChats("Gattini(tm) reminder: Have you drank water and stretched?", config, b) })
	s.Every(1).Day().At("16:30").Do(func() { sendToAllChats("Gattini(tm) reminder: Have you drank water and stretched?", config, b) })

	// its friday then
	s.Every(1).Friday().At("08:55").Do(func() { sendToAllChats("https://www.youtube.com/watch?v=1AnG04qnLqI", config, b) })
	s.Every(1).Friday().At("11:00").Do(func() { sendToAllChats(Viernes, config, b) })

	// SABADOOOOOO
	s.Every(1).Saturday().At("10:00").Do(func() { sendToAllChats(Sabado, config, b) })

	// misc shotpost
	s.Every(1).Day().At("13:12").Do(func() { sendToAllChats("A.C.A.B.", config, b) })

	// reload top 500 hacker news articles for the markov chain at midnight
	s.Every(1).Day().At("00:00").Do(func() { fakernews_mod.TrainModel() })

	// start scheduler asynchronously
	log.Info("Starting asynchronous scheduler...")
	s.StartAsync()

}
