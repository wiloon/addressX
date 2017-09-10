package main

import (
	"log"
	"net"
	"google.golang.org/grpc"
	pb "wiloon.com/ipx/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
)

const address = ":7000"

type server struct{}

func (s *server) SetIp(ctx context.Context, request *pb.AddressRequest) (*pb.AddressReply, error) {
	log.Printf("set ip: %v", request.Ip)
	return &pb.AddressReply{Reply: true}, nil
}

func main() {
	log.Println("server starting...")

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening:%v", address)

	s := grpc.NewServer()
	pb.RegisterAddressServer(s, &server{})
	reflection.Register(s)

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
