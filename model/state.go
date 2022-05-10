package model

import "math/big"

type NodeState struct {
	Start uint64
	Size  uint64
}

type MemoState struct {
	day uint64

	cfg *Config

	liquid *big.Int // 总当前可流动代币数

	paid   *big.Int // pay for order
	pledge *big.Int
	profit *big.Int

	foundation *big.Int
	kincome    *big.Int
	pincome    *big.Int

	mint        *MintInfo
	spaceTime   *big.Int
	spacePrice  *big.Int
	size        *big.Int
	reward      *big.Int            // reward to pledge pool
	subPriceMap map[uint64]*big.Int // 记录过期价格，key为天数
	subSizeMap  map[uint64]*big.Int // 记录过期空间，key为天数

	fixPledge     *big.Int
	groups        uint64
	keeperCount   uint64 // 系统中Keeper的数量
	providerCount uint64 // 系统中Provider的数量
}

func NewMemoState(config *Config) *MemoState {
	s := &MemoState{
		cfg:       config,
		liquid:    new(big.Int).Set(config.InitSupply),
		paid:      big.NewInt(0),
		pledge:    big.NewInt(0),
		profit:    big.NewInt(0),
		fixPledge: big.NewInt(0),

		foundation: big.NewInt(0),
		kincome:    big.NewInt(0),
		pincome:    big.NewInt(0),

		mint:        InitMint(config.RewardTarget),
		spaceTime:   big.NewInt(0),
		spacePrice:  big.NewInt(0),
		size:        big.NewInt(0),
		reward:      big.NewInt(0),
		subPriceMap: make(map[uint64]*big.Int),
		subSizeMap:  make(map[uint64]*big.Int),

		keeperCount:   0,
		providerCount: 0,
	}

	return s
}
