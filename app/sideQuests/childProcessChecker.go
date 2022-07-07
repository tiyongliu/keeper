package sideQuests

import (
	"keeper/app/modules"
	"time"
)

var counter int

func childProcessChecker(ch chan *modules.EchoMessage) {
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			counter = counter + 1
			ch <- &modules.EchoMessage{
				Payload: counter,
				MsgType: "ping",
			}
		}
	}()
}
