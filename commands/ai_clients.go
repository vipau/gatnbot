package commands

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/google/generative-ai-go/genai"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/pkg/errors"
	"github.com/vipau/gatnbot/settings"
	"google.golang.org/api/option"
	tb "gopkg.in/telebot.v3"
	"strconv"
	"strings"
)

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

func callOAIAPI(client gpt3.Client, format bool, c tb.Context, configmap settings.Settings, b *tb.Bot, prompt string, model string) error {
	fmt.Println("OpenAI/DeepSeek -- User ID: " + strconv.FormatInt(c.Sender().ID, 10) + " | username: " + c.Sender().Username + " | full name: " + findPrintableName(c.Sender()) + " | Chat ID: " + strconv.FormatInt(c.Chat().ID, 10))
	if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
		settings.ListContainsID(configmap.OpenAICompatibleId, c.Message().Chat.ID) {
		if !c.Message().IsReply() {
			_, err := b.Reply(c.Message(), "Need to reply to a message to use this command")
			checkPrintErr(err)
		} else {
			respo, err := client.ChatCompletion(context.Background(), gpt3.ChatCompletionRequest{
				Messages: []gpt3.ChatCompletionRequestMessage{
					{
						Role:    "system",
						Content: prompt},
					{
						Role:    "user",
						Content: c.Message().ReplyTo.Text,
					},
				},
				Model: model,
			})
			if err == nil {
				if respo.Choices[0].Message.Content == "" {
					checkSendErr(errors.New("gatnbot warning: response is empty!"), b, c, true)
				} else {
					output := respo.Choices[0].Message.Content
					if format {
						// replace standard Markdown with Telegram markdown (breaks code blocks)
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
					"Gatnbot note: If the above says *\"context deadline exceeded\"*, AI took too long to generate an answer. \n"+
						"If it says *\"Service Unavailable\"* or *\"Bad gateway\"* then the API is down, try again later.")
			}
		}

	} else {
		checkSendErr(errors.New("Error: You are not authorized to use AI in this chat :(\n"+
			"Try this command in an authorized chat, or ask the admin for access"), b, c, true)
	}
	return nil
}

func callGemini(modelname string, format bool, c tb.Context, configmap settings.Settings, b *tb.Bot) error {
	fmt.Println("GMN -- User ID: " + strconv.FormatInt(c.Sender().ID, 10) + " | username: " + c.Sender().Username + " | full name: " + findPrintableName(c.Sender()) + " | Chat ID: " + strconv.FormatInt(c.Chat().ID, 10))
	if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
		settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
		if !c.Message().IsReply() {
			_, err := b.Reply(c.Message(), "Need to reply to a message to use /gemini")
			checkPrintErr(err)
		} else {
			ctx := context.Background()
			client, err := genai.NewClient(ctx, option.WithAPIKey(configmap.GeminiApiKey))
			if err != nil {
				checkSendErr(err, b, c, true)
			}
			defer client.Close()

			model := client.GenerativeModel(modelname)
			respo, err := model.GenerateContent(ctx, genai.Text(c.Message().ReplyTo.Text))

			if err == nil {
				if format {
					opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
					fixasio := strings.ReplaceAll(buildGeminiResponse(respo), "**", "TEMP_DOUBLE_ASTERISK")
					fixasio = strings.ReplaceAll(fixasio, "*", "-")
					fixasio = strings.ReplaceAll(fixasio, "TEMP_DOUBLE_ASTERISK", "*")
					_, err = b.Reply(c.Message(), fixasio, opts)
				} else {
					opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: ""}
					_, err = b.Reply(c.Message(), buildGeminiResponse(respo), opts)
				}
				if err != nil {
					checkSendErr(err, b, c, true)
				}
			} else {
				checkSendErr(err, b, c, true)
			}
		}

	}
	return nil
}

func callClaude(modelname string, format bool, c tb.Context, configmap settings.Settings, b *tb.Bot) error {
	fmt.Println("CLD-- User ID: " + strconv.FormatInt(c.Sender().ID, 10) + " | username: " + c.Sender().Username + " | full name: " + findPrintableName(c.Sender()) + " | Chat ID: " + strconv.FormatInt(c.Chat().ID, 10))
	if settings.ListContainsID(configmap.Chatid, c.Message().Chat.ID) ||
		settings.ListContainsID(configmap.Usersid, c.Message().Chat.ID) {
		if !c.Message().IsReply() {
			_, err := b.Reply(c.Message(), "Need to reply to a message to use /claude")
			checkPrintErr(err)
		} else {
			client := anthropic.NewClient(configmap.ClaudeApiKey)

			respo, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
				Model: modelname,
				Messages: []anthropic.Message{
					anthropic.NewUserTextMessage(c.Message().ReplyTo.Text),
				},
				MaxTokens: 1200,
			})

			if err == nil {
				if format {
					opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "Markdown"}
					// replace Claude Markdown with Telegram markdown (breaks code blocks)
					output := strings.ReplaceAll(*respo.Content[0].Text, "**", "TEMP_DOUBLE_ASTERISK")
					output = strings.ReplaceAll(output, "*", "_")
					output = strings.ReplaceAll(output, "TEMP_DOUBLE_ASTERISK", "*")
					_, err = b.Reply(c.Message(), output, opts)
				} else {
					opts := &tb.SendOptions{DisableWebPagePreview: true, ParseMode: ""}
					_, err = b.Reply(c.Message(), *respo.Content[0].Text, opts)
				}
				if err != nil {
					checkSendErr(err, b, c, true)
				}
			} else {
				checkSendErr(err, b, c, true)
			}
		}

	}
	return nil
}
