package main

import (
	"fmt"
	"log"
	"net"
	"tmsservice/proto"
	"tmsservice/service"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:4000")
	if err != nil {
		log.Fatalf("failed to listen : %v ", err)
	}
	server := grpc.NewServer()
	proto.RegisterTokenServiceServer(server, &service.TokenServer{})

	fmt.Println("listen to tms server ", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v ", err)
	}

}
