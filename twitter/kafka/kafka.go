package kafka

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

var (
	writerChan chan *twitter.Tweet = make(chan *twitter.Tweet)
	topicName  string              = "tweets"
)

type Producer struct {
	ID   string
	Base *kafkago.Writer
}
type key string

func producer() (*Producer, error) {

	writerConfig := kafkago.WriterConfig{
		Brokers:          []string{"localhost:9092"},
		Topic:            topicName,
		BatchSize:        20,
		CompressionCodec: snappy.NewCompressionCodec(),
	}

	return &Producer{
		ID:   "1",
		Base: kafkago.NewWriter(writerConfig),
	}, nil
}

func Start(ctx context.Context) chan *twitter.Tweet {
	go consumeTweets(ctx)
	return writerChan
}

func (p *Producer) writeMessage(msg *kafkago.Message) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, key("id"), p.ID)
	err := p.Base.WriteMessages(ctx, *msg)

	return err
}

func createHeader(key string, value string) kafkago.Header {
	return kafkago.Header{
		Key:   key,
		Value: []byte(value),
	}
}
func toKafkaHeaders(tweet *twitter.Tweet) []kafkago.Header {
	return []kafkago.Header{
		createHeader("tweetId", tweet.IDStr),
	}
}

func fromTweet(tweet *twitter.Tweet) (*kafkago.Message, error) {
	headers := toKafkaHeaders(tweet)
	/*tweetJSON, err := json.Marshal(tweet)

	if err != nil {
		return nil, err
	}*/

	return &kafkago.Message{
		Key:     []byte(tweet.IDStr),
		Headers: headers,
		Value:   []byte(tweet.IDStr),
	}, nil
}

func consumeTweets(ctx context.Context) {
	producer, _ := producer()
	for {
		select {
		case tweet := <-writerChan:
			msg, err := fromTweet(tweet)
			if err != nil {
				panic(err)
			}
			err = producer.writeMessage(msg)
			if err != nil {
				fmt.Println(err)
			}

		case <-ctx.Done():
			fmt.Println("Shutting down")
			producer.Shutdown()
			return
		}
	}
}

func (p *Producer) Shutdown() {
	p.Base.Close()
}
