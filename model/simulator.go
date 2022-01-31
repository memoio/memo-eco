package model

import (
	"fmt"
	"math/big"
)

type SizeSimulateFunction func(state *MemoState, reward *big.Int, day int64, config *EconomicsConfig, lastDayOrder *Order) *big.Int
type PriceSimulateFunction func(state *MemoState, reward *big.Int, day int64, config *EconomicsConfig, lastDayOrder *Order) *big.Int
type DurationSimulateFunction func(state *MemoState, reward *big.Int, day int64, config *EconomicsConfig, lastDayOrder *Order) int64
type ProviderSimulateFunction func(state *MemoState, reward *big.Int, day int64, config *EconomicsConfig, lastDayOrder *Order) (int64, *big.Int)

var Level1 = false
var Level2 = false
var Level3 = false

// 每日订单空间的模拟变化函数
// 参数为，当前状态数据，调用这次函数时计算的奖励值
// 当前为项目启动后第多少天，以及昨天的订单空间
func DefaultSizeSimulate(state *MemoState, reward *big.Int, day int64, config *EconomicsConfig, lastDayOrder *Order) *big.Int {
	if day == 0 {
		return new(big.Int).Set(config.InitialSize) // 起始订单量
	}

	// 流动代币占比小于1/16的时候，订单数为0
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.0625 {
		size := big.NewInt(0)
		return size
	}

	// 流动代币突破瓶颈期，如果size为0，则设为初始值500GB
	if lastDayOrder.Size.Cmp(BigZero) <= 0 {
		fmt.Println("Iszero")
		size := big.NewInt(500)
		return size
	}

	// 流动代币占比小于1/8的时候，新增订单空间下降30%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.125 {
		size := new(big.Int).Mul(lastDayOrder.Size, big.NewInt(7_0000_0000))
		return size.Div(size, OneBillion)
	}

	// 流动代币占比小于1/4的时候，新增订单空间下降20%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) <= 0.25 {
		size := new(big.Int).Mul(lastDayOrder.Size, big.NewInt(8_0000_0000))
		return size.Div(size, OneBillion)
	}

	// 流动代币占比小于2/5的时候，新增订单空间下降5%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) <= 0.4 {
		size := new(big.Int).Mul(lastDayOrder.Size, big.NewInt(9_5000_0000))
		return size.Div(size, OneBillion)
	}

	// 100T以内，每天涨10%
	if lastDayOrder.Size.Cmp(big.NewInt(100*1024)) < 0 {
		size := new(big.Int).Mul(lastDayOrder.Size, big.NewInt(11_0000_0000))
		return size.Div(size, OneBillion)
	}

	if !Level1 {
		fmt.Println("-----------------------")
		fmt.Println("Day:", day, "Level1, OrderSize:", FormatGBytes(lastDayOrder.Size.Int64()),
			"reward:", WeiToMemo(reward, config.Decimals),
			"MintLevel:", state.MintLevel,
			"TotalSupply:", WeiToMemo(state.TotalSupply, config.Decimals),
			"TotalLiquid", WeiToMemo(state.TotalLiquid, config.Decimals))
		fmt.Println("-----------------------")
		Level1 = true
	}

	// 1PB以内，每天涨5%
	if lastDayOrder.Size.Cmp(big.NewInt(2*1024*1024)) < 0 {
		size := new(big.Int).Mul(lastDayOrder.Size, big.NewInt(10_5000_0000))
		return size.Div(size, OneBillion)
	}

	if !Level2 {
		fmt.Println("-----------------------")
		fmt.Println("Day:", day, "Level2, OrderSize:", FormatGBytes(lastDayOrder.Size.Int64()),
			"reward:", WeiToMemo(reward, config.Decimals),
			"MintLevel:", state.MintLevel,
			"TotalSupply:", WeiToMemo(state.TotalSupply, config.Decimals),
			"TotalLiquid", WeiToMemo(state.TotalLiquid, config.Decimals))
		fmt.Println("-----------------------")
		Level2 = true
	}

	// 10PB以内，每天涨2%
	if lastDayOrder.Size.Cmp(big.NewInt(10*1024*1024)) < 0 {
		size := new(big.Int).Mul(lastDayOrder.Size, big.NewInt(10_2000_0000))
		return size.Div(size, OneBillion)
	}

	if !Level3 {
		fmt.Println("-----------------------")
		fmt.Println("Day:", day, "Level3, OrderSize:", FormatGBytes(lastDayOrder.Size.Int64()),
			"reward:", WeiToMemo(reward, config.Decimals),
			"MintLevel:", state.MintLevel,
			"TotalSupply:", WeiToMemo(state.TotalSupply, config.Decimals),
			"TotalLiquid", WeiToMemo(state.TotalLiquid, config.Decimals))
		fmt.Println("-----------------------")
		Level3 = true
	}

	return new(big.Int).Set(lastDayOrder.Size)
}

// 每日订单价格的模拟变化函数
// 参数为，当前状态数据，调用这次函数时计算的奖励值
// 当前为项目启动后第多少天，以及昨天的订单价格
func DefaultPriceSimulate(state *MemoState, reward *big.Int, day int64, config *EconomicsConfig, lastDayOrder *Order) *big.Int {
	// 初始订单价格
	if day == 0 {
		return new(big.Int).Set(config.InitialPrice) // 每天，每GB 0.1Memo 订单平均价格
	}

	// 流动代币占比小于1/16的时候，Memo定价每天下降30%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.0625 {
		if lastDayOrder.Price.Cmp(config.MinimumPrice) > 0 {
			price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(7_0000_0000))
			return price.Div(price, OneBillion)
		}

		// Memo定价不变
		price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(10_0000_0000))
		return price.Div(price, OneBillion)
	}

	// 流动代币占比小于1/8的时候，Memo定价每天下降20%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.125 {
		if lastDayOrder.Price.Cmp(config.MinimumPrice) > 0 {
			price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(8_0000_0000))
			return price.Div(price, OneBillion)
		}

		// Memo定价不变
		price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(10_0000_0000))
		return price.Div(price, OneBillion)
	}

	// 流动代币占比小于1/4的时候，Memo定价每天下降10%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.25 {
		if lastDayOrder.Price.Cmp(config.MinimumPrice) > 0 {
			price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(9_0000_0000))
			return price.Div(price, OneBillion)
		}

		// Memo定价不变
		price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(10_0000_0000))
		return price.Div(price, OneBillion)
	}

	// 流动代币占比小于1/2的时候，Memo定价每天下降5%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.5 {
		if lastDayOrder.Price.Cmp(config.MinimumPrice) > 0 {
			price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(9_5000_0000))
			return price.Div(price, OneBillion)
		}
		// Memo定价不变
		price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(10_0000_0000))
		return price.Div(price, OneBillion)
	}

	// 流动代币占比小于2/3的时候，Memo定价不变
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.667 {
		price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(10_0000_0000))
		return price.Div(price, OneBillion)
	}

	// 流动代币占比小于2/3的时候，Memo定价每天上升5%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.667 {
		price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(10_5000_0000))
		return price.Div(price, OneBillion)
	}

	// 流动代币占比小于3/4的时候，Memo定价每天上升10%
	if float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()) < 0.75 {
		price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(11_0000_0000))
		return price.Div(price, OneBillion)
	}

	// 否则，每天增加
	price := new(big.Int).Mul(lastDayOrder.Price, big.NewInt(12_0000_0000))
	return price.Div(price, OneBillion)
}

// 简单策略，平均存一年
func DefaultDurationSimulate(state *MemoState, reward *big.Int, day int64, config *EconomicsConfig, lastDayOrder *Order) int64 {
	return 365 // 模拟的总周期，单位天
}

// 简单策略，每天增加50个Provider，同时返回一个Provider需要质押的数量
// TODO: 动态质押
func DefaultProviderSimulate(state *MemoState, reward *big.Int, day int64, config *EconomicsConfig, lastDayOrder *Order) (int64, *big.Int) {
	if state.ProviderCount < 8_0000 {
		return 50, new(big.Int).Set(config.InitialProviderPledge)
	}

	return 0, new(big.Int).Set(config.InitialProviderPledge)
}
