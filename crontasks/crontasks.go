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
	// make a scheduler
	tmz, _ := time.LoadLocation(config.Timezone)
	s := gocron.NewScheduler(tmz)

	// poll gmail API once per minute
	// already acts on all chats, no need to do it in the for
	s.Every(1).Minute().Do(func() { sendemail.CheckAndForward(config.Ouremail, config.Chatid, b) })

	// for every group in the array of IDs
	for _, i := range config.Chatid {
		// get group instance from ID
		group := tb.ChatID(i)

		// its friday then
		s.Every(1).Friday().At("09:00").Do(func() { b.Send(group, "https://www.youtube.com/watch?v=1AnG04qnLqI") })

		// misc shotpost
		s.Every(1).Day().At("13:12").Do(func() { b.Send(group, "A.C.A.B.") })

		// reload top 500 hacker news articles for the markov chain at midnight
		s.Every(1).Day().At("00:00").Do(func() { fakernews_mod.TrainModel() })

	}
	// start scheduler asynchronously
	s.StartAsync()
}
