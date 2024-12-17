package main

import (
	"context"
	"log"

	"github.com/jpmoraess/go-blockchain/node"
	"github.com/jpmoraess/go-blockchain/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{})
	makeNode(":4000", []string{":3000"})

	select {}
}

func makeTransaction() {
	conn, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewNodeClient(conn)

	version := &proto.Version{
		Version:    "1.0",
		Height:     1,
		ListenAddr: ":4000",
	}

	_, err = client.Handshake(context.Background(), version)
	if err != nil {
		log.Fatal(err)
	}
}

func makeNode(listenAddr string, boostrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr)
	if len(boostrapNodes) > 0 {
		if err := n.BoostrapNetwork(boostrapNodes); err != nil {
			log.Fatal(err)
		}
	}
	return n
}
