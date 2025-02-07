package settings

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
)

type Settings struct {
	Timezone       string  `hcl:"timezone"`
	Apiurl         string  `hcl:"apiurl"`
	Bottoken       string  `hcl:"bottoken"`
	Chatid         []int64 `hcl:"chatid"`
	Usersid        []int64 `hcl:"usersid"`
	Deepseekid     []int64 `hcl:"deepseekid"`
	Ouremail       string  `hcl:"ouremail"`
	DeepseekApiKey string  `hcl:"deepseekapikey"`
	Linksmsg       string  `hcl:"linksmsg"`
	GeminiApiKey   string  `hcl:"geminiapikey"`
	ClaudeApiKey   string  `hcl:"claudeapikey"`
}

// LoadSettings unmarshals the HCL config file and returns our Settings.
func LoadSettings(filename string) Settings {
	var Config Settings
	err := hclsimple.DecodeFile(filename, nil, &Config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	return Config
}

// ListContainsID checks if an array contains a specific value
func ListContainsID(list []int64, a int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
