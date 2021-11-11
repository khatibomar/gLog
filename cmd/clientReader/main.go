package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/khatibomar/gLog/api/v1"
	"google.golang.org/grpc"
)

var (
	addr   = flag.String("addr", ":6666", "the address to connect")
	offset = flag.Uint64("offset", 0, "offset of record")
)

func main() {
	flag.Parse()
	if !strings.Contains(strings.Join(os.Args, " "), "offset") {
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

	r, err := c.Read(ctx, &pb.ConsumeRequest{
		Offset: *offset,
	})

	if err != nil {
		log.Fatalf("Client : %v", err)
	}
	log.Printf("{'Value': %s , 'Offset': %d}", r.Record.Value, r.Record.Offset)
}
