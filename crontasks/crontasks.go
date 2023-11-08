package crontasks

import (
	"github.com/go-co-op/gocron"
	fakernews_mod "github.com/paualberto/gatnbot/fakernews-mod"
	"github.com/paualberto/gatnbot/sendemail"
	"github.com/paualberto/gatnbot/settings"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func StartCronProcesses(config settings.Settings, b *tb.Bot) {
	// for every group in the array of IDs
	for i, _ := range config.Chatid {

		// make a CET scheduler
		cet, _ := time.LoadLocation("Europe/Rome")
		s := gocron.NewScheduler(cet)

		// its friday then
		s.Every(1).Friday().At("09:00").Do(func() { b.Send(&tb.Chat{ID: config.Chatid[i]}, "https://www.youtube.com/watch?v=1AnG04qnLqI") })

		// misc shitpost
		s.Every(1).Day().At("04:20").Do(func() { b.Send(&tb.Chat{ID: config.Chatid[i]}, "Ricordate di blazzarla duro come lo zio Snoop") })
		s.Every(1).Day().At("13:12").Do(func() { b.Send(&tb.Chat{ID: config.Chatid[i]}, "A.C.A.B.") })

		// reload top 500 hacker news articles for the markov chain at midnight
		s.Every(1).Day().At("00:00").Do(func() { fakernews_mod.TrainModel() })

		// poll gmail API once per minute
		s.Every(1).Minute().Do(func() { sendemail.CheckAndForward(config.Ouremail, config.Chatid, b) })

		// start scheduler asynchronously
		s.StartAsync()
	}
}
