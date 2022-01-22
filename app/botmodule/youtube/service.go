package youtube

import (
    "fmt"
    "io"
    "io/ioutil"
    "net/url"
    "os"
    "strings"

    tbot "github.com/go-telegram-bot-api/telegram-bot-api"
    multierror "github.com/hashicorp/go-multierror"
    "github.com/kkdai/youtube/v2"
)

func (yt *YoutubeModule) Match(text string) []string {
    return yt.pattern.FindAllString(text, -1)
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

func (yt *YoutubeModule) Download(links []string, folder string) ([]string, error) {
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

func (ytm *YoutubeModule) MessageUpdate(message *tbot.Message) error {
    links := ytm.Match(message.Text)

    if len(links) != 0 {
        loadingMsg := tbot.NewMessage(message.Chat.ID, fmt.Sprintf("loading..."))
        loadingMsg.ReplyToMessageID = message.MessageID
        if _, err := ytm.sendMessageRepo.SendMessage(loadingMsg); err != nil {
            log.Error.Println("сan not reply loading message to the message: ", err.Error())
        }

        log.Info.PrintTKV(
            "detected youtube short links of {{length}} length from {{user}}",
            "length", len(links), "user", message.From.String())

        folder, err := ioutil.TempDir("/tmp", "yt*")
        if err != nil {
            return err
        }

        defer os.RemoveAll(folder)

        filePaths, err := ytm.Download(links, folder)
        if err != nil {
            log.Error.Println(err.Error())
            // Let's try to reply to message with error message
            v := tbot.NewMessage(message.Chat.ID, fmt.Sprintf("failed to process video: %s", err.Error()))
            v.ReplyToMessageID = message.MessageID

            if _, err := ytm.sendMessageRepo.SendMessage(v); err != nil {
                log.Error.Println("failed to reply to message: ", err.Error())
            }
            return err
        }
        var filesErrs error
        for _, filePath := range filePaths {
            v := tbot.NewVideoUpload(message.Chat.ID, filePath)
            v.ReplyToMessageID = message.MessageID

            if _, err = ytm.sendMessageRepo.SendMessage(v); err != nil {
                log.Error.Println(err.Error())
                // Let's try to reply to message with error message
                v := tbot.NewMessage(message.Chat.ID, fmt.Sprintf("failed to process video: %s", err.Error()))
                v.ReplyToMessageID = message.MessageID

                if _, err := ytm.sendMessageRepo.SendMessage(v); err != nil {
                    log.Error.Println("failed to reply to message: ", err.Error())
                    filesErrs = multierror.Append(filesErrs, err)
                }
            }
        }
        return filesErrs
    }
    return nil
}
