package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rikkuness/discord-rpc"
	"github.com/shkh/lastfm-go/lastfm"
)

func main() {
	apiKey := ""
	apiSecret := ""
	for _, env := range os.Environ() {
		name,value,_ := strings.Cut(env, "=")
		switch name {
			case "SCROBBLECORD_LASTFM_API_KEY":
				apiKey = value
			case "SCROBBLECORD_LASTFM_API_SECRET":
				apiSecret = value
		}
	}

	if apiKey == "" || apiSecret == "" {
		fmt.Println("scrobblecord failed!")
		fmt.Println("please make sure you've set up your env variables:")
		fmt.Println("put Last.fm API key in SCROBBLECORD_LASTFM_API_KEY and API secret in SCROBBLECORD_LASTFM_API_SECRET")
		return
	}

	lastfmApi := lastfm.New(apiKey, apiSecret)


	var rpc *discordrpc.Client = nil

	for {
		time.Sleep(time.Second * 2)
		
		recentTracksArgs := make(map[string]any)
		recentTracksArgs["user"] = "skill_issue_dev"
		recentTracksArgs["api_key"] = apiKey
		recentTracksArgs["limit"] = 1
		recentTracks, err := lastfmApi.User.GetRecentTracks(recentTracksArgs)
		if err != nil {
			fmt.Println("failed to GetRecentTracks of the user!")
			fmt.Println("err = " + err.Error())
		}

		if len(recentTracks.Tracks) == 0 {
			fmt.Println("can't get any tracks. skipping.")
			continue
		}

		lastTrack := recentTracks.Tracks[0]
		lastTrackNowPlaying := lastTrack.NowPlaying

		if lastTrackNowPlaying == "" {
			if rpc != nil {
				err = rpc.Socket.Close()
				if err != nil {
					fmt.Println("scrobblecord failed! failed to close the Discord RPC socket")
					fmt.Println("err = " + err.Error())
				}
				rpc = nil
			}
			continue
		}

		if rpc == nil {
			rpc, err = discordrpc.New("1426173818308788286")
			if err != nil {
				fmt.Println("scrobblecord failed! failed to create an rpc. make sure you have 'Share my activity' enabled in Discord settings.")
				fmt.Println("err = " + err.Error())
			}
			defer rpc.Socket.Close()
		}

		lastTrackName := lastTrack.Name
		lastTrackAlbum := lastTrack.Album.Name
		lastTrackArtist := lastTrack.Artist.Name
		var lastTrackArtwork string
		for _, image := range lastTrack.Images {
			if image.Size == "extralarge" {
				lastTrackArtwork = image.Url
				break
			}
			if image.Size == "large" {
				lastTrackArtwork = image.Url // just in case :D
			}
		}

		err = rpc.SetActivity(discordrpc.Activity {
			Type: 2,
			StatusType: 0,
			Name: lastTrackArtist,
			State: "by " + lastTrackArtist,
			Details: lastTrackName,
			Assets: &discordrpc.Assets{
				LargeImage: lastTrackArtwork,
				LargeText: "on " + lastTrackAlbum,
			},
		})
		if err != nil {
			fmt.Println("scrobblecord failed! Discord RPC SetActivity failed.")
			fmt.Println("err = " + err.Error())

			err = rpc.Socket.Close()
			if err != nil {
				fmt.Println("scrobblecord failed! failed to close the Discord RPC socket")
				fmt.Println("err = " + err.Error())
			}
			rpc = nil
		}
	}
}
