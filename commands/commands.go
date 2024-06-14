package commands

import (
	"bufio"
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/google/generative-ai-go/genai"
	"github.com/pkg/errors"
	"github.com/vipau/gatnbot/crontasks"
	fakernewsmod "github.com/vipau/gatnbot/fakernews-mod"
	"github.com/vipau/gatnbot/settings"
	"google.golang.org/api/option"
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

func checkSendErr(err error, b *tb.Bot, c tb.Context, isReply bool, outer ...string) {
	if err != nil {
		opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
		errmsg := "gatnbot: lol an error occurred, check it out yo...\n\n```error\n" + err.Error() + "```"
		if len(outer) != 0 {
			errmsg += "\n\n" + outer[0]
		}
		fmt.Println(errmsg)
		if isReply {
			_, err = b.Reply(c.Message(), errmsg, opts)
			checkPrintErr(err)
		} else {
			_, err = b.Send(c.Chat(), errmsg, opts)
			checkPrintErr(err)
		}
	}
}

func findPrintableName(u *tb.User) string {
	if u.LastName == "" {
		return u.FirstName
	} else {
		return u.FirstName + " " + u.LastName
	}
}

func returnFragments(path string) []string {
	//This cuts off the leading forward slash.
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	//This cuts off the trailing forward slash.
	if strings.HasSuffix(path, "/") {
		cut_off_last_char_len := len(path) - 1
		path = path[:cut_off_last_char_len]
	}
	//We need to isolate the individual components of the path.
	components := strings.Split(path, "/")
	return components
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
			fmt.Println("User ID: " + strconv.FormatInt(user.ID, 10) + " | username: " + c.Sender().Username + " | full name: " + findPrintableName(c.Sender()) + " | Chat ID: " + strconv.FormatInt(c.Chat().ID, 10))
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
						u.Host = "fxtwitter.com"
						b.Delete(c.Message())
						q := u.Query()
						if q.Has("s") {
							q.Del("s")
							u.RawQuery = q.Encode()
						}
						if q.Has("t") {
							q.Del("t")
							u.RawQuery = q.Encode()
						}
						b.Send(c.Chat(), "From: "+findPrintableName(c.Sender())+" who did not use fxtwitter... wtf\n\n"+u.String(), opts)
						checkSendErr(err, b, c, false)
					} else {
						// if it's a profile, just remove 's' and 't' trackers
						b.Delete(c.Message())
						q := u.Query()
						if q.Has("s") {
							q.Del("s")
							u.RawQuery = q.Encode()
						}
						if q.Has("t") {
							q.Del("t")
							u.RawQuery = q.Encode()
						}
						b.Send(c.Chat(), "From: "+findPrintableName(c.Sender())+" who did not remove s/t trackers from the link... wtf\n\n"+u.String(), opts)
						checkSendErr(err, b, c, false)

					}
				}

				// try for youtube
				if u.Hostname() == "youtube.com" || u.Hostname() == "www.youtube.com" || u.Hostname() == "youtu.be" {
					q := u.Query()
					if q.Has("si") {
						b.Delete(c.Message())
						q.Del("si")
						u.RawQuery = q.Encode()
						b.Send(c.Chat(), "From: "+findPrintableName(c.Sender())+" who did not remove the 'si' tracking tag... wtf\n\n"+u.String(), opts)
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
			if err != nil {
				checkSendErr(err, b, c, false)
			} else {
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					checkPrintErr(err)
				}(resp.Body)
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					checkSendErr(err, b, c, false)
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

	b.Handle("/gpt4", func(c tb.Context) error {
		return callGPT4(true, c, configmap, b)
	})

	b.Handle("/gpt4code", func(c tb.Context) error {
		return callGPT4(false, c, configmap, b)
	})

	b.Handle("/gemini", func(c tb.Context) error {
		modelname := "gemini-pro"
		if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
			settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
			if !c.Message().IsReply() {
				_, err = b.Reply(c.Message(), "Need to reply to a message to use /gemini")
				checkPrintErr(err)
			} else {
				ctx := context.Background()
				client, err := genai.NewClient(ctx, option.WithAPIKey(configmap.GeminiApiKey))
				if err != nil {
					checkSendErr(err, b, c, true)
				}
				defer client.Close()

				model := client.GenerativeModel(modelname)
				resp, err := model.GenerateContent(ctx, genai.Text(c.Message().ReplyTo.Text))

				if err == nil {
					opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
					fixasio := strings.ReplaceAll(buildGeminiResponse(resp), "**", "TEMP_DOUBLE_ASTERISK")
					fixasio = strings.ReplaceAll(fixasio, "*", "_")
					fixasio = strings.ReplaceAll(fixasio, "TEMP_DOUBLE_ASTERISK", "*")
					_, err = b.Reply(c.Message(), fixasio, opts)
					if err != nil {
						checkSendErr(err, b, c, true)
					}
				} else {
					checkSendErr(err, b, c, true)
				}
			}

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
		_, err = b.Send(c.Message().Chat, crontasks.Viernes)
		return err
	})
	b.Handle("/sabado", func(c tb.Context) error {
		_, err = b.Send(c.Message().Chat, crontasks.Sabado)
		return err
	})

	return b
}

func buildGeminiResponse(resp *genai.GenerateContentResponse) string {
	var output strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				output.WriteString(fmt.Sprintf("%v", part))
			}
		}
	}
	return output.String()
}

func callGPT4(format bool, c tb.Context, configmap settings.Settings, b *tb.Bot) error {
	model := "gpt-4o"
	if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
		settings.ListContainsID(configmap.Gpt4id, c.Message().Chat.ID) {
		if !c.Message().IsReply() {
			_, err := b.Reply(c.Message(), "Need to reply to a message to use /gpt4")
			checkPrintErr(err)
		} else {
			client := gpt3.NewClient(configmap.OpenaiApikey, gpt3.WithDefaultEngine(model))
			if len(c.Message().ReplyTo.Text) > 1024 {
				_, err := b.Reply(c.Message(), "Gatnbot warning: Prompt too long, sorry bro")
				checkPrintErr(err)
			} else {
				resp, err := client.ChatCompletion(context.Background(), gpt3.ChatCompletionRequest{
					Messages: []gpt3.ChatCompletionRequestMessage{
						{
							Role:    "system",
							Content: "You are GattiniBot, a bot in a group of people called Gattini.",
						},
						{
							Role:    "user",
							Content: c.Message().ReplyTo.Text,
						},
					},
					Model: model,
				})
				if err == nil {
					if resp.Choices[0].Message.Content == "" {
						checkSendErr(errors.New("gatnbot warning: response is empty!"), b, c, true)
					} else {
						output := resp.Choices[0].Message.Content
						if format {
							// replace GPT4 Markdown with Telegram markdown (breaks code blocks)
							output = strings.ReplaceAll(output, "**", "TEMP_DOUBLE_ASTERISK")
							output = strings.ReplaceAll(output, "*", "_")
							output = strings.ReplaceAll(output, "TEMP_DOUBLE_ASTERISK", "*")
						}
						opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
						_, err = b.Reply(c.Message(), output, opts)
						if err != nil {
							checkSendErr(err, b, c, true)
						}
					}
				} else {
					checkSendErr(err, b, c, true,
						"Gatnbot note: If the above says *\"context deadline exceeded\"*, GPT took too long to generate an answer. Please try a simpler prompt, try again later, or if this is important try /gpt4 \n"+
							"If it says *\"Service Unavailable\"* or *\"Bad gateway\"* then the API is down, try again later.")
				}
			}
		}

	} else {
		checkSendErr(errors.New("Error: You are not authorized to use GPT4 in this chat :(\n"+
			"Try /gpt3 here, or ask the admin for access to GPT4"), b, c, true)
	}
	return nil
}
