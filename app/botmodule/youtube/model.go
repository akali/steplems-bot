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

var (
	ytLinkRegex    = "(((?:https?:)?\\/\\/)?((?:www|m)\\.)?((?:youtube\\.com))(\\/(shorts\\/))([\\w\\-]+)(\\S+)?)"
	log            = logger.Factory.NewLogger("youtube")
	allowedQuality = []QualityType{HD, HD720, HD1080, MEDIUM}
	allowedType    = []VideoType{MP4, MKV, WEBM}
)

func NewModule(botApiRepo public.BotApiRepo) *YoutubeModule {
	r := regexp.MustCompile(ytLinkRegex)

	return &YoutubeModule{
		pattern:    r,
		client:     youtube.Client{},
		botApiRepo: botApiRepo,
	}
}
