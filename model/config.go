package model

import "math/big"

const (
	Day = 86400
	Memo = 1e18
	KiB = 1024
	MiB = 1048576
	GiB = 1073741824
	TiB = 1099511627776
	PiB = 1024*TiB
	EiB = 1024*PiB
)	

const (
	DefaultSize = 8*GiB // 100 KiB/second; 24 hours
	DefaultDuration = 100*Day  
	MinDuration = 100*Day
	MaxDuration = 1000*Day
	DefaultPrice = 1*GiB    // 1 attoMemo per byte  
)

type MintInfo struct {
	Ratio    *big.Int // 增发比例
	Size     *big.Int // 期望空间
	Reward   *big.Int 
}

type EconomicsConfig struct {
	MintLevel             []MintInfo // 增发阶段设计
	Decimals              *big.Int   // Memo代币精度
	TotalSupply           *big.Int   // 创世代币数量, 以 10^-8 Memo 为单位
	RewardTarget          *big.Int   // 创世第一阶段目标增发代币数, 以 10^-8 Memo 为单位，往后开始减半
	InitialKeeperPledge   *big.Int   // Keeper 初始需要质押的代币数
	InitialProviderPledge *big.Int   // Provider 初始需要质押的代币数

	TargetPrice     *big.Int                  // 最小的订单价格
	SizeSimulate     SizeSimulateFunction     // 模拟每天订单数据量大小的函数
	PriceSimulate    PriceSimulateFunction    // 模拟每天订单价格大小的函数
	DurationSimulate DurationSimulateFunction // 模拟每天订单的平均时间的函数
	ProviderSimulate ProviderSimulateFunction // 模拟每天新增Provider数据的函数

	TotalDuration int64 // 总统计周期，单位 day
}

func MintConfig(rt *big.Int) []MintInfo {
	stageRatio := big.NewInt(2_000_000);	
	stageSize := big.NewInt(TiB)
	stageReward := new(big.Int).Mul(stageSize, stageRatio)
	stageReward.Mul(stageReward, big.NewInt(864))



	mint = make([]MintInfo, 0, 32)

	

	for i:=0;i<32;i++ {
		mi := MintInfo {
			Ratio: new(big.Int).Set(stageRatio),
			Size : new(big.Int).Set(stageSize),
			Reward:  new(big.Int).Set(stageReward), 
		}
		mint = append(mint,mi)
		rt.Sub(rt, stageReward)

		fmt.Print(stageRatio, stageSize, stageReward)

		tmpSize = new(big.Int).Set(stageSize)
		if stageSize.Cmp(big.NewInt(EiB)) >= 0 {
			stageSize = stageSize.Mul(stageSize, big.NewInt(2))
		} else {
			stageSize = stageSize.Mul(stageSize, big.NewInt(32))
		}

		stageRatio.Mul(stageRatio, 9)
		stageRatio.Div(stageRatio, 10)

		tmpSize.Sub(stageSize, tmpSize) 
		stageReward.Mul(tmpSize, stageRatio)
		stageReward.Mul(stageReward, big.NewInt(864))
		if stageReward.Cmp(rt) >= 0 {
			tmpSize.Mul(tmpSize, big.NewInt(864*2))
			stageRatio = stageRatio.Div(rt, tmpSize)
			stageReward.Mul(tmpSize, stageRatio)
			stageReward.Mul(stageReward, big.NewInt(864))
		} 
	}
}

func DefaultEconomicsConfig() *EconomicsConfig {
	decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil) // 精度为8

	

	return &EconomicsConfig{
		MintLevel: mint,

		TotalSupply: new(big.Int).Mul(big.NewInt(1000_0000_0000), decimals),                   // 设置初始发行量，1000亿，精度为8，设Memo发行价为 0.0002U
		RewardTarget: new(big.Int).Mul(big.NewInt(500_0000_0000), decimals),                    // 初始增发奖励目标500亿，到达后减半至250亿，以此类推，直到达到最小增发率

		KeeperPledge:   new(big.Int).Mul(big.NewInt(5000_0000), decimals), // 五千万 Memo
		ProviderPledge: new(big.Int).Mul(big.NewInt(1000000), decimals),   // 一百万 Memo

		SizeSimulate:     DefaultSizeSimulate,
		PriceSimulate:    DefaultPriceSimulate,
		DurationSimulate: DefaultDurationSimulate,
		ProviderSimulate: DefaultProviderSimulate,

		TotalDuration: 4 * 365, // 模拟的总周期，单位天
	}
}
