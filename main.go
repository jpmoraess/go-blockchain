package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jpmoraess/go-blockchain/node"
	"github.com/jpmoraess/go-blockchain/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	listen, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	node := node.NewNode()

	proto.RegisterNodeServer(grpcServer, node)

	fmt.Println("node running on port:", ":3000")

	go func() {
		time.Sleep(3 * time.Second)
		makeTransaction()
	}()

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("grpc server fail: %+v", err)
	}
}

func makeTransaction() {
	conn, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewNodeClient(conn)

	version := &proto.Version{
		Version: "1.0",
		Height:  1,
	}

	_, err = client.Handshake(context.Background(), version)
	if err != nil {
		log.Fatal(err)
	}
}
