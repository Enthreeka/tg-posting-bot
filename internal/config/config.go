package config

import (
	"github.com/joho/godotenv"
	"os"
)

type (
	Config struct {
		Postgres Postgres `create_post.json:"postgres"`
		Telegram Telegram `create_post.json:"telegram"`
	}

	Postgres struct {
		URL string `create_post.json:"url"`
	}

	Telegram struct {
		Token string `create_post.json:"token"`
	}
)

func New() (*Config, error) {
	err := godotenv.Load("configs/bot.env")
	if err != nil {
		return nil, err
	}

	config := &Config{
		Postgres: Postgres{
			URL: os.Getenv("POSTGRES_URL"),
		},
		Telegram: Telegram{
			Token: os.Getenv("TOKEN_TG"),
		},
	}

	return config, nil
}
