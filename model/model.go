package model

import (
	"fmt"
	"math/big"

	"gonum.org/v1/plot/plotter"
)

const (
	POINTS_COUNT = 7

	SUPPLY_INDEX = 0
	LIQUID_INDEX = 1
	REWARD_INDEX = 2
	PLEDGE_INDEX = 3
	PAY_INDEX    = 4
	PAID_INDEX   = 5
	SIZE_INDEX   = 6
)

// 生成每天的模拟数据
func EcoModelSimulate(config *EconomicsConfig) []plotter.XYs {
	state := NewMemoState(config)
	state.Ratio.Set(config.MintLevel[0].Ratio)

	// 计算这么多天的增发
	pts := make([]plotter.XYs, POINTS_COUNT)

	for i := 0; i < len(pts); i++ {
		pts[i] = make(plotter.XYs, config.TotalDuration)
	}

	fmt.Println("Memo initial supply:", WeiToMemo(state.TotalSupply, config.Decimals))
	fmt.Println("Memo initial Liquid:", WeiToMemo(state.TotalLiquid, config.Decimals))

	// 设置第一天的数据
	// 假设初始有7个Keeper，计入质押，退出流动
	initialKeeperPledge := new(big.Int).Mul(big.NewInt(7), state.KeeperPledge)
	state.TotalPledge.Add(
		state.TotalPledge,
		initialKeeperPledge,
	)
	state.KeeperCount += 7
	state.TotalLiquid.Sub(state.TotalLiquid, initialKeeperPledge)

	// 假设初始有500个Provider，计入质押，退出流动
	initialProviderPledge := new(big.Int).Mul(big.NewInt(500), state.ProviderPledge)
	state.TotalPledge.Add(
		state.TotalPledge,
		initialProviderPledge,
	)
	state.ProviderCount += 500
	state.TotalLiquid.Sub(state.TotalLiquid, initialProviderPledge)

	// 代币总量数据
	pts[SUPPLY_INDEX][0].X = 0
	// 单位换算成Memo
	pts[SUPPLY_INDEX][0].Y = float64(WeiToMemo(state.TotalSupply, config.Decimals).Uint64())

	// 流通代币数据
	pts[LIQUID_INDEX][0].X = 0
	// 单位换算成Memo
	pts[LIQUID_INDEX][0].Y = float64(WeiToMemo(state.TotalLiquid, config.Decimals).Uint64())

	// 奖励数据
	pts[REWARD_INDEX][0].X = 0
	pts[REWARD_INDEX][0].Y = 0

	// 质押数据
	pts[PLEDGE_INDEX][0].X = 0
	pts[PLEDGE_INDEX][0].Y = float64(WeiToMemo(state.TotalPledge, config.Decimals).Uint64())

	// 支付数据
	pts[PAY_INDEX][0].X = 0
	pts[PAY_INDEX][0].Y = 0

	// 已支付数据
	pts[PAID_INDEX][0].X = 0
	pts[PAID_INDEX][0].Y = 0

	// 存储空间数据
	pts[SIZE_INDEX][0].X = 0
	pts[SIZE_INDEX][0].Y = 0

	// 将每天所有的订单抽象成一个订单
	// 第一天的初始订单
	order := &Order{
		Size:  config.SizeSimulate(state, BigZero, 0, config, nil),
		Price: config.PriceSimulate(state, BigZero, 0, config, nil),
		Dur:   config.DurationSimulate(state, BigZero, 0, config, nil),
	}

	order.NewProvider, state.ProviderPledge = config.ProviderSimulate(state, BigZero, 0, config, nil)

	// 开始模拟每天的订单
	for i := 1; i < int(config.TotalDuration); i++ {
		// 横坐标为开始的第多少天
		for j := 0; j < len(pts); j++ {
			pts[j][i].X = float64(i)
		}

		// 订单的起始时间与结束时间
		start := int64(i)
		// 这一天没有新订单
		if order.Size.Cmp(BigZero) > 0 {
			end := start + order.Dur

			// 该订单的 spacePrice
			spacePrice := new(big.Int).Mul(order.Price, order.Size)
			// 将 sub spacePrice and size
			subSpacePrice, ok := state.SubSpacePriceMap[end]
			if ok {
				subSpacePrice.Add(subSpacePrice, spacePrice)
			} else {
				state.SubSpacePriceMap[end] = spacePrice
			}

			subSize, ok := state.SubSizeMap[end]
			if ok {
				subSize.Add(subSize, order.Size)
			} else {
				state.SubSizeMap[end] = new(big.Int).Set(order.Size)
			}

			// 将订单的效果反应到全局状态中去
			state.TotalSize.Add(state.TotalSize, order.Size)
			state.TotalSpacePrice.Add(state.TotalSpacePrice, spacePrice)

			// 计算订单的数额
			spacetime := new(big.Int).Set(order.Size)
			spacetime.Mul(spacetime, big.NewInt(order.Dur))
			pay := new(big.Int).Mul(spacetime, order.Price)

			// 修改状态，pay代币暂时退出流通
			state.TotalLiquid.Sub(state.TotalLiquid, pay)
			state.TotalPay.Add(state.TotalPay, pay)
		}

		// 减去以前的订单数据
		sp, ok := state.SubSpacePriceMap[start]
		if ok {
			state.TotalSpacePrice.Sub(state.TotalSpacePrice, sp)
			delete(state.SubSpacePriceMap, start)
		}

		sSize, ok := state.SubSizeMap[start]
		if ok {
			state.TotalSize.Sub(state.TotalSize, sSize)
			delete(state.SubSizeMap, start)
		}

		timeNow := start
		dur := big.NewInt(int64(timeNow - state.LastMint))
		paid := new(big.Int).Mul(state.TotalSpacePrice, dur)

		// 计算到目前为止已支付的
		state.TotalPaid.Add(state.TotalPaid, paid)
		// 已支付的可以继续流动
		state.TotalLiquid.Add(state.TotalLiquid, paid)
		// 计算上次增发到目前为止的累积时空值
		state.TotalSpaceTime.Add(state.TotalSpaceTime, new(big.Int).Mul(state.TotalSize, dur))

		tempEsize := new(big.Int).Set(config.MintLevel[state.MintLevel].Size)

		// 选择当前的MintLevel
		for j := state.MintLevel + 1; j < len(config.MintLevel); j++ {
			esize := new(big.Int).Set(config.MintLevel[j].Size)

			if esize.Cmp(state.TotalSize) < 0 {
				esize.Set(state.TotalSize)
			}

			tempEsize.Set(esize)

			// 切换到新的增发阶段
			if new(big.Int).Div(state.TotalSpaceTime, esize).Cmp(big.NewInt(config.MintLevel[j].Duration)) >= 0 {
				fmt.Println("----------Change MintLevel-------------")
				fmt.Println("Day:", i, "Change MintLevel to", j, "Total Size:", state.TotalSize)
				fmt.Println("spacetime:", state.TotalSpaceTime.String(), "esize:", FormatGBytes(tempEsize.Int64()))
				fmt.Println("Calculate Dur:", new(big.Int).Div(state.TotalSpaceTime, tempEsize), "expect dur:", dur)
				state.MintLevel = j
				factor := new(big.Int).Exp(BigTwo, big.NewInt(state.HalfFactor), BigZero)
				state.Ratio.Div(config.MintLevel[state.MintLevel].Ratio, factor)
				// 如果增发比例已经小于最小比例，则设为最小比例
				if state.Ratio.Cmp(config.MinimumRation) < 0 {
					state.Ratio.Set(config.MinimumRation)
				}
			} else {
				break
			}
		}

		reward := big.NewInt(0)

		// 当该订单处理完，关闭enableMaxSize
		// 或者当前Size要大于历史上的最大Size时，才会增发
		if state.TotalSize.Cmp(state.MaxSize) > 0 || !config.EnableMaxSize {
			// 更新 MaxSize
			state.MaxSize.Set(state.TotalSize)
			// 计算当前的奖励
			reward = new(big.Int).Mul(paid, state.Ratio)
			// 将奖励除以基数
			reward.Div(reward, OneBillion)
			// 计算临时的累积奖励
			tempReward := new(big.Int).Add(reward, state.TotalReward)

			// 处理边界条件，如果奖励已经超过目标，超出部分需要除2
			isBig := tempReward.Cmp(state.TargetReward)
			// 比例达到最小增发比例时不再减半
			if isBig > 0 && state.Ratio.Cmp(config.MinimumRation) > 0 {
				// 先增发目标内的代币数
				leftReward := new(big.Int).Sub(state.TargetReward, state.TotalReward)
				overflowReward := new(big.Int).Sub(tempReward, state.TargetReward)
				reward.Add(leftReward, overflowReward.Div(overflowReward, BigTwo))

				// 增发计算后奖励
				state.TotalReward.Add(state.TotalReward, reward)
				state.TotalSupply.Add(state.TotalSupply, reward)
				state.TotalLiquid.Add(state.TotalLiquid, reward)

				// 更换目标
				state.HalfFactor += 1

				factor := new(big.Int).Exp(BigTwo, big.NewInt(state.HalfFactor), BigZero)
				state.Ratio.Div(config.MintLevel[state.MintLevel].Ratio, factor)

				// 如果增发比例已经小于最小比例，则设为最小比例
				if state.Ratio.Cmp(config.MinimumRation) < 0 {
					state.Ratio.Set(config.MinimumRation)
				} else {
					state.PeriodReward.Div(state.PeriodReward, BigTwo)
					state.TargetReward.Add(state.TargetReward, state.PeriodReward)
				}
				fmt.Println("----------Change HalfFactor-------------")
				fmt.Println("Day:", i, "Change HalfFactor to", state.HalfFactor, "reward:", WeiToMemo(reward, config.Decimals), "MintLevel:", state.MintLevel)
				fmt.Println("spacetime:", state.TotalSpaceTime.String(), "esize:", FormatGBytes(tempEsize.Int64()))
				fmt.Println("Calculate Dur:", float64(new(big.Int).Div(state.TotalSpaceTime, tempEsize).Uint64()))
				fmt.Println("TotalPay:", WeiToMemo(state.TotalPay, config.Decimals), "TotalPaid:", WeiToMemo(state.TotalPaid, config.Decimals), "TotalReward", WeiToMemo(state.TotalReward, config.Decimals))
				fmt.Println("TotalLiquid:", WeiToMemo(state.TotalLiquid, config.Decimals),
					"TotalSupply:", WeiToMemo(state.TotalSupply, config.Decimals),
					"Percent:", float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()))
				fmt.Println("TotalSize:", FormatGBytes(state.TotalSize.Int64()))
				fmt.Println("OrderSize:", FormatGBytes(order.Size.Int64()), "OrderPrice", order.Price)
				fmt.Println("Issurance ratio:", float64(state.Ratio.Int64())/float64(OneBillion.Int64()), "Providers Count", state.ProviderCount)
				fmt.Println("Total issue times:", float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64())/float64(WeiToMemo(config.InitialSupply, config.Decimals).Int64()))
				fmt.Println("TotalPledge:", WeiToMemo(state.TotalPledge, config.Decimals))
				// 刚好达到目标
			} else {
				// 直接增发所有奖励
				state.TotalReward.Add(state.TotalReward, reward)
				state.TotalSupply.Add(state.TotalSupply, reward)
				state.TotalLiquid.Add(state.TotalLiquid, reward)
				// 如果已达到目标，接下来的都会受减半因子影响
				// 比例达到最小增发比例时不再减半
				if isBig == 0 && state.Ratio.Cmp(config.MinimumRation) > 0 {
					factor := new(big.Int).Exp(BigTwo, big.NewInt(state.HalfFactor), BigZero)
					state.Ratio.Div(config.MintLevel[state.MintLevel].Ratio, factor)

					// 如果增发比例已经小于最小比例，则设为最小比例
					if state.Ratio.Cmp(config.MinimumRation) < 0 {
						state.Ratio.Set(config.MinimumRation)
					} else {
						state.PeriodReward.Div(state.PeriodReward, BigTwo)
						state.TargetReward.Add(state.TargetReward, state.PeriodReward)
					}
					fmt.Println("-----------Change HalfFactor------------")
					fmt.Println("Day:", i, "Change HalfFactor to", state.HalfFactor, "reward:", WeiToMemo(reward, config.Decimals), "MintLevel:", state.MintLevel)
					fmt.Println("spacetime:", state.TotalSpaceTime.String(), "esize:", FormatGBytes(tempEsize.Int64()))
					fmt.Println("Calculate Dur:", float64(new(big.Int).Div(state.TotalSpaceTime, tempEsize).Uint64()))
					fmt.Println("TotalPay:", WeiToMemo(state.TotalPay, config.Decimals), "TotalPaid:", WeiToMemo(state.TotalPaid, config.Decimals), "TotalReward", WeiToMemo(state.TotalReward, config.Decimals))
					fmt.Println("TotalLiquid:", WeiToMemo(state.TotalLiquid, config.Decimals),
						"TotalSupply:", WeiToMemo(state.TotalSupply, config.Decimals),
						"Percent:", float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()))
					fmt.Println("TotalSize:", FormatGBytes(state.TotalSize.Int64()))
					fmt.Println("OrderSize:", FormatGBytes(order.Size.Int64()), "OrderPrice", order.Price)
					fmt.Println("TargetReward:", WeiToMemo(state.TargetReward, config.Decimals), "PeriodRewad:", WeiToMemo(state.PeriodReward, config.Decimals))
					fmt.Println("Issurance ratio:", float64(state.Ratio.Int64())/float64(OneBillion.Int64()), "Providers Count", state.ProviderCount)
					fmt.Println("Total issue times:", float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64())/float64(WeiToMemo(config.InitialSupply, config.Decimals).Int64()))
					fmt.Println("TotalPledge:", WeiToMemo(state.TotalPledge, config.Decimals))
				}
			}
		}

		// 填充纵轴数据
		pts[SUPPLY_INDEX][i].Y = float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64())
		pts[LIQUID_INDEX][i].Y = float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())
		pts[REWARD_INDEX][i].Y = float64(WeiToMemo(state.TotalReward, config.Decimals).Int64())
		pts[PLEDGE_INDEX][i].Y = float64(WeiToMemo(state.TotalPledge, config.Decimals).Uint64())
		pts[PAY_INDEX][i].Y = float64(WeiToMemo(state.TotalPay, config.Decimals).Uint64())
		pts[PAID_INDEX][i].Y = float64(WeiToMemo(state.TotalPaid, config.Decimals).Uint64())
		pts[SIZE_INDEX][i].Y = float64(state.TotalSize.Int64())

		// 模拟下一次订单变化

		// 新增Provider数
		order.NewProvider, state.ProviderPledge = config.ProviderSimulate(state, reward, int64(i), config, order)

		providerPledge := new(big.Int).Mul(big.NewInt(order.NewProvider), state.ProviderPledge)
		state.TotalPledge.Add(
			state.TotalPledge,
			providerPledge,
		)
		state.ProviderCount += order.NewProvider
		state.TotalLiquid.Sub(state.TotalLiquid, providerPledge)

		// 下一天的订单空间数
		order.Size = config.SizeSimulate(state, reward, int64(i), config, order)
		// 下一天的订单价格
		order.Price = config.PriceSimulate(state, reward, int64(i), config, order)
		// 下一天的订单时长
		order.Dur = config.DurationSimulate(state, reward, int64(i), config, order)

		// 打印一些时间点的数据
		if i%200 >= 0 && i%200 < 3 || int64(i) >= config.TotalDuration-5 {
			fmt.Println("--------Specific---------------")
			fmt.Println("Day:", i, "reward:", WeiToMemo(reward, config.Decimals), "MintLevel:", state.MintLevel, "HalfFactor", state.HalfFactor)
			fmt.Println("spacetime:", state.TotalSpaceTime.String(), "esize:", FormatGBytes(tempEsize.Int64()))
			fmt.Println("Calculate Dur:", float64(new(big.Int).Div(state.TotalSpaceTime, tempEsize).Uint64()))
			fmt.Println("TotalPay:", WeiToMemo(state.TotalPay, config.Decimals), "TotalPaid:", WeiToMemo(state.TotalPaid, config.Decimals), "TotalReward", WeiToMemo(state.TotalReward, config.Decimals))
			fmt.Println("TotalLiquid:", WeiToMemo(state.TotalLiquid, config.Decimals),
				"TotalSupply:", WeiToMemo(state.TotalSupply, config.Decimals),
				"Percent:", float64(WeiToMemo(state.TotalLiquid, config.Decimals).Int64())/float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64()))
			fmt.Println("TotalSize:", FormatGBytes(state.TotalSize.Int64()))
			fmt.Println("OrderSize:", FormatGBytes(order.Size.Int64()), "OrderPrice", order.Price)
			fmt.Println("TargetReward:", WeiToMemo(state.TargetReward, config.Decimals), "PeriodRewad:", WeiToMemo(state.PeriodReward, config.Decimals))
			fmt.Println("Issurance ratio:", float64(state.Ratio.Int64())/float64(OneBillion.Int64()), "Providers Count", state.ProviderCount)
			fmt.Println("Total issue times:", float64(WeiToMemo(state.TotalSupply, config.Decimals).Int64())/float64(WeiToMemo(config.InitialSupply, config.Decimals).Int64()))
			fmt.Println("TotalPledge:", WeiToMemo(state.TotalPledge, config.Decimals))
		}

		state.LastMint = start
	}

	return pts
}
