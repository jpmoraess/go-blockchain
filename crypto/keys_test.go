package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), privKeyLen)

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), pubKeyLen)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("go blockchain")

	sign := privKey.Sign(msg)
	assert.True(t, sign.Verify(pubKey, msg))

	// test fake message
	assert.False(t, sign.Verify(pubKey, []byte("fake message")))

	// test with invalid public key
	fakePrivKey := GeneratePrivateKey()
	assert.False(t, sign.Verify(fakePrivKey.Public(), msg))
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "d412b9ed546ee4c3be60a6523d8cda720554d05a602284a8aa83a8fa60f440d1"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "c7f2e6e120f9a92b3688684dfb05b4ef429273b7"
	)

	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
	address := privKey.Public().Address()
	assert.Equal(t, address.String(), addressStr)
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	assert.Equal(t, addressLen, len(address.Bytes()))
	fmt.Println(address)
}
