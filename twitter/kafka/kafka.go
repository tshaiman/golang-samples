package kafka

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

var (
	writerChan chan *twitter.Tweet = make(chan *twitter.Tweet, 100)
)

func Start(ctx context.Context) chan *twitter.Tweet {
	go consumeTweets(ctx)
	return writerChan
}
func consumeTweets(ctx context.Context) {
	for {
		select {
		case tweet := <-writerChan:
			fmt.Printf("[%s]:[%s] \n", tweet.IDStr, tweet.Text)
		case <-ctx.Done():
			fmt.Println("The consumer was shut down")
			return
		}
	}
}
