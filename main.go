package main

import (
	"github.com/beaujr/plex2emby/plex"
	"github.com/beaujr/plex2emby/emby"
	"github.com/beaujr/plex2emby/process"
	"log"
)

const embyURL = "" // eg https://192.168.1.2:8096/emby
const embyAPIKEY = "" // emby api key
const embyUSERKEY = "" // emby user key

const plexURL = "" // eg https://192.168.1.2:32400
const plexAPIKEY = "" // plex api key

func main() {

	embyClient := emby.NewClient(embyURL, embyAPIKEY, embyUSERKEY)
	plexClient := plex.NewClient(plexURL, plexAPIKEY)


	client := process.NewPlex2EmbyClient(plexClient, embyClient)
	err := client.Process()
	if err != nil {
		log.Panic(err)
	}
}