package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tukejonny/trial-grpc-xds/pb"
	"google.golang.org/grpc"
)

var (
	addr        = flag.String("addr", "", "接続先アドレス")
	reqCount    = flag.Int("count", 10, "リクエスト回数")
	reqDuration = flag.Duration("duration", 5*time.Second, "リクエスト間隔")
)

func main() {
	ctx := context.Background()

	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	name := fmt.Sprintf("client-%s", id.String())
	log.SetPrefix(name)

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewPingServiceClient(conn)
	for idx := 0; idx < int(*reqCount); idx++ {
		resp, err := client.Ping(ctx, &pb.PingRequest{Id: name})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[<==] idx=%d, resp=%v\n", idx, resp)
		time.Sleep(*reqDuration)
	}

	return
}
