package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	seeds := strings.Split(os.Getenv("BROKERS"), ",")
	opts := []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("wproc"),
		kgo.ConsumeTopics(os.Getenv("TOPIC")),
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, nil)),
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		die("unable to create client: %v", err)
	}

	go consume(context.Background(), cl)

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	fmt.Println("received interrupt signal; closing client")
	done := make(chan struct{})
	go func() {
		defer close(done)
		cl.Close()
	}()

	select {
	case <-sigs:
		fmt.Println("received second interrupt signal; quitting without waiting for graceful close")
	case <-done:
	}
}

func consume(ctx context.Context, cl *kgo.Client) {
	tick := time.Tick(15 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			fmt.Println(time.Now())
			// A nil context so that PollFetches does not block indefinitely.
			fetches := cl.PollFetches(nil)
			if fetches.IsClientClosed() {
				return
			}
			fetches.EachError(func(t string, p int32, err error) {
				die("fetch err topic %s partition %d: %v", t, p, err)
			})
			if fetches.Empty() {
				fmt.Println("no records to process")
				continue
			}
			var seen int
			iter := fetches.RecordIter()
			for !iter.Done() {
				iter.Next()
				seen++
			}
			err := cl.CommitUncommittedOffsets(ctx)
			if err != nil {
				fmt.Printf("commit records failed: %v", err)
				continue
			}
			fmt.Printf("processed %d records\n", seen)
		}
	}
}

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
