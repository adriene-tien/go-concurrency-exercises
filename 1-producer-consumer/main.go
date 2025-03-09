//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, wg *sync.WaitGroup, ch chan *Tweet, errCh chan *error) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			errCh <- &err
			return
		}

		ch <- tweet
	}
}

func consumer(ch chan *Tweet, wg *sync.WaitGroup, errCh chan *error) {
	defer wg.Done()
	for {
		select {
		case _, ok := <-errCh:
			if ok {
				return
			}
		default:
		}
		t := <-ch
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	var wg sync.WaitGroup

	twtCh := make(chan *Tweet)
	errCh := make(chan *error)

	wg.Add(1)
	// Producer
	go producer(stream, &wg, twtCh, errCh)

	// Consumer
	wg.Add(1)
	go consumer(twtCh, &wg, errCh)

	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}

// My approach (closed book):
// Create an error channel and a tweets channel. The producer will send tweets to the tweet channel that the consumer will continuously poll from
// Once the tweet stream is over, we send an error to the error channel
// The consumer will also constantly poll the error channel. Once it receives a value at the error channel, we can return.
// We call Wait() on our waitgroup to ensure all the go routines (producer/consumer) have completed. wg.Done() is called at the end
// of both the goroutine calls to make sure the waitgroup counter is successfully decremented.

// Key things to remember: we need to actually make our channels with make(), not just instantiate it with var

// IMPROVEMENT: synchronous performance: 3.582s
// concurrent performance: 1.978s. pretty significant improvement!
