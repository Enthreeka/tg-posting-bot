package main

import (
	"context"
	"github.com/Enthreeka/tg-posting-bot/internal/bot"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	ch := make(chan os.Signal, 1)
	go func() {
		sig := <-ch
		log.Printf("handle signal %s", sig)
		cancel()
	}()

	tgBot := bot.NewBot()
	tgBot.Run(ctx)
}
