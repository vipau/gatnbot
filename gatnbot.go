package main

import (
	"github.com/vipau/gatnbot/commands"
	"github.com/vipau/gatnbot/crontasks"
	"github.com/vipau/gatnbot/settings"
	"log/slog"
)

func main() {
	// load bot token and other settings
	slog.Info("Loading settings...")
	configmap := settings.LoadSettings("settings.hcl")
	slog.Info("Loaded settings.")

	// Instantiate bot and set what commands to handle
	b := commands.HandleCommands(configmap)

	// start background activities (cron) while not blocking the flow
	crontasks.StartCronProcesses(configmap, b)

	slog.Info("Starting bot. Listening...")
	// start the bot polling (blocking call)
	b.Start()
}
