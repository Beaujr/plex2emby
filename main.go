package main

import (
	"github.com/beaujr/plex2emby/process"
	"log"
	"github.com/beaujr/plex2emby/emby"
	"github.com/beaujr/plex2emby/plex"
	"os"
	"fmt"
)

func main() {

	embyURL, err := getEnvVar("embyURL")
	if err != nil {
		log.Panic(err)
	}
	embyAPIKEY, err := getEnvVar("embyAPIKEY")
	if err != nil {
		log.Panic(err)
	}
	embyUSERKEY, err := getEnvVar("embyUSERKEY")
	if err != nil {
		log.Panic(err)
	}
	plexURL, err := getEnvVar("plexURL")
	if err != nil {
		log.Panic(err)
	}
	plexAPIKEY, err := getEnvVar("plexAPIKEY")
	if err != nil {
		log.Panic(err)
	}

	embyClient := emby.NewClient(embyURL, embyAPIKEY, embyUSERKEY)
	plexClient := plex.NewClient(plexURL, plexAPIKEY)

	client := process.NewPlex2EmbyClient(plexClient, embyClient)
	err = client.Process()
	if err != nil {
		log.Panic(err)
	}
}

func getEnvVar(name string) (string, error) {
	v, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("%s must be set", name)
	}
	if len(v) == 0 {
		return "", fmt.Errorf("%s must not be empty", name)
	}
	return v, nil
}