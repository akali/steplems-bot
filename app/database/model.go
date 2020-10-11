package database

import (
	"github.com/go-bongo/bongo"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	MessagesCollection = "messages"
)

type (
	Database struct {
		*bongo.Config
		*bongo.Connection
	}
	Message struct {
		bongo.DocumentBase `bson:",inline"`
		tbot.Message
	}
)

func (d *Database) Init() error {
	var err error
	d.Connection, err = bongo.Connect(d.Config)
	return err
}

func (d *Database) SaveMessage(message *Message) error {
	return d.Connection.Collection(MessagesCollection).Save(message)
}
