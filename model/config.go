package model

// 配置参数
const (
	MinDuration = 100  //  day
	MaxDuration = 1000 // Day

)

type TokenConfig struct {
	TotalSupply  int64 // 创世代币数量
	InitSupply   int64 // unlock imediately
	LockSupply   int64 // unlock at some day
	LockDay      int64 // lock days
	LinearSupply int64 // unlock linearly
	LinearDay    int64 // unlock rate
}

type PledgeConfig struct {
	InRatio  int64
	OutRatio int64
}

type MintConfig struct {
	RewardTarget int64
	RatioInit    int64
	RatioDecimal int64
	RatioAlter   int64
}
type RoleConfig struct {
	KeeperPledge         int64
	ProviderPledge       int64
	KeeperCntPerGroup    uint64
	ProviderCntPerGroup  uint64
	ProviderStorage      int64
	ProviderCreatePerDay uint64
}

type OrderConfig struct {
	DefaultSize     int64 // default size
	DefaultDuration uint64
	DefaultPrice    uint64

	LinearRate int64
	EndRate    int64
	TaxRate    int64
}

type SimuConfig struct {
	Duration uint64
	Detail   bool
}

type Config struct {
	Simu   SimuConfig
	Token  TokenConfig
	Mint   MintConfig
	Role   RoleConfig
	Order  OrderConfig
	Pledge PledgeConfig
}

func DefaultConfig() *Config {
	return &Config{
		Simu: SimuConfig{
			Duration: 1000,
			Detail:   true,
		},
		Token: TokenConfig{
			TotalSupply:  600_000_000, // issue price: 1 $/Memo
			InitSupply:   100_000_000, // all is available or unlock linearly
			LockSupply:   100_000_000,
			LinearSupply: 100_000_000, // unlock linearly
			LockDay:      540,         // lock 540 day(18 month)
			LinearDay:    180,         // unlock in 180 day
		},
		Mint: MintConfig{
			RewardTarget: 300_000_000,
			RatioInit:    1, // per second*byte
			RatioDecimal: 1_000_000_000,
			RatioAlter:   150,
		},
		Role: RoleConfig{
			KeeperPledge:         1_000,   // Memo
			ProviderPledge:       100,     // Memo
			KeeperCntPerGroup:    10,      // each group has 10 keeper
			ProviderCntPerGroup:  5_000,   //assume each provider has 4TB, one group has 5000 provider
			ProviderStorage:      8 * TiB, // provider storage
			ProviderCreatePerDay: 150,     // one day
		},
		Order: OrderConfig{
			DefaultSize:     8 * GiB, // 100 KiB/second; 24 hours
			DefaultDuration: 365,     // day
			DefaultPrice:    1 * GiB, // 1 attoMemo per byte

			LinearRate: 3,
			EndRate:    1,
			TaxRate:    1,
		},
		Pledge: PledgeConfig{
			InRatio:  100,
			OutRatio: 25,
		},
	}
}

// todo: add config.json
func LoadConfig(path string) *Config {
	if path == "" {
		return DefaultConfig()
	}

	return nil
}
