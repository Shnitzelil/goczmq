package main

import (
	"flag"
	czmq "github.com/taotetek/goczmq"
	"log"
	"time"
)

func main() {
	var messageSize = flag.Int("message_size", 0, "size of message")
	var messageCount = flag.Int("message_count", 0, "number of messages")
	flag.Parse()

	pullSock := czmq.NewZsock(czmq.PULL)
	defer pullSock.Destroy()

	_, err := pullSock.Bind("inproc://test")
	if err != nil {
		panic(err)
	}

	go func() {
		pushSock := czmq.NewZsock(czmq.PUSH)
		defer pushSock.Destroy()
		err := pushSock.Connect("inproc://test")
		if err != nil {
			panic(err)
		}

		for i := 0; i < *messageCount; i++ {
			payload := make([]byte, *messageSize)
			err = pushSock.SendBytes(payload, 0)
			if err != nil {
				panic(err)
			}
		}
	}()

	startTime := time.Now()
	for i := 0; i < *messageCount; i++ {
		msg, _, err := pullSock.RecvBytes()
		if err != nil {
			panic(err)
		}
		if len(msg) != *messageSize {
			panic("msg too small")
		}
	}
	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	throughput := float64(*messageCount) / elapsed.Seconds()
	megabits := float64(throughput*float64(*messageSize)*8.0) / 1e6

	log.Printf("message size: %d", *messageSize)
	log.Printf("message count: %d", *messageCount)
	log.Printf("test time (seconds): %f", elapsed.Seconds())
	log.Printf("mean throughput: %f [msg/s]", throughput)
	log.Printf("mean throughput: %f [Mb/s]", megabits)
}
