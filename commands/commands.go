package commands

import (
	"bufio"
	"fmt"
	fakernews_mod "github.com/paualberto/gatnbot/fakernews-mod"
	"github.com/paualberto/gatnbot/settings"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

	b.Handle("/links", func(m *tb.Message) {
		if settings.ListContainsID(configmap.Chatid, m.Chat.ID) ||
			settings.ListContainsID(configmap.Adminid, m.Chat.ID) {
			opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
			b.Send(m.Chat, configmap.Linksmsg, opts)
		}
	})

	b.Handle("/turbo", func(m *tb.Message) {
		if settings.ListContainsID(configmap.Chatid, m.Chat.ID) ||
			settings.ListContainsID(configmap.Adminid, m.Chat.ID) {
			rand.Seed(time.Now().UnixNano())
			min := 4
			max := 57
			rando := rand.Intn(max-min+1) + min
			b.Send(m.Chat, fmt.Sprintf("this chat is now cringe-protected for %d minutes thanks the power of TURBO", rando))
		}
	})

	b.Handle("/hackernews", func(m *tb.Message) {
		if settings.ListContainsID(configmap.Chatid, m.Chat.ID) ||
			settings.ListContainsID(configmap.Adminid, m.Chat.ID) {
			if _, err := os.Stat("model.json"); err == nil {
			} else if os.IsNotExist(err) {
				fakernews_mod.TrainModel()
			} // train the model first if it doesn't exist

			out := fakernews_mod.GenerateNews()
			b.Send(m.Chat, string(out))
		}
	})

	b.Handle("/supercazzola", func(m *tb.Message) {
		if settings.ListContainsID(configmap.Chatid, m.Chat.ID) ||
			settings.ListContainsID(configmap.Adminid, m.Chat.ID) {
			// query the BS generator
			resp, err := http.Get("http://ftrv.se/bullshit")
			defer resp.Body.Close()
			if err != nil {
				errmsg := "lol an error occurred\ncheck it out bro\n\n" + err.Error()
				fmt.Println(errmsg)
				b.Send(m.Chat, errmsg)
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					errmsg := "lol an error occurred\ncheck it out bro\n\n" + err.Error()
					fmt.Println(errmsg)
					b.Send(m.Chat, errmsg)
				} else {
					// here we enter a loop to strip the HTML tags from the response
					scanner := bufio.NewScanner(strings.NewReader(string(body)))
					for scanner.Scan() {
						line := scanner.Text()
						// simply check that the line does not start with <
						if !strings.HasPrefix(line, "<") {
							b.Send(m.Chat, line)
						}
					}
				}

			}
		}
	})

	return b
}
