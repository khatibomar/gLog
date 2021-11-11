package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/khatibomar/gLog/api/v1"
	glog "github.com/khatibomar/gLog/internal/log"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 6666, "the server port")
)

type server struct {
	Log glog.Log
	pb.UnimplementedLoggerServer
}

func (s *server) Append(ctx context.Context, in *pb.ProduceRequest) (*pb.ProduceResponse, error) {
	var record glog.Record
	record.Offset = in.Record.Offset
	record.Value = in.Record.Value
	off, err := s.Log.Append(record)
	if err != nil {
		return &pb.ProduceResponse{}, err
	}
	return &pb.ProduceResponse{
		Offset: off,
	}, nil
}

func (s *server) Read(ctx context.Context, in *pb.ConsumeRequest) (*pb.ConsumeResponse, error) {
	record, err := s.Log.Read(in.Offset)
	if err != nil {
		return &pb.ConsumeResponse{}, err
	}

	return &pb.ConsumeResponse{
		Record: &pb.Record{
			Offset: record.Offset,
			Value:  record.Value,
		},
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := &server{
		Log: glog.Log{},
	}

	pb.RegisterLoggerServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
