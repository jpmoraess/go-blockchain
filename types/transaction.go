package types

import (
	"crypto/sha256"

	"github.com/jpmoraess/go-blockchain/crypto"
	"github.com/jpmoraess/go-blockchain/proto"

	pb "google.golang.org/protobuf/proto"
)

func SignTansaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		var (
			sign   = crypto.SignatureFromBytes(input.Signature)
			pubKey = crypto.PublicKeyFromBytes(input.PublicKey)
		)
		// TODO: make sure we dont run into problems after verification
		// cause we have set the signature to nil.
		input.Signature = nil
		if !sign.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
	}
	return true
}
