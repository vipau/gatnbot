package main

import (
	"github.com/prometheus/common/log"
	"github.com/vipau/gatnbot/commands"
	"github.com/vipau/gatnbot/crontasks"
	"github.com/vipau/gatnbot/settings"
)

func main() {
	// load bot token and other settings
	log.Info("Loading settings...")
	configmap := settings.LoadSettings("settings.hcl")
	log.Info("Loaded settings.")

	// Instantiate bot and set what commands to handle
	b := commands.HandleCommands(configmap)

	// start background activities (cron) while not blocking the flow
	crontasks.StartCronProcesses(configmap, b)

	log.Info("Starting bot. Listening...")
	// start the bot polling (blocking call)
	b.Start()
}
