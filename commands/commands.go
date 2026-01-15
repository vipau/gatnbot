package commands

import (
	"github.com/vipau/gatnbot/settings"
	tb "gopkg.in/telebot.v3"
	"log"
	"time"
)

// HandleCommands sets endpoints handled by the bot
func HandleCommands(configmap settings.Settings) *tb.Bot {
	// create new bot
	b, err := tb.NewBot(tb.Settings{
		// If field is empty it equals to "https://api.telegram.org".
		URL: configmap.Apiurl,

		Token:  configmap.Bottoken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	// die if bot is unable to initialize
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// start handling our custom commands

	// URL rewriting handler
	b.Handle(tb.OnText, func(c tb.Context) error {
		// All the text messages that weren't captured by existing handlers.
		return handleURLRewriting(c, b, configmap)
	})

	// Register all command handlers
	registerSimpleCommands(b, configmap)
	registerAIHandlers(b, configmap)

	return b
}
