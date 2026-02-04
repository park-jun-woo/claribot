package telegram

import (
	"context"
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Message represents an incoming telegram message
type Message struct {
	ChatID   int64
	UserID   int64
	Username string
	Text     string
}

// Handler is called when a message is received
type Handler func(msg Message)

// Bot wraps telegram bot API
type Bot struct {
	api     *tgbotapi.BotAPI
	handler Handler
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// New creates a new telegram bot
func New(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Bot{
		api:    api,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// SetHandler sets the message handler
func (b *Bot) SetHandler(h Handler) {
	b.handler = h
}

// Start begins listening for messages
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.ctx.Done():
				return
			case update := <-updates:
				if update.Message == nil {
					continue
				}
				if b.handler != nil {
					msg := Message{
						ChatID:   update.Message.Chat.ID,
						UserID:   update.Message.From.ID,
						Username: update.Message.From.UserName,
						Text:     update.Message.Text,
					}
					b.handler(msg)
				}
			}
		}
	}()

	log.Printf("Telegram bot started: @%s", b.api.Self.UserName)
	return nil
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.cancel()
	b.api.StopReceivingUpdates()
	b.wg.Wait()
	log.Println("Telegram bot stopped")
}

// Send sends a text message to a chat
func (b *Bot) Send(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}

// SendMarkdown sends a markdown formatted message
func (b *Bot) SendMarkdown(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("send markdown: %w", err)
	}
	return nil
}

// Reply sends a reply to a specific message
func (b *Bot) Reply(chatID int64, replyToID int, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = replyToID
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("reply message: %w", err)
	}
	return nil
}

// Username returns the bot's username
func (b *Bot) Username() string {
	return b.api.Self.UserName
}
