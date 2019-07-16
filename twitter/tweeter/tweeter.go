package tweeter

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"ts-tweets/kafka"
)

func createClient() *twitter.Client {
	consumerKey := flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flag.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flag.String("access-token", "", "Twitter Access Token")
	accessSecret := flag.String("access-secret", "", "Twitter Access Secret")
	flag.Parse()

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	// Twitter Client
	return twitter.NewClient(httpClient)

}
func Start() {
	client := createClient()
	demux := twitter.NewSwitchDemux()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	count := 0
	tweetChannels := kafka.Start(ctx)

	fmt.Println("Starting Stream...")

	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"cats", "trump", "usa", "iran", "summer"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	demux.Tweet = func(tweet *twitter.Tweet) {
		tweetChannels <- tweet
		count++
		if count >= 1000 {
			cancel()
			stream.Stop()
		}
	}
	demux.HandleChan(stream.Messages)

}
