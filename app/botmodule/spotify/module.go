package spotify

import (
	"strings"

	"github.com/akali/steplems-bot/app/bot/public"
	"github.com/akali/steplems-bot/app/logger"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	NOW_PLAYING_COMMAND = "/nowplaying"
)

type Client interface {
	Authorize(string) error
	GetNowPlaying(string) error
}

type SpotifyModule struct {
	client     Client
	botApiRepo public.BotApiRepo
}

func NewModule(botApiRepo public.BotApiRepo, clientId, clientSecret string) *SpotifyModule {
	return &SpotifyModule{
		botApiRepo: botApiRepo,
		client:     nil,
	}
}

var (
	log = logger.Factory.NewLogger("spotify")
)

func (s *SpotifyModule) MessageUpdate(message *tbot.Message) error {
	if strings.HasPrefix(message.Text, NOW_PLAYING_COMMAND) {
		_, err := s.botApiRepo.SendMessage(tbot.NewMessage(message.Chat.ID, "Not implemented yet :("))
		return err
	}
	return nil
}
