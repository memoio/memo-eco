# Memo经济模型简版

[TOC]

# 1. Memo协议介绍

Memo协议是一个基于区块链技术的分散式存储基础设施。Memo结合了分层架构、高效的存储证明以及智能合约结算，从而构建了一个可靠、低成本且可扩展的存储服务。Memo致力于服务需要可靠数据存储的项目以及开发者，使他们可以简单地将存储基础替换为基于区块链的存储，让产品进一步去中心化。

从架构上看，Memo可以分成三个层次，分别为结算层、验证层以及存储层。

结算层将聚合所有的订单信息，并且订单金额逐步发放给存储节点，并且代币增发机制也实现在结算层之上，而这一切都由代码驱动。

验证层中的Keeper节点将挑战Provider节点，并验证证明结果，以决定是否给Provider节点取款凭证，验证层的所有流程都将经过验证层节点的拜占庭容错共识。

存储层的分散Provider节点存储真实的数据，并定期向验证层提交存储证明，同时使用取款凭证在结算层获取一个阶段的存储收益。

# 2. 基本设计

## 2.1 代币定义

Memo是Memo协议使用的代币，用于驱动整个Memo存储协议的运转以及规模增长。

在Memo协议中，Keeper和Provider需要首先质押Memo以获得赚取收益的资格，同时用户需要支付Memo代币以存储、检索数据，为了激励数据存储，根据订单的增长，Memo会增发分配给Keeper、Provider以及长期持币者。

Memo不仅是协议的使用代币，也是协议的权益代币，持有Memo的持币者将可以从Memo生态的增长中获益，同时参与Memo的治理。

## 2.2 设计原则

Memo代币旨在给用户以及生态参与方创造价值，提高生产力，并持续激励生态增长。Memo的经济模型将按照以下原则设计：

+ 参与方获得的激励应与其付出的成本成正比，并尽量保持公平。
+ 经济模型逻辑应尽量由代码控制，减少人为调控的部分。
+ 经济模型应鼓励参与方实施对协议有利的行为，满足激励相容。

## 2.3 代币分配

### 2.3.1 Memo经济循环

代币将在用户、开发者、Keeper以及存储节点之间循环流动，首先，用户需要购买Memo才能使用Memo的存储服务，Keeper和Provider需要先购买Memo质押才可以提供服务获取收益，而Memo的存储付费又会流动到Keeper和Provider中，同时，Keeper、Provider以及长期持币人也可以在质押池内质押Memo，从而获得Memo生态增长的激励增发。

由此，Memo的产生是一个消费再生产的循环，并将在此过程中逐步扩大存储规模，而规模的扩大，又会提高Memo的价值。

### 2.3.2 Memo使用场景

+ Keeper和Provider首先需要在质押池内deposit一定数量的Memo才可以参与协议提供服务。
+ 用户或者开发者需要购买Memo，并将其转移至订单池内，才可以签订订单使用Memo的存储服务，同时，检索数据也需要付出一定量的Memo。
+ 由于Memo同时还代表着对整个协议的权益，所以Memo的持有人可以将Memo质押到质押池，获取增发收益，同时，未来可能会设计治理机制，使用Memo代币可以参与到协议的治理过程。

### 2.3.3 代币创世发行分配

代币分配分为创世分配与激励增发分配，创世代币分配份额由以下组成：

+ 投资人: 29%，将在四百天内，每天线性解锁。
+ 核心团队: 27%，由核心团队控制，一半立刻释放，另一半在400天内线性解锁。
+ 合作伙伴及顾问: 14%，一部分将锁定在订单池，其他部分线性解锁。
+ 生态基金: 29%，暂时锁定，留待以后通过治理提取使用。

## 2.4 增发模型概述

在初始代币发行结束后，需要设计机制激励生态发展，尽快激励存储节点提供存储空间，同时激励用户存储数据。

首先，维护者（Keeper）和存储节点（Provider）需要质押Memo代币才能开始提供服务，而且任何持币人都可以在质押池质押Memo以获得增发收益。

除质押池外，订单池记录了每个存储节点和维护节点的订单金额，在每次用户与存储节点签订订单时，用户需要提前将代币金额转移到订单池中Provider以及Keeper节点的账户，但是节点不能立刻提取资金，只能在存储周期内按时上传存储证明，慢慢取走订单池内的资金。

发送给存储节点订单池的金额，在整个存储周期由存储节点逐步领取；发送给Keeper的，一部分随着存储节点领取收益而逐步发送给Keeper，另一部分在订单到期时，直接发送给Keeper。

此外在签订订单的时候，会根据订单信息以及参数计算一个增发量，然后将增发的代币发送给质押池，由所有在质押池内质押代币的账户按照份额平等分配奖励。

由上可知，Memo的代币增发将和Memo协议的使用情况绑定，当Memo协议规模扩大，才会增发更多的Memo以满足需求。Memo的增发将会在触发目标增发值后减半，所以早期的生态参与方可以获得更多红利。

### 2.4.1 质押机制

Keeper和Provider在提供服务前都需要质押，质押最低限度的主代币`MinimumKeeperPledge`以及`MinimumProviderPledge`，此外，任何持有主代币的用户都可以把自己的代币质押到质押池，以获得代币增发奖励。

### 2.4.2 质押池分润机制

质押池为以主代币GToken为质押代币的资金池，假设在T时刻，有个人$k_1$往质押池内质押主代币GToken $x_1$个，得到$x_1$个mGToken，此时池内GToken数量$\Zeta_{GToken}= x_1$，mGToken的数量$\Zeta_{mGToken}=x_1$。

在后续的系统运行中，假设系统得到了$\Delta_x$个GToken的奖励，则$\Zeta_{GToken1}= x_1+\Delta_x, \Zeta_{mGToken1}=x_1$，则此时$x_1$个mGToken可以兑换回$(x_1+\Delta_x)$个GToken，即每个mGToken值$\frac{x_1+\Delta_x}{x_1}$个GToken。

若在B时刻，另外一人$k_2$质押了$x_2$个GToken，则此时它将获得$\frac{\Zeta_{mGToken1}}{\Zeta_{GToken1}}x_2=\frac{x_1x_2}{x_1+\Delta_x}$个mGToken，此时$\Zeta_{GToken2}=x_1+x_2+\Delta_x,\Zeta_{mGToken2}=x_1+\frac{x_2x_1}{x_1+\Delta_x}$

若此时没有奖励，则此时$k_1$可以兑换回的代币数量为$x_1 \cdot \frac{x_1+x_2+\Delta_x}{x_1+\frac{x_2x_1}{x_1+\Delta_x}}=x_1+\Delta_x$，$k_2$可以兑回的数量为$\frac{x_1x_2}{x_1+\Delta_x}\cdot\frac{x_1+x_2+\Delta_x}{x_1+\frac{x_2x_1}{x_1+\Delta_x}}=x_2$个，与其贡献等同，而一旦有新的奖励打入池子，改变了$\Zeta_{GToken1}$，那么$k_2$将兑换得比$x_2$更多的GToken。

后面的新增质押以及奖励分配与上面类似。

### 2.4.3 增发机制

激励增发分配，包括长期的生态激励以及增发奖励，以激励维护者、存储节点、用户以及代币持有者为生态做贡献，这部分分配会随着订单量的增长逐步释放。

激励增发分配的首次目标释放为创世发行的一半，并在达到目标后减半，直至到达设定的最低增发率，保持稳定。

代币增发与存储量挂钩，并设定一些增发目标，在每次触发订单时，会根据当前的存储状态以及订单数据，设定一个增发量，增发到质押池，并由所有的质押者平分。

每个订单的参数包括订单开始时间`start`，订单结束时间`end`，订单大小`size`和订单单价`price`，时间单位为秒，订单大小单位为字节。

在处理一个订单时，计算当前的增发奖励，首先，存在一个累积当前总存储大小`totalSize`，将当前订单的Size考虑进去 $totalSize = totalSize + size$，计算整个订单的时空值$spaceTime=size\cdot (start-end)$，和时间价格值$spacePrice=price\cdot (start-end)$， 同时计算时间价格累积值$totalSpacePrice=totalSpacePrice+spacePrice$，以及订单支付总金额$pay=spaceTime \cdot price $。

同时，获取当前时间`nowTime`以及上一次触发增发的时间`lastMintTime`，计算$paid = (lastMintTime-nowTime)\cdot totalSpacePrice$，可得据上次增发后已支付了多少代币。

增发分为多个`mintLevel`，每个`mintLevel`包括的参数为增发比例`Ratio`，空间`Size`，周期长度`Duration`，其中越往后的`mintLevel`拥有更大的`size`。

然后计算订单所处的增发阶段，从第一个`mintLevel`开始，寻找最后一个符合以下条件的`mintLevel`：

+ 首先取出`mintLevel`设定的空间`eSize`，假如`eSize`小于累积空间，则令$eSize=totalSize$。
+ 然后取出`mintLevel`设定的周期长度`Dur`，计算平均订单周期长度$avarageDur=totalSpaceTime \div eSize$，假如$avarageDur>Dur$ ，则选取下一个`mintLevel`继续检查，直到不符合要求，选取上一个`mintLevel`为最终选择。
+ 然后计算本订单触发的增发奖励，取出`mintLevel`的增发比例`ratio`，计算$reward = paid \cdot ratio$，同时统计所有已增发的totalReward，如果已经超过本阶段目标reward，则触发减半，ratio需要除2重新计算reward，并将增发的代币转移给质押池。

问题：`mintLevel`内应该加入增发数量的属性，限制每个阶段可以增发的代币上限，超过后，自动跳到下一步。

# 3. 生态构成

## 3.1 基金会

Memo基金会主要用于生态系统建设、市场推广和社区维护等工作，同时付出部分资金用于投资促进生态发展，并保持基金会的长期可持续运行。

基金会Memo协议负有以下责任：

+ 组织开发团队或者外包任务以完成Memo协议的实现以及迭代升级。
+ 支持以及资助基于Memo开发的生态应用。
+ 长期参与Memo的社区运营，并持续贡献。

Memo基金会有权发起关于系统治理的提案，然后由社区决定该提案是否最终实施。基金会可以发起并包括但不限于以下建议:

+ 修改系统的经济参数
+ 提议对技术方案进行升级。
+ 惩罚作恶或不作为的Keeper节点。
+ 惩罚作恶或不作为的Provider节点。

## 3.2 社区生态

作为一个结合区块链的分散式去中心化存储方案，Memo协议的开发与发展离不开社区的支持。Memo基金会将积极组织和建立具有不同功能的社区，包括生态治理、开发者社区和持币人社区等，从而在多方面促进协议和生态健康稳定地发展。

## 3.3 Keeper

节点首先质押足额的保证金，然后才能在协议内成为Keeper节点。
作为Keeper节点，必须履行以下义务:

+ 股份不得少于规定金额的保证金。

+ 保证长期在线活动，并维护验证层的历史数据。

+ 检查Provider的存储证明，并通过拜占庭容错共识给Provider签发取款单。

+ 在Provider丢失数据时，调度数据修复流程。

同时，Keeper节点具有以下权益：

+ 从自己管理的订单内抽取一定比例作为收益。

+ 质押的保证金将获取协议的增发奖励。

如果Keeper节点作恶或者不作为，则可能会遭受惩罚，即丢失保证金，并且不能再继续提供服务获取收益。

## 3.4 Provider

拥有空闲存储空间的节点可以成为Provider，首先质押足额的保证金，然后可以与用户签订订单获取收益。
作为Provider，需要履行以下义务:

+ 在质押池质押足额的保证金。

+ 按照订单规定的标准存储数据，并保证数据的可靠性和可用性。

+ 按时向验证层提交存储证明。

同时，Provider具有以下权益

+ 逐步领取用户的订单付费。
+ 从质押池中获取增发收益。

如果Provider不履行义务，如丢失数据，不按时提交存储证明，则会遭受惩罚，扣除保证金以及转让订单。

## 3.5 开发者

开发者可以在自己的项目中使用Memo协议，更多的项目采用Memo的存储服务，会给Memo带来更多的用户和存储需求。

基金会可以通过一些开发者活动筛选对Memo协议有价值的开发者，并给予代币激励，并提供孵化服务。
开发人员需履行以下义务:

# 3. 经济模型仿真

## 3.1 参数选取

首先选取一个经济模型配置，如下所示：

```go
type MintInfo struct {
	Ratio    *big.Int // 增发比例
	Size     *big.Int // 期望空间
	Duration int64    // 期望周期
}

type EconomicsConfig struct {
	MintLevel             []MintInfo // 增发阶段设计
	Decimals              *big.Int   // Memo代币精度
	MinimumRation         *big.Int   // 最小增发率，减半到最小增发率后保持稳定
	InitialSupply         *big.Int   // 创世代币数量, 以 10^-8 Memo 为单位
	InitialTarget         *big.Int   // 创世第一阶段目标增发代币数, 以 10^-8 Memo 为单位，往后开始减半
	InitialKeeperPledge   *big.Int   // Keeper 初始需要质押的代币数
	InitialProviderPledge *big.Int   // Provider 初始需要质押的代币数

	InitialSize      *big.Int                 // 初始的订单空间
	InitialPrice     *big.Int                 // 初始的价格
	MinimumPrice     *big.Int                 // 最小的订单价格
	SizeSimulate     SizeSimulateFunction     // 模拟每天订单数据量大小的函数
	PriceSimulate    PriceSimulateFunction    // 模拟每天订单价格大小的函数
	DurationSimulate DurationSimulateFunction // 模拟每天订单的平均时间的函数
	ProviderSimulate ProviderSimulateFunction // 模拟每天新增Provider数据的函数
	EnableMaxSize    bool                     // 是否只有当前totalSize大于maxSize才增发？

	TotalDuration int64 // 总统计周期，单位 day
}
```

下面是一个示例的参数选取：

```go
func DefaultEconomicsConfig() *EconomicsConfig {
	decimals := big.NewInt(1_0000_0000)
	return &EconomicsConfig{
		MintLevel: []MintInfo{
			{
				Ratio:    big.NewInt(5_0000_0000), // 增发比例 50%
				Size:     big.NewInt(100 * 1024),  // 100T
				Duration: 100,                     // 100 days
			},
			{
				Ratio:    big.NewInt(8_0000_0000), // 增发比例 80%
				Size:     big.NewInt(1024 * 1024), // 1PB
				Duration: 150,                     // 150 days
			},
			{
				Ratio:    big.NewInt(10_0000_0000),     // 增发比例 100%
				Size:     big.NewInt(50 * 1024 * 1024), // 50 PB
				Duration: 200,                          // 200 days
			},
			{
				Ratio:    big.NewInt(6_0000_0000),        // 增发比例 80%
				Size:     big.NewInt(1024 * 1024 * 1024), // 1EB
				Duration: 300,                            // 300 days
			},
			{
				Ratio:    big.NewInt(20_0000_0000),            // 增发比例 50%
				Size:     big.NewInt(50 * 1024 * 1024 * 1024), // 50EB
				Duration: 730,                                 // 730 days
			},
		},

		Decimals:      decimals,
		MinimumRation: new(big.Int).Mul(new(big.Int).Div(OneBillion, OneHudred), big.NewInt(3)), // 最小增发比例3%
		InitialSupply: new(big.Int).Mul(big.NewInt(1000_0000_0000), decimals),                   // 设置初始发行量，1000亿，精度为8，设Memo发行价为 0.0002U
		InitialTarget: new(big.Int).Mul(big.NewInt(500_0000_0000), decimals),                    // 初始增发奖励目标500亿，到达后减半至250亿，以此类推，直到达到最小增发率

		InitialKeeperPledge:   new(big.Int).Mul(big.NewInt(5000_0000), decimals), // 五千万 Memo
		InitialProviderPledge: new(big.Int).Mul(big.NewInt(100_0000), decimals),  // 一百万 Memo

		InitialSize:  big.NewInt(500),                               // 500GB
		InitialPrice: new(big.Int).Div(decimals, big.NewInt(1)),     // 1 GB*Day/Memo
		MinimumPrice: new(big.Int).Div(decimals, big.NewInt(10000)), // 0.00001 Memo

		SizeSimulate:     DefaultSizeSimulate,
		PriceSimulate:    DefaultPriceSimulate,
		DurationSimulate: DefaultDurationSimulate,
		ProviderSimulate: DefaultProviderSimulate,
		EnableMaxSize:    false,

		TotalDuration: 4 * 365, // 模拟的总周期，单位天
	}
}
```

即初始发行1000亿Memo，首次增发目标500亿，增发到达500亿后，减半目标为250亿，以此类推，直到达到最小增发比例，参数暂设置为3%，同时通过仿真函数去模拟每天的平均订单量以及平均订单价格等参数。

## 3.2 参数候选

### 3.2.1 选项一

Memo代币精度：8

初始发行总量：1000亿

第一目标增发量：500亿

最小增发率：3%

增发阶段：

+ 阶段一

  目标空间：100TB

  目标时间：100天

  增发比例Ratio：50%

+ 阶段二

  目标空间：1PB

  目标时间：150天

  增发率Ratio：80%

+ 阶段三

  目标空间：50PB

  目标时间：200天

  增发率Ratio：100%

+ 阶段四

  目标空间：1EB

  目标时间：300天

  增发率Ratio：80%

+ 阶段五

  增发率Ratio：50%

Keeper初始质押：5000万Memo

Provider初始质押：100万Memo

注意，增发率还需要结合减半因子$\lambda$，所以真实增发率为$Ratio/2^\lambda$。

计算是否达到目标增发阶段的方法为：累积时空值除目标空间大于目标时间。

### 3.2.2 选项二

Memo代币精度：12

初始发行总量：6亿

第一目标增发量：3亿

最小增发率：3%

增发阶段：

+ 阶段一

  目标空间：100TB

  目标时间：100天

  增发比例Ratio：50%

+ 阶段二

  目标空间：1PB

  目标时间：150天

  增发率Ratio：80%

+ 阶段三

  目标空间：50PB

  目标时间：200天

  增发率Ratio：100%

+ 阶段四

  目标空间：1EB

  目标时间：300天

  增发率Ratio：80%

+ 阶段五

  增发率Ratio：50%

Keeper初始质押：30万Memo

Provider初始质押：6000Memo

## 3.3 仿真结果

### 3.3.1 选项一

以选项一运行一次仿真，首先将期限定为4年，运行一次，得到的最终结果如下：

```json
Day: 1459 reward: 10254646 MintLevel: 3 HalfFactor 3
spacetime: 2987941389080 esize: 50.00 EiB
Calculate Dur: 55
TotalPay: 315366497744 TotalPaid: 305414361560 TotalReward 93597314744
TotalLiquid: 109845178559 TotalSupply: 193597314744 Percent: 0.5673899904254965
TotalSize: 3.63 EiB
OrderSize: 10.19 PiB
TargetReward: 93750000000 PeriodRewad: 6250000000
Issurance ratio: 0.075 Providers Count 73450
Total issue times: 1.93597314744
TotalPledge: 73800000000
```

得出仿真结果：四年后总空间为：3.63EiB，总增发935.9亿Memo，流动Memo占比56.7%。


### 3.3.2 选项二

首先将期限定为4年，运行一次，得到的最终结果如下：

```json
Day: 1459 reward: 32975 MintLevel: 3 HalfFactor 4
spacetime: 3111877297490 esize: 50.00 EiB
Calculate Dur: 57
TotalPay: 1937660346 TotalPaid: 1894302227 TotalReward 575316762
TotalLiquid: 689158644 TotalSupply: 1175316762 Percent: 0.5863599212413854
TotalSize: 3.58 EiB
OrderSize: 10.04 PiB
TargetReward: 581250000 PeriodRewad: 18750000
Issurance ratio: 0.0375 Providers Count 73450
Total issue times: 1.95886127
TotalPledge: 442800000
```

相比创世发行，增发了95.88%的代币，5.75亿的代币，空间为3.58EiB。
