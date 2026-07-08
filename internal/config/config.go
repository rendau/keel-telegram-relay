package config

import (
	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

var Conf = struct {
	Debug          bool   `env:"DEBUG" envDefault:"false"`
	HttpPort       string `env:"HTTP_PORT" envDefault:"80"`
	TelegramApiUrl string `env:"TELEGRAM_API_URL"` // optional, override for proxy; empty => https://api.telegram.org
	TelegramToken  string `env:"TELEGRAM_TOKEN"`
	TelegramChatId string `env:"TELEGRAM_CHAT_ID"`
}{}

func init() {
	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
