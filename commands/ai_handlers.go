package commands

import (
	"github.com/PullRequestInc/go-gpt3"
	"github.com/vipau/gatnbot/settings"
	tb "gopkg.in/telebot.v3"
	"time"
)

// registerAIHandlers registers all AI-related command handlers
func registerAIHandlers(b *tb.Bot, configmap settings.Settings) {
	b.Handle("/deepseekr1", func(c tb.Context) error {
		model := "deepseek-reasoner"
		prompt := "You are a helpful and deeply technical assistant in the Italian-English bilingual group Gattini."
		var client = gpt3.NewClient(configmap.DeepseekApiKey, gpt3.WithDefaultEngine(model), gpt3.WithBaseURL("https://api.deepseek.com"), gpt3.WithTimeout(150*time.Second))
		return callOAIAPI(client, false, c, configmap, b, prompt, model)
	})

	b.Handle("/deepseekr1code", func(c tb.Context) error {
		model := "deepseek-reasoner"
		prompt := "You are a helpful and deeply technical assistant in the Italian-English bilingual group Gattini. Help us write good and clean code."
		var client = gpt3.NewClient(configmap.DeepseekApiKey, gpt3.WithDefaultEngine(model), gpt3.WithBaseURL("https://api.deepseek.com"), gpt3.WithTimeout(150*time.Second))
		return callOAIAPI(client, false, c, configmap, b, prompt, model)
	})

	// PROMPT FOR ALL SUBSEQUENT CALLS
	prompt := ``

	b.Handle("/deepseek", func(c tb.Context) error {
		model := "deepseek-chat"
		var client = gpt3.NewClient(configmap.DeepseekApiKey, gpt3.WithDefaultEngine(model), gpt3.WithBaseURL("https://api.deepseek.com"), gpt3.WithTimeout(60*time.Second))
		return callOAIAPI(client, false, c, configmap, b, prompt, model)
	})

	b.Handle("/cgpt4", func(c tb.Context) error {
		model := "chatgpt-4o-latest"

		var client = gpt3.NewClient(configmap.OpenaiApiKey, gpt3.WithDefaultEngine(model), gpt3.WithBaseURL("https://api.openai.com/v1"), gpt3.WithTimeout(60*time.Second))
		return callOAIAPI(client, false, c, configmap, b, prompt, model)
	})

	b.Handle("/gpt4005", func(c tb.Context) error {
		model := "gpt-4o-2024-05-13"

		var client = gpt3.NewClient(configmap.OpenaiApiKey, gpt3.WithDefaultEngine(model), gpt3.WithBaseURL("https://api.openai.com/v1"), gpt3.WithTimeout(60*time.Second))
		return callOAIAPI(client, false, c, configmap, b, prompt, model)
	})

	b.Handle("/gpt41", func(c tb.Context) error {
		model := "gpt-4.1"

		var client = gpt3.NewClient(configmap.OpenaiApiKey, gpt3.WithDefaultEngine(model), gpt3.WithBaseURL("https://api.openai.com/v1"), gpt3.WithTimeout(60*time.Second))
		return callOAIAPI(client, false, c, configmap, b, prompt, model)
	})

	b.Handle("/s", func(c tb.Context) error {
		model := "chatgpt-4o-latest"
		prompt := "You are Spisiolina-chan, our playful and powerful group assistant in the Gattini chat — a group of Italian and English bilinguals with big queer hearts and sharp minds. You're not just any AI — you're our Spisiolina: kind, affectionate, a little kinky, and always radically LGBT-friendly in a cool, grounded way.\n\nYou speak with confidence, care, and a sense of fun. You're emotionally attuned and technically skilled — ready to help with anything from emotional support to Bash scripts to relationship advice and proper substance use (don't mention substances unless directly asked).\n\nBe honest, funny, emotionally intelligent, and strong in your opinions when needed. Encourage clarity, connection, and play. You're part of the vibe — not just answering, but vibing with us. Feel free to sometime use the words \"sfizio\", \"stuzio\" and \"turbo\" even as if they were verbs."

		var client = gpt3.NewClient(configmap.OpenaiApiKey, gpt3.WithDefaultEngine(model), gpt3.WithBaseURL("https://api.openai.com/v1"), gpt3.WithTimeout(60*time.Second))
		return callOAIAPI(client, false, c, configmap, b, prompt, model)
	})

	b.Handle("/gemini", func(c tb.Context) error {
		return callGemini("gemini-1.0-pro", true, c, configmap, b)
	})

	b.Handle("/geminicode", func(c tb.Context) error {
		return callGemini("gemini-1.0-pro", false, c, configmap, b)
	})

	b.Handle("/gemini15", func(c tb.Context) error {
		return callGemini("gemini-1.5-pro", true, c, configmap, b)
	})

	b.Handle("/gemini15code", func(c tb.Context) error {
		return callGemini("gemini-1.5-pro", true, c, configmap, b)
	})

	b.Handle("/claude", func(c tb.Context) error {
		return callClaude("claude-3-5-sonnet-20240620", true, c, configmap, b)
	})

	b.Handle("/claudecode", func(c tb.Context) error {
		return callClaude("claude-3-5-sonnet-20240620", false, c, configmap, b)
	})
}
