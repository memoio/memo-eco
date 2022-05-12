package model

import (
	"fmt"
	"math/big"
	"time"

	"gonum.org/v1/plot/plotter"
)

const (
	POINTS_COUNT = 7

	SUPPLY_INDEX = 0
	LIQUID_INDEX = 1
	REWARD_INDEX = 2
	PLEDGE_INDEX = 3
	PAID_INDEX   = 4
	SIZE_INDEX   = 5
	ASIZE_INDEX  = 6
)

func (s *MemoState) updateGroup() {
	// current group has exceed 70%; create a new group
	// provider average storage > 80% of proivderStorage
	create := false
	ts := new(big.Int).Mul(big.NewInt(s.cfg.Role.ProviderStorage), big.NewInt(int64(s.providerCount)))
	ts.Mul(ts, big.NewInt(4))
	ts.Div(ts, big.NewInt(5))

	if s.size.Cmp(ts) > 0 {
		create = true
	}

	gs, ok := s.gState[s.groups]
	if s.groups == 0 || create || (ok && gs.PCnt >= s.cfg.Role.ProviderCntPerGroup*70/100) {
		s.groups++

		knt := s.cfg.Role.KeeperCntPerGroup
		ngs := &GroupState{
			Born:       s.day,
			Index:      s.groups,
			KCnt:       knt,
			Size:       big.NewInt(0),
			subSizeMap: make(map[uint64]*big.Int),
		}
		s.gState[s.groups] = ngs

		s.keeperCount += knt
		kp := new(big.Int).Mul(big.NewInt(s.cfg.Role.KeeperPledge), big.NewInt(int64(knt)))
		kp.Mul(kp, big.NewInt(Memo))
		s.fixPledge.Add(s.fixPledge, kp)
		s.pledge.Add(s.pledge, kp)
		s.liquid.Sub(s.liquid, kp)
	}

	// 50-150
	for i := uint64(1); i <= s.groups; i++ {
		gs, ok := s.gState[i]
		if ok && gs.PCnt < s.cfg.Role.ProviderCntPerGroup {
			pn := s.cfg.Role.ProviderCreatePerDay - 2*(s.day-gs.Born) // todo: add create distribution
			if pn <= 0 {
				continue
			}
			gs.PCnt += pn

			s.providerCount += pn
			pp := new(big.Int).Mul(big.NewInt(s.cfg.Role.ProviderPledge), big.NewInt(int64(pn)))
			pp.Mul(pp, big.NewInt(Memo))
			s.fixPledge.Add(s.fixPledge, pp)
			s.pledge.Add(s.pledge, pp)
			s.liquid.Sub(s.liquid, pp)
		}
	}
}

// each provider has one order
// todo: size distribute, duration distribute
func (s *MemoState) updateOrder() {
	for j := uint64(1); j <= s.groups; j++ {
		gs, ok := s.gState[j]
		if !ok {
			continue
		}

		for i := uint64(0); i < gs.PCnt; i++ {

			durDay := uint64(MinDuration + s.r.Int63n(2*int64(s.cfg.Order.DefaultDuration)-MinDuration))
			size := big.NewInt(s.cfg.Order.DefaultSize/10 + s.r.Int63n(10*s.cfg.Order.DefaultSize-s.cfg.Order.DefaultSize/10))
			price := int64(s.cfg.Order.DefaultPrice)

			// add to group
			gs.Size.Add(gs.Size, size)
			gss, ok := gs.subSizeMap[s.day+durDay]
			if ok {
				gss.Add(gss, size)
			} else {
				gs.subSizeMap[s.day+durDay] = new(big.Int).Set(size)
			}

			ss, ok := s.subSizeMap[s.day+durDay]
			if ok {
				ss.Add(ss, size)
			} else {
				s.subSizeMap[s.day+durDay] = new(big.Int).Set(size)
			}

			sprice := new(big.Int).Mul(size, big.NewInt(price))
			sprice.Div(sprice, big.NewInt(GiB))
			s.spacePrice.Add(s.spacePrice, sprice)

			sp, ok := s.subPriceMap[s.day+durDay]
			if ok {
				sp.Add(sp, sprice)
			} else {
				s.subPriceMap[s.day+durDay] = new(big.Int).Set(sprice)
			}

			s.accSize.Add(s.accSize, size)
			s.size.Add(s.size, size)

			dur := new(big.Int).SetUint64(durDay * Day)

			st := new(big.Int).Mul(size, dur)
			s.spaceTime.Add(s.spaceTime, st)

			pay := new(big.Int).Mul(sprice, dur)

			// 1% for tax
			pay.Div(pay, big.NewInt(100))
			tax := new(big.Int).Mul(pay, big.NewInt(s.cfg.Order.TaxRate))
			s.foundation.Add(s.foundation, tax)

			// %5 for manage and tax
			pay.Mul(pay, big.NewInt(100+s.cfg.Order.LinearRate+s.cfg.Order.EndRate+s.cfg.Order.TaxRate))
			s.liquid.Sub(s.liquid, pay)
			s.paid.Add(s.paid, pay)

			ts := new(big.Int).Mul(big.NewInt(int64(gs.PCnt)), big.NewInt(s.cfg.Role.ProviderStorage))
			if gs.Size.Cmp(ts) >= 0 { // group is full, break
				break
			}
		}
	}
}

func (s *MemoState) updateReward() {
	reward := new(big.Int).Mul(s.spacePrice, big.NewInt(Day))
	reward.Mul(reward, s.mint.Ratio)
	reward.Div(reward, big.NewInt(s.cfg.Mint.RatioDecimal))

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

	s.reward.Add(s.reward, reward)

	s.profits[s.day] = new(big.Int).Set(reward)
}

func (s *MemoState) updateIncome() {
	// update income for provider and keeper
	income := new(big.Int).Mul(s.spacePrice, big.NewInt(Day))
	s.pincome.Add(s.pincome, income)
	income.Mul(income, big.NewInt(s.cfg.Order.LinearRate))
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
		income.Mul(income, big.NewInt(s.cfg.Order.EndRate))
		income.Div(income, big.NewInt(100))
		s.kincome.Add(s.kincome, income)
	}

	for j := uint64(1); j <= s.groups; j++ {
		gs, ok := s.gState[j]
		if ok {
			gss, ok := gs.subSizeMap[s.day]
			if ok {
				gs.Size.Sub(gs.Size, gss)
			}
		}
	}
}

// depend on profit
func (s *MemoState) updatePledge() {
	if s.pledge.BitLen() == 0 {
		return
	}

	profit := new(big.Int)
	pt := new(big.Int).Set(s.pledge)

	cnt := int64(0)
	for i := s.day; i > 0; i-- {
		profit.Add(profit, s.profits[i])
		cnt++
		if cnt >= 30 {
			break
		}
	}

	if cnt > 0 {
		profit.Div(profit, big.NewInt(cnt))
	}

	profit.Mul(profit, big.NewInt(10000))

	// profit > 1% per day, pledge more
	if new(big.Int).Div(profit, pt).Cmp(big.NewInt(s.cfg.Pledge.InRatio)) > 0 {
		for {
			pt.Mul(pt, big.NewInt(101))
			pt.Div(pt, big.NewInt(100))
			if new(big.Int).Div(profit, pt).Cmp(big.NewInt(s.cfg.Pledge.InRatio)) < 0 {
				if pt.Cmp(s.pledge) > 0 {
					// pledge more
					pt.Sub(pt, s.pledge)
					s.pledge.Add(s.pledge, pt)
					s.liquid.Sub(s.liquid, pt)
					if s.cfg.Simu.Detail {
						fmt.Println("pledge: ", WeiToMemo(pt))
					}
				}
				break
			}
		}

		return
	}

	// profit < 0.25%, withdraw
	if new(big.Int).Div(profit, pt).Cmp(big.NewInt(s.cfg.Pledge.OutRatio)) < 0 {
		for {
			pt.Mul(pt, big.NewInt(99))
			pt.Div(pt, big.NewInt(100))
			if pt.BitLen() == 0 {
				return
			}
			if new(big.Int).Div(profit, pt).Cmp(big.NewInt(s.cfg.Pledge.OutRatio)) > 0 {
				if pt.Cmp(s.pledge) < 0 {
					// withdraw
					if pt.Cmp(s.fixPledge) < 0 {
						pt.Set(s.fixPledge)
					}

					pt.Sub(s.pledge, pt)

					if pt.BitLen() == 0 {
						return
					}

					s.liquid.Add(s.liquid, pt)
					s.pledge.Sub(s.pledge, pt)
					if s.cfg.Simu.Detail {
						fmt.Println("withdraw: ", WeiToMemo(pt))
					}
				}

				break
			}
		}
	}
}

func (s *MemoState) checkMint() {
	dur := new(big.Int).Div(s.spaceTime, s.accSize) // average duration
	//nsize := new(big.Int).Div(s.spaceTime, big.NewInt(MinDuration*Day))
	nsize := new(big.Int).Set(s.accSize)
	s.mint.Check(nsize, dur)
}

// 生成每天的模拟数据
func Simulate(cfg *Config) []plotter.XYs {
	ss := time.Now()
	fmt.Println("============ simulate start ============")
	s := NewMemoState(cfg)

	// 计算这么多天的增发
	pts := make([]plotter.XYs, POINTS_COUNT)

	for i := 0; i < len(pts); i++ {
		pts[i] = make(plotter.XYs, cfg.Simu.Duration)
	}

	// 开始模拟每天的订单
	for i := uint64(0); i < cfg.Simu.Duration; i++ {
		nt := time.Now()

		s.day = i

		s.updateGroup()
		s.updateIncome()
		s.updateOrder()
		s.updateReward()

		s.checkMint()
		s.updatePledge()

		dp := new(big.Int).Mul(s.profits[s.day], big.NewInt(10000))
		dp.Div(dp, s.pledge)

		nt.Day()

		if s.cfg.Simu.Detail {
			fmt.Println(s.day, s.groups, s.providerCount, ",liquid:", WeiToMemo(s.liquid), ",pledge:", WeiToMemo(s.pledge), ",reward:", WeiToMemo(s.reward), ",daily: ", dp, ",paid:", WeiToMemo(s.paid), ",income:", WeiToMemo(s.pincome), ",kincome:", WeiToMemo(s.kincome), ",size:", new(big.Int).Div(s.size, big.NewInt(TiB)), new(big.Int).Div(s.accSize, big.NewInt(TiB)), time.Since(nt))
		}

		// 填充纵轴数据
		for j := 0; j < len(pts); j++ {
			pts[j][i].X = float64(i)
		}

		pts[SUPPLY_INDEX][i].Y = float64(s.cfg.Token.TotalSupply)
		pts[LIQUID_INDEX][i].Y = float64(WeiToMemo(s.liquid).Int64())
		pts[REWARD_INDEX][i].Y = float64(WeiToMemo(s.reward).Int64())
		pts[PLEDGE_INDEX][i].Y = float64(WeiToMemo(s.pledge).Uint64())
		pts[PAID_INDEX][i].Y = float64(WeiToMemo(s.paid).Uint64())
		pts[SIZE_INDEX][i].Y = float64(new(big.Int).Div(s.size, big.NewInt(GiB)).Int64())
		pts[ASIZE_INDEX][i].Y = float64(new(big.Int).Div(s.accSize, big.NewInt(GiB)).Int64())
	}

	fmt.Println("============ simulate end ============")

	fmt.Println("simulate cost:", time.Since(ss))

	fmt.Println(s.day, s.groups, s.providerCount, ",liquid:", WeiToMemo(s.liquid), ",pledge:", WeiToMemo(s.pledge), ",reward:", WeiToMemo(s.reward), ",paid:", WeiToMemo(s.paid), ",income:", WeiToMemo(s.pincome), ",kincome:", WeiToMemo(s.kincome), ",size:", new(big.Int).Div(s.size, big.NewInt(TiB)), new(big.Int).Div(s.accSize, big.NewInt(TiB)))

	fmt.Println(s.mint.Ratio, new(big.Int).Div(s.mint.Size, big.NewInt(TiB)), WeiToMemo(s.mint.Reward), WeiToMemo(s.mint.Residual))

	return pts
}
