package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"time"
)

func main() {
	zctx, _ := zmq.NewContext()

	s, _ := zctx.NewSocket(zmq.PULL)
	if err := s.Connect("tcp://127.0.0.1:30456"); err != nil {
		panic(err)
	}

	var i = 0
	for {
		if msg, err := s.Recv(0); err != nil {
			panic(err)
		} else {
			i++

			if i%10000000 == 0 {
				println(time.Now().String())
				fmt.Printf("Received %s\n", msg)
			}

		}
	}
}
