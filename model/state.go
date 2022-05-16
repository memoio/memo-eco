package model

import (
	"math/big"
	"math/rand"
	"time"
)

type NodeState struct {
	Born  uint64
	Group uint64
	Size  uint64
}

type GroupState struct {
	Born       uint64
	Index      uint64
	KCnt       uint64
	PCnt       uint64
	Size       *big.Int
	subSizeMap map[uint64]*big.Int
}

type MemoState struct {
	r   *rand.Rand
	day uint64

	cfg *Config

	liquid       *big.Int // 总当前可流动代币数
	unlockPerDay *big.Int

	paid   *big.Int // pay for order
	fs     *big.Int
	pledge *big.Int

	foundation *big.Int
	kincome    *big.Int
	pincome    *big.Int

	mint        *MintInfo
	spaceTime   *big.Int
	spacePrice  *big.Int
	accSize     *big.Int
	size        *big.Int
	reward      *big.Int // reward to pledge pool
	profits     []*big.Int
	subPriceMap map[uint64]*big.Int // 记录过期价格，key为天数
	subSizeMap  map[uint64]*big.Int // 记录过期空间，key为天数

	fixPledge     *big.Int
	groups        uint64
	gState        map[uint64]*GroupState
	keeperCount   uint64 // 系统中Keeper的数量
	providerCount uint64 // 系统中Provider的数量
}

func NewMemoState(cfg *Config) *MemoState {
	upd := new(big.Int).Mul(big.NewInt(cfg.Token.LinearSupply), big.NewInt(Memo))
	upd.Div(upd, big.NewInt(cfg.Token.LinearDay))
	s := &MemoState{
		r:            rand.New(rand.NewSource(time.Now().UnixNano())),
		cfg:          cfg,
		liquid:       new(big.Int).Mul(big.NewInt(cfg.Token.InitSupply), big.NewInt(Memo)),
		unlockPerDay: upd,
		paid:         big.NewInt(0),
		fs:           big.NewInt(0),
		pledge:       big.NewInt(0),
		fixPledge:    big.NewInt(0),

		foundation: big.NewInt(0),
		kincome:    big.NewInt(0),
		pincome:    big.NewInt(0),

		mint:        InitMint(cfg.Mint),
		spaceTime:   big.NewInt(0),
		spacePrice:  big.NewInt(0),
		accSize:     big.NewInt(0), // acc size
		size:        big.NewInt(0), // current size
		reward:      big.NewInt(0),
		subPriceMap: make(map[uint64]*big.Int),
		subSizeMap:  make(map[uint64]*big.Int),

		profits: make([]*big.Int, cfg.Simu.Duration),

		keeperCount:   0,
		providerCount: 0,
		gState:        make(map[uint64]*GroupState),
	}

	return s
}
