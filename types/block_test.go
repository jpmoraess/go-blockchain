package types

import (
	"testing"

	"github.com/jpmoraess/go-blockchain/crypto"
	"github.com/jpmoraess/go-blockchain/util"
	"github.com/stretchr/testify/assert"
)

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	assert.Equal(t, 32, len(hash))
}

func TestSignBlock(t *testing.T) {
	var (
		block   = util.RandomBlock()
		privKey = crypto.GeneratePrivateKey()
		pubKey  = privKey.Public()
	)

	sign := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sign.Bytes()))
	assert.True(t, sign.Verify(pubKey, HashBlock(block)))
}
