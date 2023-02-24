package main

import (
	"log"

	"github.com/zeromq/goczmq"
)

func main() {

	// Create a dealer socket and connect it to the router.
	dealer, err := goczmq.NewDealer("tcp://127.0.0.1:30456")
	if err != nil {
		log.Fatal(err)
	}
	defer dealer.Destroy()
	request, err := dealer.RecvMessage()
	log.Println(request)

	log.Println("dealer created and connected")

	err = dealer.SendFrame([]byte("22"), goczmq.FlagNone)
	if err != nil {
		log.Fatal(err)
	}

	// Send a 'Hello' message from the dealer to the router.
	// Here we send it as a frame ([]byte), with a FlagNone
	// flag to indicate there are no more frames following.
	err = dealer.SendFrame([]byte("Hello"), goczmq.FlagNone)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("dealer sent 'Hello'")

	//
	//// Receive the reply.
	//reply, err := dealer.RecvMessage()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("dealer received '%s'", string(reply[0]))
}
