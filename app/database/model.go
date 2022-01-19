package database

import (
	"context"
	"time"

	"github.com/akali/steplems-bot/app/logger"
	"github.com/go-bongo/bongo"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	log = logger.Factory.NewLogger("model")
)

const (
	MessagesCollection = "messages"
)

type (
	Database interface {
		Init(updateTimeout time.Duration) (error, func())
		SaveMessage(message *Message) error
	}

	DatabaseImpl struct {
		client        *mongo.Client
		Url           string
		Database      string
		updateTimeout time.Duration
	}
	DatabaseNoOp struct{}
	Message      struct {
		bongo.DocumentBase `bson:",inline"`
		tbot.Message
	}
)

func (d *DatabaseImpl) Init(updateTimeout time.Duration) (error, func()) {
	d.updateTimeout = updateTimeout

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error

	d.client, err = mongo.NewClient(options.Client().ApplyURI(d.Url))
	if err != nil {
		return err, nil
	}
	return d.client.Connect(ctx), func() {
		if err = d.client.Disconnect(ctx); err != nil {
			log.Panic.Println(err)
		}
	}
}

func (d *DatabaseImpl) SaveMessage(message *Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.updateTimeout*time.Second)
	defer cancel()
	_, err := d.client.Database(d.Database).Collection(MessagesCollection).InsertOne(ctx, message)
	return err
}

func (d *DatabaseNoOp) Init(updateTimeout time.Duration) (error, func()) {
	return nil, nil
}

func (d *DatabaseNoOp) SaveMessage(message *Message) error {
	return nil
}

//mongodb://localhost:27019
