package model

import "math/big"

type EconomicsConfig struct {
	MintLevel             []MintInfo // 增发阶段设计
	MinimumRation         *big.Int   // 最小增发率，减半到最小增发率后保持稳定
	InitialSupply         *big.Int   // 创世代币数量, 以 10^-8 Memo 为单位
	InitialTarget         *big.Int   // 创世第一阶段目标增发代币数, 以 10^-8 Memo 为单位，往后开始减半
	InitialKeeperPledge   *big.Int   // Keeper 初始需要质押的代币数
	InitialProviderPledge *big.Int   // Provider 初始需要质押的代币数

	InitialSize      *big.Int                 // 初始的订单空间
	InitialPrice     *big.Int                 // 初始的价格
	MinimumPrice     *big.Int                 // 最小的订单价格
	SizeSimulate     SizeSimulateFunction     // 模拟每天订单数据量大小的函数
	PriceSimulate    PriceSimulateFunction    // 模拟每天订单价格大小的函数
	DurationSimulate DurationSimulateFunction // 模拟每天订单的平均时间的函数
	ProviderSimulate ProviderSimulateFunction // 模拟每天新增Provider数据的函数

	TotalDuration int64 // 总统计周期，单位 day
}

func DefaultEconomicsConfig() *EconomicsConfig {
	return &EconomicsConfig{
		MintLevel: []MintInfo{
			{
				Ratio:    big.NewInt(5_0000_0000), // 增发比例 50%
				Size:     big.NewInt(100 * 1024),  // 100T
				Duration: 100,                     // 100 days
			},
			{
				Ratio:    big.NewInt(8_0000_0000), // 增发比例 80%
				Size:     big.NewInt(1024 * 1024), // 1PB
				Duration: 150,                     // 150 days
			},
			{
				Ratio:    big.NewInt(10_0000_0000),     // 增发比例 100%
				Size:     big.NewInt(50 * 1024 * 1024), // 50 PB
				Duration: 200,                          // 200 days
			},
			{
				Ratio:    big.NewInt(6_0000_0000),        // 增发比例 80%
				Size:     big.NewInt(1024 * 1024 * 1024), // 1EB
				Duration: 300,                            // 300 days
			},
			{
				Ratio:    big.NewInt(20_0000_0000),            // 增发比例 50%
				Size:     big.NewInt(50 * 1024 * 1024 * 1024), // 50EB
				Duration: 730,                                 // 730 days
			},
		},

		MinimumRation: big.NewInt(5000_0000),               // 最小增发比例5%
		InitialSupply: big.NewInt(100_0000_0000_0000_0000), // 设置初始发行量，100亿，精度为8，设Memo发行价为 0.0006U
		InitialTarget: big.NewInt(50_0000_0000_0000_0000),  // 初始增发奖励目标，50亿，然后25亿，以此类推，直到达到最小增发率

		InitialKeeperPledge:   big.NewInt(5000_0000_0000_0000), // 五千万 Memo
		InitialProviderPledge: big.NewInt(10_0000_0000_0000),   // 十万 Memo

		InitialSize:  big.NewInt(500),        // 500GB
		InitialPrice: big.NewInt(10000_0000), // 1 GB*Day/Memo
		MinimumPrice: big.NewInt(1_0000),     // 0.00001 Memo

		SizeSimulate:     DefaultSizeSimulate,
		PriceSimulate:    DefaultPriceSimulate,
		DurationSimulate: DefaultDurationSimulate,
		ProviderSimulate: DefaultProviderSimulate,

		TotalDuration: 4 * 365, // 模拟的总周期，单位天
	}
}
