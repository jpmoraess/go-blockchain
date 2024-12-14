package types

import (
	"crypto/sha256"

	"github.com/jpmoraess/go-blockchain/crypto"
	"github.com/jpmoraess/go-blockchain/proto"

	pb "google.golang.org/protobuf/proto"
)

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	return pk.Sign(HashBlock(b))
}

// HashBlock returns SHA256 of the header.
func HashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}
