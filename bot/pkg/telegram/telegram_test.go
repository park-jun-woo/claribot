package telegram

import (
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
}

func loadToken() string {
	home, _ := os.UserHomeDir()
	data, err := os.ReadFile(home + "/.claribot/config.yaml")
	if err != nil {
		return ""
	}
	var cfg Config
	yaml.Unmarshal(data, &cfg)
	return cfg.Telegram.Token
}

func TestBot(t *testing.T) {
	token := loadToken()
	if token == "" || token == "BOT_TOKEN" {
		t.Skip("Set real token in ~/.claribot/config.yaml")
	}

	bot, err := New(token)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	t.Logf("Bot created: @%s", bot.Username())

	// Set handler
	received := make(chan Message, 1)
	bot.SetHandler(func(msg Message) {
		t.Logf("Received: [%s] %s", msg.Username, msg.Text)
		received <- msg
	})

	// Start bot
	if err := bot.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer bot.Stop()

	t.Log("Bot running. Send a message to the bot within 30 seconds...")

	// Wait for message or timeout
	select {
	case msg := <-received:
		// Echo back
		err := bot.Send(msg.ChatID, "Echo: "+msg.Text)
		if err != nil {
			t.Errorf("Failed to send: %v", err)
		} else {
			t.Log("Reply sent successfully")
		}
	case <-time.After(30 * time.Second):
		t.Log("Timeout - no message received (this is OK for manual test)")
	}
}
