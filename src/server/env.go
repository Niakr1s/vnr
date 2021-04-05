package server

import (
	"log"
	"os"
)

var isDevMode bool

func init() {
	_, dev := os.LookupEnv("DEV")
	isDevMode = dev
	log.Printf("server is running in dev mode = %v\n", isDevMode)
}
