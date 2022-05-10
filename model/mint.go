package model

import (
	"fmt"
	"math/big"
)

type MintInfo struct {
	Residual *big.Int
	Ratio    *big.Int // 增发比例
	Size     *big.Int // 期望空间
	Reward   *big.Int // 本阶段奖励
}

func InitMint(rt *big.Int) *MintInfo {
	stageRatio := big.NewInt(RatioInit * RatioDecimal)
	stageSize := big.NewInt(TiB)
	stageReward := new(big.Int).Mul(stageSize, stageRatio)
	stageReward.Mul(stageReward, big.NewInt(MinDuration))
	stageReward.Div(stageReward, big.NewInt(RatioDecimal))
	mi := &MintInfo{
		Residual: new(big.Int).Set(rt),
		Ratio:    new(big.Int).Set(stageRatio),
		Size:     new(big.Int).Set(stageSize),
		Reward:   new(big.Int).Set(stageReward),
	}

	return mi
}

func (m *MintInfo) Check(nsize *big.Int) {
	if nsize.Cmp(m.Size) < 0 {
		return
	}

	stageRatio := new(big.Int).Set(m.Ratio)
	stageSize := new(big.Int).Set(m.Size)
	stageReward := new(big.Int).Set(m.Reward)
	rt := new(big.Int).Set(m.Residual)

	tmpSize := new(big.Int).Set(stageSize)
	if stageSize.Cmp(big.NewInt(EiB)) >= 0 {
		stageSize = stageSize.Mul(stageSize, big.NewInt(2))
	} else {
		stageSize = stageSize.Mul(stageSize, big.NewInt(32))
	}

	stageRatio.Mul(stageRatio, big.NewInt(RatioReduce))
	stageRatio.Div(stageRatio, big.NewInt(RatioReduceDecimal))

	tmpSize.Sub(stageSize, tmpSize)
	stageReward.Mul(tmpSize, stageRatio)
	stageReward.Mul(stageReward, big.NewInt(MinDuration))
	stageReward.Div(stageReward, big.NewInt(RatioDecimal))
	if stageReward.Cmp(rt) >= 0 {
		stageRatio.Mul(rt, big.NewInt(RatioDecimal))
		stageRatio.Div(stageRatio, tmpSize)
		stageRatio.Div(stageRatio, big.NewInt(MinDuration*2)) // half of rest

		stageReward.Mul(tmpSize, stageRatio)
		stageReward.Mul(stageReward, big.NewInt(MinDuration))
		stageReward.Div(stageReward, big.NewInt(RatioDecimal))
	}

	m.Ratio.Set(stageRatio)
	m.Reward.Add(m.Reward, stageReward)
	m.Size.Set(stageSize)

	fmt.Println(stageRatio, new(big.Int).Div(stageSize, big.NewInt(TiB)), new(big.Int).Div(stageReward, big.NewInt(Memo)), new(big.Int).Div(rt, big.NewInt(Memo)))
}

func Mint(rt *big.Int) []MintInfo {
	stageRatio := big.NewInt(RatioInit * RatioDecimal)
	stageSize := big.NewInt(TiB)
	stageReward := new(big.Int).Mul(stageSize, stageRatio)
	stageReward.Mul(stageReward, big.NewInt(MinDuration))
	stageReward.Div(stageReward, big.NewInt(RatioDecimal))

	mint := make([]MintInfo, 0, 32)

	for i := 0; i < 32; i++ {
		mi := MintInfo{
			Ratio:  new(big.Int).Set(stageRatio),
			Size:   new(big.Int).Set(stageSize),
			Reward: new(big.Int).Set(stageReward),
		}
		mint = append(mint, mi)
		rt.Sub(rt, stageReward)

		fmt.Println(i, stageRatio, new(big.Int).Div(stageSize, big.NewInt(TiB)), new(big.Int).Div(stageReward, big.NewInt(Memo)), new(big.Int).Div(rt, big.NewInt(Memo)))

		tmpSize := new(big.Int).Set(stageSize)
		if stageSize.Cmp(big.NewInt(EiB)) >= 0 {
			stageSize = stageSize.Mul(stageSize, big.NewInt(2))
		} else {
			stageSize = stageSize.Mul(stageSize, big.NewInt(32))
		}

		stageRatio.Mul(stageRatio, big.NewInt(RatioReduce))
		stageRatio.Div(stageRatio, big.NewInt(RatioReduceDecimal))

		tmpSize.Sub(stageSize, tmpSize)
		stageReward.Mul(tmpSize, stageRatio)
		stageReward.Mul(stageReward, big.NewInt(MinDuration))
		stageReward.Div(stageReward, big.NewInt(RatioDecimal))
		if stageReward.Cmp(rt) >= 0 {
			stageRatio.Mul(rt, big.NewInt(RatioDecimal))
			stageRatio.Div(stageRatio, tmpSize)
			stageRatio.Div(stageRatio, big.NewInt(MinDuration*2)) // half of rest

			stageReward.Mul(tmpSize, stageRatio)
			stageReward.Mul(stageReward, big.NewInt(MinDuration))
			stageReward.Div(stageReward, big.NewInt(RatioDecimal))
		}
	}

	return mint
}
