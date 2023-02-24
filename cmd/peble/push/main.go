package main

import (
	zmq "github.com/pebbe/zmq4"
)

//func reseive(s2 *zmq.Socket) {
//	for {
//		msg, _ :=
//		println(msg[0])
//	}
//
//}

func reseive(s2 *zmq.Socket) {
	var i = 0
	for {
		msg3, _ := s2.RecvMessage(0)

		i++
		if i%1000 == 0 {
			println(msg3[0])
			println(i)
		}

	}
}

func main() { // рабочая хрень.

	zctx, _ := zmq.NewContext()

	s, _ := zctx.NewSocket(zmq.DEALER)
	s.Connect("tcp://127.0.0.1:30456")

	s2, _ := zctx.NewSocket(zmq.DEALER)
	s2.Connect("tcp://127.0.0.1:30457")
	bytes := []byte("fromServ")
	go reseive(s2)
	for i := 0; i < 10000000000; i++ {

		if _, err := s.SendBytes(bytes, 0); err != nil {
			panic(err)
		}

	}

}
