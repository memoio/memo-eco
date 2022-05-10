package model

import (
	"math/big"
)

// 配置参数
const (
	TotalSupply  = 600_000_000 // issue price: 1 $/Memo
	InitSupply   = 300_000_000 // all is available or unlock linearly
	RewardTarget = 300_000_000

	KeeperPledge   = 1_000 // Memo
	ProviderPledge = 10    // Memo
	KCntPerGroup   = 10    // each group has 10 keeper
	PCntPerGroup   = 5_000 //assume each provider has 1TB, one group has 5000 provider
	ProCreate      = 100   // one day

	DefaultSize     = 8 * GiB // 100 KiB/second; 24 hours
	DefaultDuration = 365 * Day
	MinDuration     = 100 * Day
	MaxDuration     = 1000 * Day
	DefaultPrice    = 1 * GiB // 1 attoMemo per byte

	RatioInit          = 1 // per second*byte
	RatioDecimal       = 1_000_000_000
	RatioReduce        = 115
	RatioReduceDecimal = 100

	LinearRate = 3
	EndRate    = 1
	TaxRate    = 1
)

type Config struct {
	TotalSupply    *big.Int // 创世代币数量
	InitSupply     *big.Int
	RewardTarget   *big.Int // 目标增发代币数
	KeeperPledge   *big.Int // Keeper 初始需要质押的代币数
	ProviderPledge *big.Int // Provider 初始需要质押的代币数

	TargetPrice *big.Int // 最小的订单价格

	TotalDuration int64 // 总统计周期，单位 day
}

func DefaultEconomicsConfig() *Config {
	return &Config{
		TotalSupply:  new(big.Int).Mul(big.NewInt(TotalSupply), big.NewInt(Memo)),
		InitSupply:   new(big.Int).Mul(big.NewInt(InitSupply), big.NewInt(Memo)),
		RewardTarget: new(big.Int).Mul(big.NewInt(RewardTarget), big.NewInt(Memo)),

		KeeperPledge:   new(big.Int).Mul(big.NewInt(KeeperPledge), big.NewInt(Memo)),
		ProviderPledge: new(big.Int).Mul(big.NewInt(ProviderPledge), big.NewInt(Memo)),

		TotalDuration: 2 * 365, // 模拟的总周期，单位天
	}
}
