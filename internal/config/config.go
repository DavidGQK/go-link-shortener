package config

import (
	"flag"
)

var (
	FlagServerPort   string
	FlagShortURLBase string
)

func ParseFlags() {
	flag.StringVar(&FlagServerPort, "a", "localhost:8080", "port where server runs on")
	flag.StringVar(&FlagShortURLBase, "b", "http://localhost:8080", "base url for shortened link")
	flag.Parse()
}
