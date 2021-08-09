package settings

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
)

type Settings struct {
	Timezone string  `hcl:"timezone"`
	Apiurl   string  `hcl:"apiurl"`
	Bottoken string  `hcl:"bottoken"`
	Chatid   []int64 `hcl:"chatid"`
	Adminid  []int64 `hcl:"adminid"`
	Ouremail string  `hcl:"ouremail"`
	Linksmsg string  `hcl:"linksmsg"`
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

// Has checks if an array contains a specific value
func Has(list []int64, a int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
