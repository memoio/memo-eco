package model

import (
	"math/big"
	"testing"
)

func TestMintLevel(t *testing.T) {

	rt := big.NewInt(300_000_000)
	rt.Mul(rt, big.NewInt(Memo))

	t.Fatal()
}
