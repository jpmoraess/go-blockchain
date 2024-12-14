package types

import (
	"testing"

	"github.com/jpmoraess/go-blockchain/crypto"
	"github.com/jpmoraess/go-blockchain/proto"
	"github.com/jpmoraess/go-blockchain/util"
	"github.com/stretchr/testify/assert"
)

// my balance 100 coins
// want to send 5 coins to "AAA address"
// 2 outputs
// 5 to the dude we wanna send
// 95 back to our address
func TestNewTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	fromAddress := fromPrivKey.Public().Address().Bytes()

	toPrivKey := crypto.GeneratePrivateKey()
	toAddress := toPrivKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}

	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}
	sign := SignTansaction(fromPrivKey, tx)
	input.Signature = sign.Bytes()

	assert.True(t, VerifyTransaction(tx))
}
