package model

import (
	"fmt"
	"math/big"
)

type MintInfo struct {
	mc       MintConfig
	Residual *big.Int
	Ratio    *big.Int // 增发比例
	Size     *big.Int // 期望空间
	Reward   *big.Int // 本阶段奖励
}

func InitMint(mc MintConfig) *MintInfo {
	stageRatio := big.NewInt(mc.RatioInit * mc.RatioDecimal)
	stageSize := big.NewInt(TiB)
	stageReward := new(big.Int).Mul(stageSize, stageRatio)
	stageReward.Mul(stageReward, big.NewInt(MaxDuration*Day))
	stageReward.Div(stageReward, big.NewInt(mc.RatioDecimal))
	mi := &MintInfo{
		mc:       mc,
		Residual: new(big.Int).Mul(big.NewInt(mc.RewardTarget), big.NewInt(Memo)),
		Ratio:    new(big.Int).Set(stageRatio),
		Size:     new(big.Int).Set(stageSize),
		Reward:   new(big.Int).Set(stageReward),
	}

	return mi
}

func (m *MintInfo) Check(nsize, dur *big.Int) {
	if nsize.Cmp(m.Size) < 0 {
		return
	}

	fmt.Println(m.Ratio, new(big.Int).Div(m.Size, big.NewInt(TiB)), new(big.Int).Div(dur, big.NewInt(Day)), WeiToMemo(m.Reward), WeiToMemo(m.Residual))

	stageRatio := new(big.Int).Set(m.Ratio)
	stageSize := new(big.Int).Set(m.Size)
	stageReward := new(big.Int).Set(m.Reward)

	if stageSize.Cmp(big.NewInt(EiB)) >= 0 {
		stageSize = stageSize.Mul(stageSize, big.NewInt(2))
	} else {
		stageSize = stageSize.Mul(stageSize, big.NewInt(32))
	}

	stageRatio.Mul(stageRatio, big.NewInt(m.mc.RatioReduce))
	stageRatio.Div(stageRatio, big.NewInt(m.mc.RatioReduceDecimal))

	stageReward.Mul(stageSize, stageRatio)
	stageReward.Mul(stageReward, dur)
	stageReward.Div(stageReward, big.NewInt(m.mc.RatioDecimal))
	if stageReward.Cmp(m.Residual) >= 0 {
		stageRatio.Mul(m.Residual, big.NewInt(m.mc.RatioDecimal))
		stageRatio.Div(stageRatio, stageSize)
		stageRatio.Div(stageRatio, dur)
		stageRatio.Div(stageRatio, big.NewInt(2)) // half of rest

		stageReward.Mul(stageSize, stageRatio)
		stageReward.Mul(stageReward, dur)
		stageReward.Div(stageReward, big.NewInt(m.mc.RatioDecimal))
	}

	m.Ratio.Set(stageRatio)
	m.Reward.Set(stageReward)
	m.Size.Set(stageSize)

	fmt.Println(m.Ratio, new(big.Int).Div(m.Size, big.NewInt(TiB)), new(big.Int).Div(dur, big.NewInt(Day)), WeiToMemo(m.Reward), WeiToMemo(m.Residual))
}
