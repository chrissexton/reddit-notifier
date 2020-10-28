package main

import (
	"github.com/chrissexton/reddit-notifier"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := redditnotifier.New().Execute(); err != nil {
		log.Fatal().Err(err).Msgf("could not execute notifier")
	}
}