package youtube

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/kkdai/youtube/v2"
)

func (yt *YoutubeModule) Match(text string) []string {
	return yt.pattern.FindAllString(text, -1)
}

func (yt *YoutubeModule) retryIfErr(f func() error) {
	for i := 1; i <= RETRY_TIMES; i++ {
		if err := f(); err != nil {
			time.Sleep(time.Second * time.Duration(2<<i))
		} else {
			return
		}
	}
}

func (yt *YoutubeModule) chooseFormat(formats youtube.FormatList) *youtube.Format {
	formats = formats.WithAudioChannels()

	for i := range formats {
		for _, q := range allowedQuality {
			for _, t := range allowedType {
				if (formats[i].Quality == string(q) || formats[i].QualityLabel == string(q)) &&
					strings.Contains(formats[i].MimeType, string(t)) {

					return &formats[i]
				}
			}
		}
	}

	return &formats[0]
}

func (yt *YoutubeModule) downloadPerLinkBackedOff(
	v *youtube.Video,
	format *youtube.Format,
	folder string,
) (s string, err error) {
	yt.retryIfErr(func() error {
		s, err = yt.downloadPerLink(v, format, folder)

		return err
	})

	return
}

func (yt *YoutubeModule) downloadPerLink(
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

func (m *YoutubeMessage) SanitizeTitle() {
	m.Title = strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						strings.Replace(
							strings.Replace(m.Title, "*", "\\*", -1),
							"_", "\\_", -1),
						"~", "\\~", -1),
					"`", "\\`", -1),
				"|", "\\|", -1),
			"[", "\\[", -1),
		"]", "\\]", -1)

}

func (m *YoutubeMessage) FormCaption() string {
	m.SanitizeTitle()

	b := strings.Builder{}

	b.WriteRune('*')
	b.WriteString(m.Title)
	b.WriteRune('*')

	b.WriteString("\n\n[link | сілтеме]")
	b.WriteRune('(')
	b.WriteString(m.Link)
	b.WriteRune(')')
	return b.String()
}

func (yt *YoutubeModule) Download(links []string, folder string) ([]YoutubeMessage, error) {
	msgs := make([]YoutubeMessage, 0)

	for _, l := range links {
		v, err := yt.client.GetVideo(l)
		if err != nil {
			log.Note.PrintTKV("can't get metadata for link: {{error}}", "error", err)

			return nil, err
		}

		chosenFormat := yt.chooseFormat(v.Formats)
		filename, err := yt.downloadPerLinkBackedOff(v, chosenFormat, folder)
		if err != nil {
			return nil, err
		}

		// there has to be ffmpeg stuff

		msgs = append(msgs, YoutubeMessage{
			Link:  l,
			Title: v.Title,
			Path:  filename})
	}

	return msgs, nil
}

func (ytm *YoutubeModule) MessageUpdate(message *tbot.Message) error {
	links := ytm.Match(message.Text)
	matches := len(links) > 0

	log.Info.Println("MessageUpdate ", message.Text, " ", links, " ", strings.HasPrefix(message.Text, downloadCommand))

	if !matches {
		// Nothing to do here.
		return nil
	}

	if strings.HasPrefix(message.Text, downloadCommand) {
		if err := ytm.botApiRepo.DeleteMessage(message.Chat.ID, message.MessageID); err != nil {
			log.Error.PrintT("failed to remove message in Chat ", message.Chat, " with links ", links, ", error: ", err.Error())
		}
		return ytm.download(message, links)
	}

	if ytm.askToDownload {
		promptMsgCfg := tbot.NewMessage(message.Chat.ID, "wanna download it?")
		promptMsgCfg.ReplyToMessageID = message.MessageID

		promptMsgCfg.ReplyMarkup = tbot.NewInlineKeyboardMarkup(
			tbot.NewInlineKeyboardRow(
				tbot.NewInlineKeyboardButtonData("Yes", fmt.Sprintf("%s %s", downloadCommand, links[0])),
				tbot.NewInlineKeyboardButtonData("No", "-"),
			),
		)
		if _, err := ytm.botApiRepo.SendMessage(promptMsgCfg); err != nil {
			log.Error.Println("can not reply prompt message to the message: ", err.Error())
			return err
		}
		return nil
	}
	return ytm.download(message, links)
}

func (ytm *YoutubeModule) download(message *tbot.Message, links []string) error {
	loadingMsgCfg := tbot.NewMessage(message.Chat.ID, "loading...")
	loadingMsgCfg.ReplyToMessageID = message.MessageID
	var loadingMsg *tbot.Message = nil
	if newLoadingMsg, err := ytm.botApiRepo.SendMessage(loadingMsgCfg); err != nil {
		log.Error.Println("сan not reply loading message to the message: ", err.Error())
	} else {
		loadingMsg = &newLoadingMsg
	}

	defer func() {
		if err := ytm.botApiRepo.DeleteMessage(loadingMsg.Chat.ID, loadingMsg.MessageID); err != nil {
			log.Error.PrintT("failed to remove message in Chat ", loadingMsg.Chat, " with links ", links, ", error: ", err.Error())
		}
	}()

	log.Info.PrintTKV(
		"detected youtube short links of {{length}} length from {{user}}",
		"length", len(links), "user", message.From.String())

	folder, err := ioutil.TempDir("/tmp", "yt*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(folder)

	yms, err := ytm.Download(links, folder)
	if err != nil {
		log.Error.Println(err.Error())
		// Let's try to reply to message with error message
		v := tbot.NewMessage(message.Chat.ID, fmt.Sprintf("failed to process video: %s", err.Error()))
		v.ReplyToMessageID = message.MessageID

		if _, err := ytm.botApiRepo.SendMessage(v); err != nil {
			log.Error.Println("failed to reply to message: ", err.Error())
		}
		return err
	}
	var filesErrs error
	for _, ym := range yms {
		v := tbot.NewVideoUpload(message.Chat.ID, ym.Path)
		v.Caption = ym.FormCaption()
		v.ParseMode = tbot.ModeMarkdown
		v.ReplyToMessageID = message.MessageID

		if _, err = ytm.botApiRepo.SendMessage(v); err != nil {
			log.Error.Println(err.Error())
			// Let's try to reply to message with error message
			v := tbot.NewMessage(message.Chat.ID, fmt.Sprintf("failed to process video: %s", err.Error()))
			v.ReplyToMessageID = message.MessageID

			if _, err := ytm.botApiRepo.SendMessage(v); err != nil {
				log.Error.Println("failed to reply to message: ", err.Error())
				if err := ytm.botApiRepo.DeleteMessage(loadingMsg.Chat.ID, loadingMsg.MessageID); err != nil {
					log.Error.Println(err.Error())
				}
				filesErrs = multierror.Append(filesErrs, err)
			}
		}
	}
	return filesErrs
}