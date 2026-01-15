package commands

import "time"

// File paths
const (
	ModelFilePath = "model.json"
	GladosDir     = "glados"
)

// AI API Configuration
const (
	DeepSeekBaseURL = "https://api.deepseek.com"
	OpenAIBaseURL   = "https://api.openai.com/v1"
)

// AI Models
const (
	ModelDeepSeekReasoner   = "deepseek-reasoner"
	ModelDeepSeekChat       = "deepseek-chat"
	ModelChatGPT4oLatest    = "chatgpt-4o-latest"
	ModelGPT4o20240513      = "gpt-4o-2024-05-13"
	ModelGPT41              = "gpt-4.1"
	ModelGemini10Pro        = "gemini-1.0-pro"
	ModelGemini15Pro        = "gemini-1.5-pro"
	ModelClaude35Sonnet     = "claude-3-5-sonnet-20240620"
)

// Timeouts
const (
	DefaultAITimeout  = 60 * time.Second
	ReasonerTimeout   = 150 * time.Second
)

// Turbo command range
const (
	TurboMinMinutes = 4
	TurboMaxMinutes = 57
)
