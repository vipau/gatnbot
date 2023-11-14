package commands

import (
	"bufio"
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/vipau/gatnbot/crontasks"
	fakernewsmod "github.com/vipau/gatnbot/fakernews-mod"
	"github.com/vipau/gatnbot/settings"
	tb "gopkg.in/telebot.v3"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func checkPrintErr(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}

func findPrintableName(c tb.Context) string {
	if c.Sender().Username == "" {
		return c.Sender().FirstName + " " + c.Sender().LastName
	} else {
		return c.Sender().Username
	}
}

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

	b.Handle(tb.OnText, func(c tb.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.

		// Print user ID and username on terminal, if message doesn't come from group
		var user = c.Sender()
		if !settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) {
			fmt.Println("User ID: " + strconv.FormatInt(user.ID, 10) + " username: " + findPrintableName(c))
		}

		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {

			// Detect if message is a link
			msg := strings.TrimSpace(c.Message().Text)
			u, err := url.Parse(msg)
			if err != nil {
				return nil
			} else {
				// it's a link

				// send link with the telegram preview
				opts := &tb.SendOptions{DisableWebPagePreview: false}

				// try for instagram
				if u.Hostname() == "instagram.com" || u.Hostname() == "www.instagram.com" {
					u.Host = "ddinstagram.com"
					b.Delete(c.Message())
					q := u.Query()
					q.Del("igshid")
					u.RawQuery = q.Encode()
					b.Send(c.Chat(), "From: "+c.Sender().FirstName+" "+c.Sender().LastName+" who did not use ddinstagram and/or remove the 'igshid' tracking tag... wtf\n\n"+u.String(), opts)
				}

				// try for twitter
				if u.Hostname() == "twitter.com" || u.Hostname() == "www.twitter.com" {
					u.Host = "fxtwitter.com"
					b.Delete(c.Message())
					b.Send(c.Chat(), "From: "+c.Sender().FirstName+" "+c.Sender().LastName+" who did not use fxtwitter... wtf\n\n"+u.String(), opts)
				}

				// try for youtube
				if u.Hostname() == "youtube.com" || u.Hostname() == "www.youtube.com" || u.Hostname() == "youtu.be" {
					q := u.Query()
					if q.Has("si") {
						b.Delete(c.Message())
						q.Del("si")
						u.RawQuery = q.Encode()
						b.Send(c.Chat(), "From: "+c.Sender().FirstName+" "+c.Sender().LastName+" who did not use remove the 'si' tracking tag... wtf\n\n"+u.String(), opts)
					}
				}
			}
		}
		return nil
	})

	b.Handle("/myid", func(c tb.Context) error {
		opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
		_, err = b.Send(c.Message().Chat, "Your ID: `"+strconv.FormatInt(c.Sender().ID, 10)+"`", opts)
		return err
	})

	b.Handle("/links", func(c tb.Context) error {
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
			opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
			_, err = b.Send(c.Message().Chat, configmap.Linksmsg, opts)
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
			_, err = b.Send(c.Message().Chat, fmt.Sprintf("this chat is now cringe-protected for %d minutes thanks the power of TURBO", rando))
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
			_, err = b.Send(c.Message().Chat, out)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		}
		return nil
	})

	b.Handle("/supercazzola", func(c tb.Context) error {
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
			// query the BS generator
			resp, err := http.Get("http://ftrv.se/bullshit")
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				checkPrintErr(err)
			}(resp.Body)
			if err != nil {
				opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
				errmsg := "lol an error occurred\ncheck it out bro\n\n```error\n" + err.Error() + "```"
				fmt.Println(errmsg)
				_, err = b.Send(c.Message().Chat, errmsg, opts)
				checkPrintErr(err)
			} else {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
					errmsg := "lol an error occurred\ncheck it out bro\n\n```error\n" + err.Error() + "```"
					fmt.Println(errmsg)
					_, err = b.Send(c.Message().Chat, errmsg, opts)
					checkPrintErr(err)
				} else {
					// here we enter a loop to strip the HTML tags from the response
					scanner := bufio.NewScanner(strings.NewReader(string(body)))
					for scanner.Scan() {
						line := scanner.Text()
						// simply check that the line does not start with <
						if !strings.HasPrefix(line, "<") {
							_, err = b.Send(c.Message().Chat, line)
							checkPrintErr(err)
						}
					}
				}

			}
		}
		return nil
	})

	b.Handle("/gpt3", func(c tb.Context) error {
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
			if !c.Message().IsReply() {
				_, err = b.Reply(c.Message(), "Need to reply to a message to use /gpt3")
				checkPrintErr(err)
			} else {
				client := gpt3.NewClient(configmap.OpenaiApikey, gpt3.WithDefaultEngine("gpt-3.5-turbo"), gpt3.WithTimeout(45*time.Second))
				if len(c.Message().ReplyTo.Text) > 2048 {
					_, err = b.Reply(c.Message(), "Gatnbot warning: Prompt too long, sorry bro")
					checkPrintErr(err)
				} else {
					resp, err := client.ChatCompletion(context.Background(), gpt3.ChatCompletionRequest{
						Messages: []gpt3.ChatCompletionRequestMessage{
							{
								Role:    "system",
								Content: "",
							},
							{
								Role:    "user",
								Content: c.Message().ReplyTo.Text,
							},
						},
						//					Functions:	  nil,
						Model: "gpt-3.5-turbo",
						//						MaxTokens: 512,
						//Stop:      []string{"."},
						//					Temperature:      gpt3.Float32Ptr(0.7),
						//					TopP:             gpt3.Float32Ptr(1),
						//					N:                gpt3.Float32Ptr(1),
						//					PresencePenalty:  0,
						//					FrequencyPenalty: 0,
					})
					if err == nil {
						if resp.Choices[0].Message.Content == "" {
							_, err = b.Reply(c.Message(), "gatnbot warning: response is empty!")
							checkPrintErr(err)
						} else {
							_, err = b.Reply(c.Message(), resp.Choices[0].Message.Content)
							if err != nil {
								opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
								_, err2 := b.Reply(c.Message(), "gatnbot error: \n\n```error\n"+err.Error()+"\n```", opts)
								checkPrintErr(err2)
							}
						}
					} else {
						opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
						_, err = b.Reply(c.Message(), "Gatnbot: error occurred :(( details:\n\n```go\n"+err.Error()+
							"```\n\nGatnbot note: If the above says *\"context deadline exceeded\"*, GPT took too long to generate an answer. Please try a simpler prompt or try again later. \n"+
							"If it says *\"Service Unavailable\"* or *\"Bad gateway\"* then the API is down, try again later.", opts)
						checkPrintErr(err)
					}
				}
			}

		}
		return nil
	})

	b.Handle("/gpt4", func(c tb.Context) error {
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Gpt4id, c.Message().Chat.ID) {
			if !c.Message().IsReply() {
				_, err = b.Reply(c.Message(), "Need to reply to a message to use /gpt4")
				checkPrintErr(err)
			} else {
				client := gpt3.NewClient(configmap.OpenaiApikey, gpt3.WithDefaultEngine("gpt-4-1106-preview"))
				if len(c.Message().ReplyTo.Text) > 512 {
					_, err = b.Reply(c.Message(), "Gatnbot warning: Prompt too long, sorry bro")
					checkPrintErr(err)
				} else {
					resp, err := client.ChatCompletion(context.Background(), gpt3.ChatCompletionRequest{
						Messages: []gpt3.ChatCompletionRequestMessage{
							{
								Role: "system",
								Content: "You are GattiniBot, a bot in a group of people called Gattini. Be the most helpful but concise." +
									//"Output Markdown if needed, but using single * for *bold* and single _ for _italics_.",
									" Output simple HTML. If formatting is needed, you can make use of the HTML tags a, b, i, s, u, code (for monospace text)." +
									" Do NOT use ANY other tag or your message will not go through.",
							},
							{
								Role:    "user",
								Content: c.Message().ReplyTo.Text,
							},
						},
						//					Functions:	  nil,
						Model: "gpt-4-1106-preview",
						//						MaxTokens: 96,
						//Stop:      []string{"."},
						//					Temperature:      gpt3.Float32Ptr(0.7),
						//					TopP:             gpt3.Float32Ptr(1),
						//					N:                gpt3.Float32Ptr(1),
						//					PresencePenalty:  0,
						//					FrequencyPenalty: 0,
					})
					if err == nil {
						if resp.Choices[0].Message.Content == "" {
							_, err = b.Reply(c.Message(), "gatnbot warning: response is empty!")
							checkPrintErr(err)
						} else {
							opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "HTML"}
							_, err = b.Reply(c.Message(), resp.Choices[0].Message.Content, opts)
							if err != nil {
								optsMd := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
								_, err2 := b.Reply(c.Message(), "gatnbot error: \n\n```error\n"+err.Error()+"\n```", optsMd)
								checkPrintErr(err2)
							}
						}
					} else {
						opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
						_, err = b.Reply(c.Message(), "Gatnbot: error occurred :(( details:\n\n```go\n"+err.Error()+
							"```\n\nGatnbot note: If the above says *\"context deadline exceeded\"*, GPT took too long to generate an answer. Please try a simpler prompt or try again later. \n"+
							"If it says *\"Service Unavailable\"* or *\"Bad gateway\"* then the API is down, try again later.", opts)
						checkPrintErr(err)
					}
				}
			}

		} else {
			_, err = b.Reply(c.Message(), "Error: You are not authorized to use GPT4 in this chat :(\n"+
				"Try /gpt3 here, or ask the admin for access to GPT4")
			checkPrintErr(err)
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
				opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
				_, err = b.Send(c.Message().Chat, "Error occurred while playing "+gladosLine+" :( details: \n\n```error\n"+err.Error()+"```", opts)
				checkPrintErr(err)
			}
		}
		return nil
	})

	// manual viernes/sabado invocation
	b.Handle("/viernes", func(c tb.Context) error {
		_, err = b.Send(c.Message().Chat, crontasks.Viernes)
		return err
	})
	b.Handle("/sabado", func(c tb.Context) error {
		_, err = b.Send(c.Message().Chat, crontasks.Sabado)
		return err
	})

	return b
}
