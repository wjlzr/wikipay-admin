package common

const (
	//coins
	BTC        = "BTC"
	ETH        = "ETH"
	USDT_ERC20 = "USDT-ERC20"
	USDT_OMNI  = "USDT-OMNI"
	USDT       = "USDT"
	USD        = "USD"

	//rpc coin
	RpcBTC   = "btc"
	RpcETH   = "eth"
	RpcOMNI  = "omni"
	RPCERC20 = "erc20"

	//currency
	CURRENCY_CNY = "CNY"

	//
	RedisWithDrawingOrders = "wikipay_withdrawing"
)

//
var (
	CoinBits = map[string]int{
		USDT: 4,
		BTC:  8,
		ETH:  5,
		USD:  2,
	}

	Coins = []string{USDT, BTC, ETH, USD}

	CNYCoinBit = map[string]int{
		USDT: 0,
		BTC:  4,
		ETH:  3,
	}

	CoinAccountType = map[string]int{
		USDT: 1,
		BTC:  1,
		ETH:  1,
		USD:  2,
	}
)
