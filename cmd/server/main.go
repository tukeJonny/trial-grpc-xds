package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

var (
	port = flag.String("port", "8080", "リッスンするポート番号")
)

func main() {
	flag.Parse()

	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	name := fmt.Sprintf("server-%s", id.String())
	log.SetPrefix(name)

	lis, err := net.Listen("tcp", net.JoinHostPort("localhost", *port))
	if err != nil {
		log.Fatal(err)
	}

	pingSvc := newPingService(name)
	if err := pingSvc.Serve(lis); err != nil {
		log.Fatal(err)
	}

	return
}
