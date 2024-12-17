package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/jpmoraess/go-blockchain/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version    string
	listenAddr string

	logger *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.Version

	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	logger, _ := zap.NewProduction()

	return &Node{
		peers:   make(map[proto.NodeClient]*proto.Version),
		version: "1.0",
		logger:  logger.Sugar(),
	}
}

func (n *Node) addPeer(client proto.NodeClient, version *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.logger.Infow("new peer connected", "addr", version.ListenAddr, "height", version.Height)

	n.peers[client] = version
}

func (n *Node) deletePeer(client proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, client)
}

func (n *Node) BoostrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		client, err := makeNodeClient(addr)
		if err != nil {
			return err
		}

		version, err := client.Handshake(context.Background(), n.getVersion())
		if err != nil {
			n.logger.Error("handshake error:", err)
			continue
		}

		n.addPeer(client, version)
	}

	return nil
}

func (n *Node) Start(listenAddr string) error {
	n.listenAddr = listenAddr

	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	var (
		opts       = []grpc.ServerOption{}
		grpcServer = grpc.NewServer(opts...)
	)

	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Infow("node started...", "port", n.listenAddr)

	return grpcServer.Serve(listen)
}

func (n *Node) Handshake(ctx context.Context, version *proto.Version) (*proto.Version, error) {
	client, err := makeNodeClient(version.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(client, version)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)

	fmt.Println("receiveed tx from:", peer)
	return &proto.Ack{}, nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "1.0",
		Height:     0,
		ListenAddr: n.listenAddr,
	}
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	conn, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(conn), nil
}
