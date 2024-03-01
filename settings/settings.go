package settings

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
)

type Settings struct {
	Timezone     string  `hcl:"timezone"`
	Apiurl       string  `hcl:"apiurl"`
	Bottoken     string  `hcl:"bottoken"`
	Chatid       []int64 `hcl:"chatid"`
	Usersid      []int64 `hcl:"usersid"`
	Gpt4id       []int64 `hcl:"gpt4id"`
	Ouremail     string  `hcl:"ouremail"`
	OpenaiApikey string  `hcl:"openaiapikey"`
	Linksmsg     string  `hcl:"linksmsg"`
	GeminiApiKey string  `hcl:"geminiapikey"`
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
