package youtube

import (
	"github.com/akali/steplems-bot/app/logger"
	"github.com/kkdai/youtube/v2"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type (
	Youtube struct {
		pattern *regexp.Regexp
		client  youtube.Client
	}

	VideoType   string
	QualityType string
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

func NewYoutube() *Youtube {
	r := regexp.MustCompile(ytLinkRegex)

	return &Youtube{
		pattern: r,
		client:  youtube.Client{},
	}
}

func (yt *Youtube) Match(text string) []string {
	return yt.pattern.FindAllString(text, -1)
}

func (yt *Youtube) chooseFormat(formats youtube.FormatList) *youtube.Format {
	formats = formats.WithAudioChannels()

	isFailedToFind := true
	var chosenFormat *youtube.Format

	for i := range formats {
		for _, q := range allowedQuality {
			for _, t := range allowedType {
				if (formats[i].Quality == string(q) || formats[i].QualityLabel == string(q)) &&
					strings.Contains(formats[i].MimeType, string(t)) {
					isFailedToFind = false
					chosenFormat = &formats[i]

					break
				}
			}
		}
	}
	if isFailedToFind {
		chosenFormat = &formats[0]
	}

	return chosenFormat
}

func (yt *Youtube) downloadPerLink(
	v *youtube.Video,
	format *youtube.Format,
	folder string,
) (string, error) {
	stream, _, err := yt.client.GetStream(v, format)
	if err != nil {
		return "", err
	}

	a := strings.Split(format.MimeType, "/")
	fileExt := strings.Split(a[1], ";")[0]

	filename := folder + "/" + url.PathEscape(v.ID) + "." + fileExt

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, stream)
	if err != nil {
		return "", err
	}

	log.Succ.PrintTKV(
		"downloaded short by id {{id}} and saved it into {{path}}",
		"id", v.ID, "path", filename)

	return filename, nil
}

func (yt *Youtube) Download(links []string, folder string) ([]string, error) {
	paths := make([]string, 0)
	for _, l := range links {
		v, err := yt.client.GetVideo(l)
		if err != nil {
			log.Note.PrintTKV("can't get metadata for link: {{error}}", "error", err)

			return nil, err
		}

		chosenFormat := yt.chooseFormat(v.Formats)
		filename, err := yt.downloadPerLink(v, chosenFormat, folder)
		if err != nil {
			return nil, err
		}

		// there has to be ffmpeg stuff

		paths = append(paths, filename)
	}

	return paths, nil
}
