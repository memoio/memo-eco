package model

import (
	"fmt"
	"math/big"
)

var OneHudred = big.NewInt(100)

// 以十亿为基础计算比例
var OneBillion = big.NewInt(10_0000_0000)

var BigZero = big.NewInt(0) // 0
var BigTwo = big.NewInt(2)  // 2

func WeiToMemo(m *big.Int) *big.Int {
	return new(big.Int).Div(m, big.NewInt(Memo))
}

const (
	Day  = 86400
	Memo = 1e18
	KiB  = 1024
	MiB  = 1048576
	GiB  = 1073741824

	TiB = 1024 * GiB
	PiB = 1024 * TiB
	EiB = 1024 * PiB
)

func FormatGBytes(i int64) (result string) {
	switch {
	case i >= EiB:
		result = fmt.Sprintf("%.02f EiB", float64(i)/EiB)
	case i >= PiB:
		result = fmt.Sprintf("%.02f PiB", float64(i)/PiB)
	case i >= TiB:
		result = fmt.Sprintf("%.02f TiB", float64(i)/TiB)
	default:
		result = fmt.Sprintf("%d GiB", i)
	}
	return
}
