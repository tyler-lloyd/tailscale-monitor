package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
	"tswatcher/ratelimiter"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"tailscale.com/client/tailscale"
)

var localclient tailscale.LocalClient

func main() {
	ctx := context.Background()
	consecutiveOffline := 0
	bucket := ratelimiter.NewBucket(5 * time.Minute)
	go bucket.Start()

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	apiKey := os.Getenv("TWILIO_API_KEY")
	apiSecret := os.Getenv("TWILIO_API_SECRET")

	for {
		st, err := localclient.Status(ctx)
		if err != nil {
			panic(err)
		}

		online := false
		for _, p := range st.Peer {
			if p.HostName == "Pixel 5a" && p.Online {
				online = true
			}
		}

		if !online {
			log.Default().Printf("Pixel 5a not online")
			consecutiveOffline++
		} else {
			log.Default().Printf("pixel 5a online")
			consecutiveOffline = 0
		}

		if consecutiveOffline == 3 && bucket.RequestToken() {
			client := twilio.NewRestClientWithParams(twilio.ClientParams{
				Username:   apiKey,
				Password:   apiSecret,
				AccountSid: accountSid,
			})
			log.Default().Printf("sending alert for being offline for 3 consecutive checks")
			params := &twilioApi.CreateMessageParams{}
			params.SetTo("+17046617694")
			params.SetFrom("+19785413960")
			params.SetBody("🚨🚨🚨")

			resp, err := client.Api.CreateMessage(params)
			if err != nil {
				panic(err)
			}

			b, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}
			log.Default().Printf("response: %s", string(b))
		}

		time.Sleep(10 * time.Second)
	}
}
