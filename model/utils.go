package model

import (
	"fmt"
	"math/big"
)

// 以十亿为基础计算比例
var OneBillion = big.NewInt(10_0000_0000)

// Memo的代币精度
var Decimals = big.NewInt(1_0000_0000) // 用户换算成Memo

var BigZero = big.NewInt(0) // 0
var BigTwo = big.NewInt(2)  // 2

func WeiToMemo(m *big.Int) *big.Int {
	return new(big.Int).Div(m, Decimals)
}

const (
	TiB = 1024
	PiB = 1024 * 1024
	EiB = 1024 * 1024 * 1024
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
