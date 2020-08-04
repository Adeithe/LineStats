package manager

import (
	"LineStats/bot/discordbot"
	"LineStats/bot/twitchbot"
)

type BotManager struct {
	Twitch  *twitchbot.TwitchBot
	Discord *discordbot.DiscordBot
}

func New() *BotManager {
	return &BotManager{
		Twitch:  twitchbot.New(),
		Discord: discordbot.New(),
	}
}

func (mgr *BotManager) Connect() error {
	if err := mgr.Twitch.Start(); err != nil {
		return err
	}
	return mgr.Discord.Start()
}

func (mgr *BotManager) Close() {
	mgr.Twitch.Close()
	mgr.Discord.Close()
}
