package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/khatibomar/gLog/api/v1"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", ":6666", "the address to connect")
	data = flag.String("data", "", "data to write to log")
)

func main() {
	flag.Parse()

	if *data == "" {
		flag.PrintDefaults()
		return
	}

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewLoggerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Append(ctx, &pb.ProduceRequest{
		Record: &pb.Record{
			Value: []byte(*data),
		},
	})

	if err != nil {
		log.Fatalf("Client : %v", err)
	}
	log.Printf("write succeed at offset: %d", r.Offset)
}
