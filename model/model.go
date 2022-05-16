package model

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"
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
	// create a new group when:
	// 1. is first group
	// 2. provider count of current group has exceed 70%;
	// 3. provider average storage > 80% of proivderStorage

	stop := false
	if s.groups >= s.cfg.Role.GroupCount {
		// not add group
		stop = true
	}

	create := false
	ts := new(big.Int).Mul(big.NewInt(s.cfg.Role.ProviderStorage), big.NewInt(int64(s.providerCount)))
	ts.Mul(ts, big.NewInt(4))
	ts.Div(ts, big.NewInt(5))
	if s.size.Cmp(ts) > 0 {
		create = true
	}

	gs, ok := s.gState[s.groups]
	if !stop && (s.groups == 0 || create || (ok && gs.PCnt >= s.cfg.Role.ProviderCntPerGroup*70/100)) {
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

	// create provider
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
			s.fs.Add(s.fs, pay)

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

// update income for providers and keepers
func (s *MemoState) updateIncome() {
	// linear pay
	income := new(big.Int).Mul(s.spacePrice, big.NewInt(Day))
	s.pincome.Add(s.pincome, income)
	s.fs.Sub(s.fs, income)

	income.Mul(income, big.NewInt(s.cfg.Order.LinearRate))
	income.Div(income, big.NewInt(100))
	s.kincome.Add(s.kincome, income)

	ss, ok := s.subSizeMap[s.day]
	if ok {
		s.size.Sub(s.size, ss)
	}

	sp, ok := s.subPriceMap[s.day]
	if ok {
		// endpay
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

// update pledge amount depend on profit
func (s *MemoState) updatePledge() {
	if s.pledge.BitLen() == 0 {
		return
	}

	profit := new(big.Int)

	cnt := int64(0)
	for i := s.day; i > 0; i-- {
		profit.Add(profit, s.profits[i])
		cnt++
		if cnt >= s.cfg.Pledge.DayBack {
			break
		}
	}

	if cnt > 0 {
		profit.Div(profit, big.NewInt(cnt))
	}

	// 年化收益
	profit.Mul(profit, big.NewInt(36500))

	// profit > 1% per day, pledge more
	if new(big.Int).Div(profit, s.pledge).Cmp(big.NewInt(s.cfg.Pledge.InRatio)) > 0 {
		// pledge more
		pt := new(big.Int).Div(s.pledge, big.NewInt(100))
		s.pledge.Add(s.pledge, pt)
		s.liquid.Sub(s.liquid, pt)
		if s.cfg.Simu.Detail {
			fmt.Println("pledge: ", WeiToMemo(pt))
		}

		return
	}

	// profit < 0.25%, withdraw
	if new(big.Int).Div(profit, s.pledge).Cmp(big.NewInt(s.cfg.Pledge.OutRatio)) < 0 {

		pt := new(big.Int).Div(s.pledge, big.NewInt(100))

		lt := new(big.Int).Sub(s.pledge, s.fixPledge)

		// withdraw
		if pt.Cmp(lt) > 0 {
			pt.Set(lt)
		}

		if pt.BitLen() == 0 {
			return
		}

		s.liquid.Add(s.liquid, pt)
		s.pledge.Sub(s.pledge, pt)
		if s.cfg.Simu.Detail {
			fmt.Println("withdraw: ", WeiToMemo(pt))
		}
	}
}

// check mint level
func (s *MemoState) checkMint() {
	dur := new(big.Int).Div(s.spaceTime, s.accSize) // average duration
	//nsize := new(big.Int).Div(s.spaceTime, big.NewInt(MinDuration*Day))
	nsize := new(big.Int).Set(s.accSize)
	s.mint.Check(nsize, dur)
}

// unlock memo
func (s *MemoState) updateLiquid() {
	// linear unlock
	if s.day < uint64(s.cfg.Token.LinearDay) {
		s.liquid.Add(s.liquid, s.unlockPerDay)
	}

	// unlock at some day
	if s.day == uint64(s.cfg.Token.LockDay) {
		uv := new(big.Int).Mul(big.NewInt(s.cfg.Token.LockSupply), big.NewInt(Memo))
		s.liquid.Add(s.liquid, uv)
	}
}

// plot
var PlotX []string
var PlotData [][]opts.LineData

func Simulate(cfg *Config) {
	ss := time.Now()

	s := NewMemoState(cfg)

	PlotX = make([]string, cfg.Simu.Duration)
	PlotData = make([][]opts.LineData, POINTS_COUNT)

	for i := 0; i < POINTS_COUNT; i++ {
		PlotData[i] = make([]opts.LineData, cfg.Simu.Duration)
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

		s.updateLiquid()

		cnt := int64(0)
		profit := new(big.Int)
		for i := s.day; i > 0; i-- {
			profit.Add(profit, s.profits[i])
			cnt++
			if cnt >= s.cfg.Pledge.DayBack {
				break
			}
		}

		if cnt > 0 {
			profit.Div(profit, big.NewInt(cnt))
		}

		// 年化收益
		profit.Mul(profit, big.NewInt(36500))
		profit.Div(profit, s.pledge)

		if s.cfg.Simu.Detail {
			fmt.Println(s.day, s.groups, s.providerCount, ",liquid:", WeiToMemo(s.liquid), ",pledge:", WeiToMemo(s.pledge), ",reward:", WeiToMemo(s.reward), ",yearly: ", profit, ",paid:", WeiToMemo(s.paid), ",fs:", WeiToMemo(s.fs), ",income:", WeiToMemo(s.pincome), ",kincome:", WeiToMemo(s.kincome), ",size:", new(big.Int).Div(s.size, big.NewInt(TiB)), new(big.Int).Div(s.accSize, big.NewInt(TiB)), time.Since(nt))
		}

		PlotX[i] = strconv.Itoa(int(i))

		liquid := new(big.Int).Add(s.liquid, s.pincome)
		liquid.Add(liquid, s.kincome)

		PlotData[SUPPLY_INDEX][i].Value = s.cfg.Token.TotalSupply
		PlotData[LIQUID_INDEX][i].Value = int(WeiToMemo(liquid).Int64())
		PlotData[REWARD_INDEX][i].Value = int(WeiToMemo(s.reward).Int64())
		PlotData[PLEDGE_INDEX][i].Value = WeiToMemo(s.pledge).Uint64()
		PlotData[PAID_INDEX][i].Value = WeiToMemo(s.paid).Uint64()
		PlotData[SIZE_INDEX][i].Value = new(big.Int).Div(s.size, big.NewInt(TiB)).Int64()
		PlotData[ASIZE_INDEX][i].Value = new(big.Int).Div(s.accSize, big.NewInt(TiB)).Int64()
	}

	fmt.Println("simulate cost:", time.Since(ss))

	fmt.Println(s.day, s.groups, s.providerCount, ",liquid:", WeiToMemo(s.liquid), ",pledge:", WeiToMemo(s.pledge), ",reward:", WeiToMemo(s.reward), ",paid:", WeiToMemo(s.paid), ",income:", WeiToMemo(s.pincome), ",kincome:", WeiToMemo(s.kincome), ",size:", new(big.Int).Div(s.size, big.NewInt(TiB)), new(big.Int).Div(s.accSize, big.NewInt(TiB)))

	fmt.Println(s.mint.Ratio, new(big.Int).Div(s.mint.Size, big.NewInt(TiB)), WeiToMemo(s.mint.Reward), WeiToMemo(s.mint.Residual))
}
