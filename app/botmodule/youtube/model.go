package youtube

import (
	"regexp"

	"github.com/akali/steplems-bot/app/bot/public"
	"github.com/akali/steplems-bot/app/logger"
	"github.com/kkdai/youtube/v2"
)

type (
	YoutubeModule struct {
		pattern    *regexp.Regexp
		client     youtube.Client
		botApiRepo public.BotApiRepo
		askToDownload bool
	}

	YoutubeMessage struct {
		Title, Link, Path string
	}

	VideoType   string
	QualityType string
)

const (
	RETRY_TIMES = 5
)

const (
	HD     QualityType = "hd"
	HD720  QualityType = "hd720"
	HD1080 QualityType = "hd1080"
	MEDIUM QualityType = "medium"
)

const (
	MP4  VideoType = "video/mp4"
	MKV  VideoType = "video/mkv"
	WEBM VideoType = "video/webm"
)

const (
	downloadCommand = "/ytdl" 
)

var (
	ytShortsLinkRegex    = "(((?:https?:)?\\/\\/)?((?:www|m)\\.)?((?:youtube\\.com))(\\/(shorts\\/))([\\w\\-]+)(\\S+)?)"
	ytLinkRegex    = `(?m)(((?:https?:)?//)?((?:www|m).)?(((?:youtube\.com)/watch\?v=([\w]+)*))|(((?:https?:)?//)?((?:www|m).)?((?:youtu\.be/)([\w]+)*)))`
	log            = logger.Factory.NewLogger("youtube")
	allowedQuality = []QualityType{HD, HD720, HD1080, MEDIUM}
	allowedType    = []VideoType{MP4, MKV, WEBM}

	rShorts = regexp.MustCompile(ytShortsLinkRegex)
	rFull = regexp.MustCompile(ytLinkRegex)
)

func NewModule(botApiRepo public.BotApiRepo) *YoutubeModule {
	return &YoutubeModule{
		pattern:    rShorts,
		client:     youtube.Client{},
		botApiRepo: botApiRepo,
	}
}

func NewModuleFull(botApiRepo public.BotApiRepo) *YoutubeModule {
	return &YoutubeModule{
		pattern:    rFull,
		client:     youtube.Client{},
		botApiRepo: botApiRepo,
		askToDownload: true,
	}
}
