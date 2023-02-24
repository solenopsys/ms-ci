package main

import (
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	dial, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080", nil)
	if err != nil {
		log.Panic(err)
	}
	for i := 0; i < 100000000; i++ {
		dial.WriteMessage(websocket.TextMessage, []byte("1663535884616-1663535873175 j324lkj23lj42l3j4lj23l4jl23j4l23"))
	}

}
