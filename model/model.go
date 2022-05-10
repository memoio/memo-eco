package model

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"gonum.org/v1/plot/plotter"
)

const (
	POINTS_COUNT = 7

	SUPPLY_INDEX = 0
	LIQUID_INDEX = 1
	REWARD_INDEX = 2
	PLEDGE_INDEX = 3

	PAID_INDEX = 5
	SIZE_INDEX = 6
)

// assume each provider has 4TB, one group has 5000 provider
func (s *MemoState) updateGroup() {
	if s.day%30 == 0 {
		s.groups++
		s.keeperCount += KCntPerGroup
		kp := new(big.Int).Mul(s.cfg.KeeperPledge, big.NewInt(KCntPerGroup))
		s.fixPledge.Add(s.fixPledge, kp)
		s.pledge.Add(s.pledge, kp)
		s.liquid.Sub(s.liquid, kp)
	}

	if s.providerCount < s.groups*PCntPerGroup {
		s.providerCount += ProCreate
		pp := new(big.Int).Mul(s.cfg.ProviderPledge, big.NewInt(ProCreate))
		s.fixPledge.Add(s.fixPledge, pp)
		s.pledge.Add(s.pledge, pp)
		s.liquid.Sub(s.liquid, pp)
	}
}

// each provider has one order
func (s *MemoState) updateOrder() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < int(s.providerCount); i++ {
		//dur := uint64(MinDuration + r.Int63n(4*MinDuration-MinDuration))
		dur := uint64(DefaultDuration)
		size := big.NewInt(DefaultSize + r.Int63n(10*DefaultSize-DefaultSize/10))
		price := int64(DefaultPrice)

		sprice := new(big.Int).Mul(size, big.NewInt(price))
		sprice.Div(sprice, big.NewInt(GiB))

		ss, ok := s.subSizeMap[s.day+dur]
		if ok {
			ss.Add(ss, size)
		} else {
			s.subSizeMap[s.day+dur] = size
		}

		sp, ok := s.subPriceMap[s.day+dur]
		if ok {
			sp.Add(sp, sprice)
		} else {
			s.subPriceMap[s.day+dur] = sprice
		}

		s.size.Add(s.size, size)
		s.spacePrice.Add(s.spacePrice, sprice)

		st := new(big.Int).SetUint64(dur)
		st.Mul(st, size)
		s.spaceTime.Add(s.spaceTime, st)

		pay := new(big.Int).Mul(sprice, new(big.Int).SetUint64(dur))
		s.paid.Add(s.paid, pay)

		// 1% for tax
		pay.Div(pay, big.NewInt(100))
		s.foundation.Add(s.foundation, pay)

		// %5 for manage and tax
		pay.Mul(pay, big.NewInt(100+LinearRate+EndRate+TaxRate))
		s.liquid.Sub(s.liquid, pay)
	}
}

func (s *MemoState) updateReward() {
	reward := new(big.Int).Mul(s.spacePrice, big.NewInt(Day))
	reward.Mul(reward, s.mint.Ratio)
	reward.Div(reward, big.NewInt(RatioDecimal))

	if reward.Cmp(s.mint.Residual) > 0 {
		fmt.Println("reward not enough:", WeiToMemo(reward), WeiToMemo(s.mint.Residual))
		reward.Set(s.mint.Residual)
	}

	if reward.Cmp(s.mint.Reward) > 0 {
		fmt.Println("stage reward not enough:", WeiToMemo(reward), WeiToMemo(s.mint.Reward))
		reward.Set(s.mint.Reward)
	}

	s.mint.Reward.Sub(s.mint.Reward, reward)
	s.mint.Residual.Sub(s.mint.Residual, reward)

	s.profit.Set(reward)

	s.reward.Add(s.reward, reward)
}

func (s *MemoState) updateIncome() {
	// update income for provider and keeper
	income := new(big.Int).Mul(s.spacePrice, big.NewInt(Day))
	s.pincome.Add(s.pincome, income)
	income.Mul(income, big.NewInt(LinearRate))
	income.Div(income, big.NewInt(100))
	s.kincome.Add(s.kincome, income)

	ss, ok := s.subSizeMap[s.day]
	if ok {
		s.size.Sub(s.size, ss)
	}

	sp, ok := s.subPriceMap[s.day]
	if ok {
		s.spacePrice.Sub(s.spacePrice, sp)
		income.Mul(sp, big.NewInt(Day))
		income.Mul(income, big.NewInt(EndRate))
		income.Div(income, big.NewInt(100))
		s.kincome.Add(s.kincome, income)
	}
}

// depend on profit
func (s *MemoState) updatePledge() {
	pt := WeiToMemo(s.pledge)

	if pt.BitLen() == 0 {
		return
	}

	// profit > 1% per day
	if new(big.Int).Div(s.profit, pt).Cmp(big.NewInt(Memo/100)) > 0 {
		for {
			pt.Mul(pt, big.NewInt(11))
			pt.Div(pt, big.NewInt(10))
			if new(big.Int).Div(s.profit, pt).Cmp(big.NewInt(Memo/100)) < 0 {
				pt.Mul(pt, big.NewInt(Memo))
				if pt.Cmp(s.pledge) > 0 {
					// pledge more
					pt.Sub(pt, s.pledge)
					s.pledge.Add(s.pledge, pt)
					s.liquid.Sub(s.liquid, pt)
					fmt.Println("pledge: ", WeiToMemo(pt))
				}
				break
			}
		}
	}

	pt = WeiToMemo(s.pledge)

	// profit < 0.25%
	if new(big.Int).Div(s.profit, pt).Cmp(big.NewInt(Memo/400)) < 0 {
		for {
			pt.Mul(pt, big.NewInt(9))
			pt.Div(pt, big.NewInt(10))
			if pt.BitLen() == 0 {
				return
			}
			if new(big.Int).Div(s.profit, pt).Cmp(big.NewInt(Memo/400)) > 0 {
				if pt.Cmp(s.pledge) < 0 {
					// withdraw
					if pt.Cmp(s.fixPledge) < 0 {
						pt.Set(s.fixPledge)
					}

					pt.Sub(s.pledge, pt)
					s.liquid.Add(s.liquid, pt)
					s.pledge.Sub(s.pledge, pt)
					fmt.Println("withdraw: ", WeiToMemo(pt))
				}

				break
			}
		}
	}
}

func (s *MemoState) checkMint() {
	nsize := new(big.Int).Div(s.spaceTime, big.NewInt(MinDuration))
	s.mint.Check(nsize)
}

// 生成每天的模拟数据
func EcoModelSimulate(config *Config) []plotter.XYs {
	s := NewMemoState(config)

	// 计算这么多天的增发
	pts := make([]plotter.XYs, POINTS_COUNT)

	for i := 0; i < len(pts); i++ {
		pts[i] = make(plotter.XYs, config.TotalDuration)
	}

	fmt.Println("Memo initial supply:", WeiToMemo(s.cfg.TotalSupply))
	fmt.Println("Memo initial liquid:", WeiToMemo(s.liquid))

	// 代币总量数据
	pts[SUPPLY_INDEX][0].X = 0
	// 单位换算成Memo
	pts[SUPPLY_INDEX][0].Y = float64(WeiToMemo(s.cfg.TotalSupply).Uint64())

	// 流通代币数据
	pts[LIQUID_INDEX][0].X = 0
	// 单位换算成Memo
	pts[LIQUID_INDEX][0].Y = float64(WeiToMemo(s.liquid).Uint64())

	// 奖励数据
	pts[REWARD_INDEX][0].X = 0
	pts[REWARD_INDEX][0].Y = 0

	// 质押数据
	pts[PLEDGE_INDEX][0].X = 0
	pts[PLEDGE_INDEX][0].Y = 0

	// 已支付数据
	pts[PAID_INDEX][0].X = 0
	pts[PAID_INDEX][0].Y = 0

	// 存储空间数据
	pts[SIZE_INDEX][0].X = 0
	pts[SIZE_INDEX][0].Y = 0

	// 开始模拟每天的订单
	for i := uint64(0); i < uint64(s.cfg.TotalDuration); i++ {
		s.day = i

		s.updateGroup()
		s.updateOrder()
		s.updateReward()
		s.updateIncome()

		s.checkMint()
		s.updatePledge()

		dp := new(big.Int).Mul(s.profit, big.NewInt(10000))
		dp.Div(dp, s.pledge)

		fmt.Println(s.day, s.groups, ",liquid:", WeiToMemo(s.liquid), ",pledge:", WeiToMemo(s.pledge), ",reward:", WeiToMemo(s.reward), ",daily: ", dp, ",paid:", WeiToMemo(s.paid), ",income:", WeiToMemo(s.pincome), ",kincome:", WeiToMemo(s.kincome), ",size:", new(big.Int).Div(s.size, big.NewInt(TiB)))

		// 填充纵轴数据
		pts[SUPPLY_INDEX][i].Y = float64(WeiToMemo(s.cfg.TotalSupply).Int64())
		pts[REWARD_INDEX][i].Y = float64(WeiToMemo(s.reward).Int64())
		pts[PLEDGE_INDEX][i].Y = float64(WeiToMemo(s.pledge).Uint64())
		pts[PAID_INDEX][i].Y = float64(WeiToMemo(s.paid).Uint64())
		pts[SIZE_INDEX][i].Y = float64(s.size.Int64())
	}

	return pts
}
