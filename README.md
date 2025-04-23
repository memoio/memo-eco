# Memo Economic Model Simplified Edition

[TOC]

# 1. Introduction to Memo Protocol

The Memo Protocol is a decentralized storage infrastructure based on blockchain technology. Memo combines a layered architecture, efficient storage proofs, and smart contract settlement to build a reliable, low-cost, and scalable storage service. Memo is committed to serving projects and developers who need reliable data storage, allowing them to simply replace the storage foundation with blockchain-based storage to further decentralize the product.

From an architectural point of view, Memo can be divided into three levels: settlement layer, verification layer, and storage layer.

The settlement layer will aggregate all order information, and the order amount will be gradually distributed to the storage nodes, and the token reward mechanism is also implemented on the settlement layer, all of which is driven by code.

The Keeper node in the verification layer will challenge the Provider node and verify the proof result to decide whether to give the Provider node a withdrawal certificate. All processes in the verification layer will go through the Byzantine fault-tolerant consensus of the verification layer nodes.

The decentralized Provider nodes in the storage layer store real data and regularly submit storage proofs to the verification layer, while using withdrawal certificates to obtain storage income for a stage in the settlement layer.

# 2. Basic Design

## 2.1 Token Definition

Memo is the token used by the Memo protocol, which is used to drive the operation and scale growth of the entire Memo storage protocol.

In the Memo protocol, Keepers and Providers need to first stake Memo to qualify for earning income. At the same time, users need to pay Memo tokens to store and retrieve data. In order to incentivize data storage, Memo will reward Keepers, Providers and long-term holders according to the growth of orders.

Memo is not only the usage token of the protocol, but also the equity token of the protocol. Holders of Memo will benefit from the growth of the Memo ecosystem and participate in the governance of Memo.

## 2.2 Design Principles

The Memo token is designed to create value for users and ecosystem participants, improve productivity, and continuously incentivize ecosystem growth. The economic model of Memo will be designed according to the following principles:

+ The incentives obtained by participants should be proportional to the costs they pay and try to keep it fair.

+ The logic of the economic model should be controlled by code as much as possible to reduce the part of human regulation.
+ The economic model should encourage participants to implement behaviors that are beneficial to the protocol and meet incentive compatibility.

## 2.3 Token Distribution

### 2.3.1 Memo Economic Cycle

Tokens will circulate between users, developers, Keepers, and storage nodes. First, users need to purchase Memo to use Memo's storage services. Keepers and Providers need to purchase Memo pledges before they can provide services to obtain income, and Memo storage fees will flow to Keepers and Providers. At the same time, Keepers, Providers, and long-term coin holders can also pledge Memo in the pledge pool to obtain incentives for the growth of the Memo ecosystem.

Therefore, the generation of Memo is a cycle of consumption and reproduction, and the storage scale will be gradually expanded in the process. The expansion of scale will increase the value of Memo.

### 2.3.2 Memo Usage Scenarios

+ Keepers and Providers first need to deposit a certain amount of Memo in the pledge pool before they can participate in the protocol to provide services.
+ Users or developers need to purchase Memo and transfer it to the order pool before they can sign an order to use Memo's storage service. At the same time, a certain amount of Memo is required to retrieve data.
+ Since Memo also represents the rights and interests of the entire protocol, Memo holders can pledge Memo to the pledge pool to obtain reward income. At the same time, a governance mechanism may be designed in the future, and the use of Memo tokens can participate in the governance process of the protocol.

### 2.3.3 Token Genesis Distribution

Token distribution is divided into Genesis distribution and Incentive Reward Distribution. The Genesis Token Distribution Share is composed of the following:

+ Investors: 29%, which will be unlocked linearly every day within 400 days.
+ Core Team: 27%, controlled by the core team, half will be released immediately, and the other half will be unlocked linearly within 400 days.
+ Partners and Consultants: 14%, part of which will be locked in the order pool, and the rest will be unlocked linearly.
+ Ecological Fund: 29%, temporarily locked, to be extracted and used through governance in the future.

## 2.4 Overview of the Reward Model

After the initial token issuance, a mechanism needs to be designed to stimulate ecological development, incentivize storage nodes to provide storage space as soon as possible, and incentivize users to store data.

First, maintainers (Keepers) and storage nodes (Providers) need to stake Memo tokens to start providing services, and any coin holder can stake Memo in the staking pool to obtain reward income.

In addition to the staking pool, the order pool records the order amount of each storage node and maintenance node. Every time a user signs an order with a storage node, the user needs to transfer the token amount to the account of the Provider and Keeper nodes in the order pool in advance, but the node cannot withdraw funds immediately. It can only upload storage proofs on time during the storage period and slowly take away the funds in the order pool.

The amount sent to the storage node order pool is gradually collected by the storage node throughout the storage period; the amount sent to the Keeper is gradually sent to the Keeper as the storage node collects income, and the other part is sent directly to the Keeper when the order expires.

In addition, when signing an order, a reward amount will be calculated based on the order information and parameters, and then the rewarded tokens will be sent to the pledge pool, and all accounts that pledge tokens in the pledge pool will share the rewards equally according to their shares.

As can be seen from the above, Memo's token rewards will be tied to the usage of the Memo protocol. When the scale of the Memo protocol expands, more Memo will be rewarded to meet the demand. Memo's rewards will be halved after the target reward value is triggered, so early ecological participants can get more dividends.

### 2.4.1 Staking Mechanism

Keeper and Provider need to pledge before providing services, staking the minimum main tokens `KeeperPledge` and `ProviderPledge`. In addition, any user holding the main token can pledge their own tokens to the pledge pool to obtain token rewards.

### 2.4.2 Profit sharing mechanism of pledge pool

The pledge pool is a fund pool with the main token Memo as the pledge token. Assume that at time T, a person $k_1$ pledges $x_1$ of the main token Memo in the pledge pool and obtains $x_1$ mMemo. At this time, the number of Memos in the pool is $\Zeta_{Memo}= x_1$, and the number of mMemos is $\Zeta_{mMemo}=x_1$.

In the subsequent system operation, assuming that the system receives $\Delta_x$ Memo rewards, then $\Zeta_{Memo1}= x_1+\Delta_x, \Zeta_{mMemo1}=x_1$, then $x_1$ mMemo can be exchanged for $(x_1+\Delta_x)$ Memos, that is, each mMemo is worth $\frac{x_1+\Delta_x}{x_1}$ Memos.

If at time B, another person $k_2$ pledges $x_2$ Memos, then it will receive $\frac{\Zeta_{mMemo1}}{\Zeta_{Memo1}}x_2=\frac{x_1x_2}{x_1+\Delta_x}$ mMemos, and $\Zeta_{Memo2}=x_1+x_2+\Delta_x,\Zeta_{mMemo2}=x_1+\frac{x_2x_1}{x_1+\Delta_x}$

If there is no reward at this time, the number of tokens that $k_1$ can redeem is $x_1 \cdot \frac{x_1+x_2+\Delta_x}{x_1+\frac{x_2x_1}{x_1+\Delta_x}}=x_1+\Delta_x$, the amount that $k_2$ can redeem is $\frac{x_1x_2}{x_1+\Delta_x}\cdot\frac{x_1+x_2+\Delta_x}{x_1+\frac{x_2x_1}{x_1+\Delta_x}}=x_2$, which is equal to its contribution. Once new rewards are added to the pool and $\Zeta_{Memo1}$ is changed, $k_2$ will be redeemed for more Memo than $x_2$.

The subsequent new pledges and reward distribution are similar to the above.

### 2.4.3 Staking Pool Reward Mechanism

Reward distribution, including long-term ecological incentives and reward rewards, to encourage maintainers, storage nodes, users and token holders to contribute to the ecosystem. This part of the distribution will be gradually released as the order volume grows.

The upper limit of reward distribution is half of the Genesis issuance.

The reward is linked to the storage volume, and some reward targets are set. Every time an order is triggered, a reward amount will be set according to the current storage status and order data, and issued to the pledge pool, and all pledgers will be allocated according to the pledge amount.

# 3. Ecological composition

## 3.1 Foundation

The Memo Foundation is mainly used for ecosystem construction, marketing and community maintenance, and at the same time, part of the funds are used for investment to promote ecological development and maintain the long-term sustainable operation of the foundation.

The Foundation Memo Protocol has the following responsibilities:

+ Organize development teams or outsource tasks to complete the implementation and iterative upgrade of the Memo protocol.

+ Support and fund ecological applications developed based on Memo.
+ Long-term participation in Memo's community operations and continuous contributions.

The Memo Foundation has the right to initiate proposals on system governance, and then the community decides whether the proposal will be finally implemented. The Foundation can initiate and include but are not limited to the following suggestions:

+ Modify the economic parameters of the system
+ Propose to upgrade the technical solution.
+ Punish Keeper nodes that do evil or do nothing.
+ Punish Provider nodes that do evil or do nothing.

## 3.2 Community Ecology

As a decentralized storage solution combined with blockchain, the development and development of the Memo protocol cannot be separated from the support of the community. The Memo Foundation will actively organize and establish communities with different functions, including ecological governance, developer communities, and coin holder communities, so as to promote the healthy and stable development of the protocol and ecology in many aspects.

## 3.3 Keeper

The node must first pledge a sufficient margin before it can become a Keeper node within the protocol.
As a Keeper node, the following obligations must be fulfilled:

+ The shares must not be less than the specified amount of margin.

+ Ensure long-term online activities and maintain historical data of the verification layer.

+ Check the storage proof of the Provider and issue a withdrawal order to the Provider through the Byzantine Fault Tolerance consensus.

+ When the Provider loses data, schedule the data repair process.

At the same time, the Keeper node has the following rights and interests:

+ Extract a certain proportion of the orders it manages as income.

+ The pledged margin will obtain the reward of the protocol.

If the Keeper node does evil or does nothing, it may be punished, that is, lose the margin, and can no longer continue to provide services to obtain income.

## 3.4 Provider

Nodes with free storage space can become Providers. First, pledge a sufficient amount of margin, and then sign orders with users to obtain income.

As a Provider, you need to fulfill the following obligations:

+ Pledge a sufficient amount of margin in the pledge pool.

+ Store data according to the standards specified in the order, and ensure the reliability and availability of the data.

+ Submit storage proof to the verification layer on time.

At the same time, the Provider has the following rights and interests

+ Gradually collect user order payments.

+ Obtain reward income from the staking pool.

If the Provider fails to fulfill its obligations, such as losing data or failing to submit storage proofs on time, it will be punished, with deposits deducted and orders transferred.

## 3.5 Developers

Developers can use the Memo protocol in their own projects. More projects adopt Memo's storage services, which will bring more users and storage needs to Memo.

The Foundation can screen developers who are valuable to the Memo protocol through some developer activities, give token incentives, and provide incubation services.
Developers need to fulfill the following obligations:

# 3. Economic model simulation

## 3.1 Parameter selection

First, select an economic model configuration as follows:

```go
type TokenConfig struct {
	TotalSupply int64 // Genesis token quantity, default 600_000_000 memo
	InitSupply int64 // Initial issuance, default 100_000_000
	LockSupply int64 // Locked quantity, default 100_000_000
	LinearSupply int64 // Linear release quantity, default 100_000_000
	LockDay int64 // Lock release time, default 540 days, that is, it can only be used after 540 days
	LinearDay int64 // Linear release period, default 180 days, 1/180 released every day
}

type MintConfig struct {
	RewardTarget int64 // Reward cap, Default 600_000_000
	RatioInit int64 // Reward ratio, default 2
	RatioDecimal int64
	RatioAlter int64 // Reward ratio adjustment parameter, default 150;
}

type PledgeConfig struct {
	InRatio int64 // If the annualized return on x days exceeds this value, the pledge will increase, default 100%
	OutRatio int64 // If the annualized return on x days is less than this value, the pledge will be reduced, default 25%
}

type RoleConfig struct {
	KeeperPledge int64 // Minimum pledge of keeper, default 1_000
	ProviderPledge int64 // Minimum pledge of provider, default 1_00
	KeeperCntPerGroup uint64 // Upper limit of keeper per group, default 10
	ProviderCntPerGroup uint64 // Upper limit of provider per group, default 5000
	ProviderStorage uint64 // Upper limit of provider storage, default 8TB, too small will cause simulation time to increase
	ProviderCreatePerDay uint64 // Upper limit of provider added per group per day, default 150, 150-200 will be added per day in the model
}

type OrderConfig struct {
	DefaultSize int64 // Order size, default 8GB
	DefaultDuration uint64 // Order length, default 365 days, greater than 100 days
	DefaultPrice uint64 // Price per GB of order, default is 1073741824

	LinearRate int64 // Linear release rate, default is 4, i.e. 4%, paid to keeper
	EndRate int64 // Expiration release rate, default is 1, i.e. 1%, paid to keeper when order expires
	TaxRate int64 // Management fee, default is 1, i.e. 1%
}

type SimuConfig struct {
	Duration uint64 // Simulation time length, default is 1000 days
	Detail bool // Whether to print intermediate information
}
```

## compile and run

```
## compile
> go build
## config is stored in ~/.simu/config.toml; see result by opening localhost:10888 in web browser
> ./memo-eco
## modify config params, run again
> ./memo-eco --config=~/.simu/config.toml
```
