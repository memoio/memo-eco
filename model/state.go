package model

import "math/big"

type MintInfo struct {
	Ratio    *big.Int // 增发比例
	Size     *big.Int // 期望空间
	Duration int64    // 期望周期
}

type Order struct {
	Size        *big.Int // 每日新增订单空间，单位为 (GB)
	Price       *big.Int // 订单平均价格，单位为 (GB * Day / 10^-8 Memo)
	Dur         int64    // 订单平均时长，单位为 (Day)
	NewProvider int64    // 当天新增Provider数目
}

type MemoState struct {
	MintLevel       int      // 当前的增发阶段
	Ratio           *big.Int // 当前的增发比例
	TotalLiquid     *big.Int // 总当前可流动代币数
	TotalSupply     *big.Int // 当前代币发行总量 单位为 10^-8 Memo
	TotalReward     *big.Int // 总增发奖励数
	TotalPledge     *big.Int // 总质押代币数 单位为 10^-8 Memo
	TotalSize       *big.Int // 总存储空间 单位为GB
	TotalSpaceTime  *big.Int // 总存储时空值 单位为 GB * Day
	TotalSpacePrice *big.Int // 总存储空间价格值 单位为 GB * 10^-8 Memo
	TotalPay        *big.Int // 总支付代币数 单位为 10^-8 Memo
	TotalPaid       *big.Int // 总已支付代币数 单位为 10^-8 Memo
	TargetReward    *big.Int // 当前目标减半的累积奖励金额
	PeriodReward    *big.Int // 当前阶段的增发总和，增发完则跳到下一个阶段减半
	HalfFactor      int64    // 当前增发率除一个2^f

	KeeperPledge     *big.Int           // Keeper需要质押的代币数
	ProviderPledge   *big.Int           // Provider需要质押的代币数
	KeeperCount      int64              // 系统中Keeper的数量
	ProviderCount    int64              // 系统中Provider的数量
	LastMint         int64              // 上一次Mint的时间，单位Day
	SubSpacePriceMap map[int64]*big.Int // 记录过期价格，key为天数
	SubSizeMap       map[int64]*big.Int // 记录过期空间，key为天数
}

func NewMemoState(config *EconomicsConfig) *MemoState {
	return &MemoState{
		MintLevel:       0,
		Ratio:           big.NewInt(0),
		TotalSupply:     new(big.Int).Set(config.InitialSupply),
		TotalLiquid:     new(big.Int).Set(config.InitialSupply),
		TotalReward:     big.NewInt(0),
		TotalPledge:     big.NewInt(0),
		TotalSize:       big.NewInt(0),
		TotalSpacePrice: big.NewInt(0),
		TotalSpaceTime:  big.NewInt(0),
		TotalPay:        big.NewInt(0),
		TotalPaid:       big.NewInt(0),
		TargetReward:    new(big.Int).Set(config.InitialTarget),
		PeriodReward:    new(big.Int).Set(config.InitialTarget),
		HalfFactor:      0,

		KeeperPledge:     new(big.Int).Set(config.InitialKeeperPledge),
		ProviderPledge:   new(big.Int).Set(config.InitialProviderPledge),
		KeeperCount:      0,
		ProviderCount:    0,
		LastMint:         0,
		SubSpacePriceMap: make(map[int64]*big.Int),
		SubSizeMap:       make(map[int64]*big.Int),
	}
}
