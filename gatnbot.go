package main

import (
	"github.com/paualberto/gatnbot/commands"
	"github.com/paualberto/gatnbot/crontasks"
	"github.com/paualberto/gatnbot/settings"
)

func main() {
	// load bot token and other settings
	configmap := settings.LoadSettings("settings.hcl")

	// Tell bot what commands to handle
	b := commands.HandleCommands(configmap)

	// start background activities (cron) while non blocking the flow
	// see internal/crontasks for more info
	crontasks.StartCronProcesses(configmap, b)

	// start the bot (blocking call)
	b.Start()
}
