package monitor

import (
	"wikipay-admin/common"
	"wikipay-admin/models/monitor"
	"wikipay-admin/models/wallet"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/rpc/client"
	"wikipay-admin/utils"
)

//统计
func CalcSatistical() ([]monitor.SatisticaInfo, error) {
	var data monitor.MonitorAddress

	infos, err := data.Satistical()
	if err != nil {
		return nil, err
	}

	mInfo := map[string]map[int]float64{
		common.BTC: map[int]float64{
			1: 0,
			2: 0,
		},
		common.ETH: map[int]float64{
			1: 0,
			2: 0,
		},
		common.USDT_ERC20: map[int]float64{
			1: 0,
			2: 0,
		},
		common.USDT_OMNI: map[int]float64{
			1: 0,
			2: 0,
		},
	}

	for _, v := range infos {
		mInfo[v.Coin][v.Type] = v.Amount
	}
	satisticaInfo := make([]monitor.SatisticaInfo, 0)
	//rpcClient := client.NewClient()

	for i := 0; i < len(utils.UpperCoins); i++ {
		info := monitor.SatisticaInfo{}
		info.Coin = utils.UpperCoins[i]
		info.HoldWalletAmount = mInfo[utils.UpperCoins[i]][1]
		info.ColdWalletAmount = mInfo[utils.UpperCoins[i]][2]
		info.Total = xfloat64.Add(info.HoldWalletAmount, info.ColdWalletAmount)
		if info.Total > 0 {
			coin := utils.GetCoinName(info.Coin, false)
			// usdPrice := rpcClient.GetUsdPrice(&client.BaseReq{
			// 	Coin: coin,
			// })
			usdPrice := common.GetCoinUsdPrice(coin)
			// cnyPrice := rpcClient.GetCurrency(&client.BaseReq{
			// 	Coin: "CNY",
			// })
			cnyPrice := common.GetCurrencyPrice(common.CURRENCY_CNY)
			info.TotalCny = xfloat64.Float64Truncate(xfloat64.Mul(info.Total, xfloat64.Mul(usdPrice, cnyPrice)), 2)
		} else {
			info.TotalCny = 0

		}
		satisticaInfo = append(satisticaInfo, info)
	}
	return satisticaInfo, nil
}

//资产对比统计
func AssetComparison() (assetComparisonInfo monitor.AssetComparisonInfo, err error) {
	var (
		data         monitor.MonitorAddress
		hotWallet    monitor.HotWallet
		addressAsset monitor.AddressAssets
	)

	infos, err := data.Satistical()
	if err != nil {
		return assetComparisonInfo, err
	}

	mInfo := map[string]map[int]float64{
		common.BTC: map[int]float64{
			1: 0,
			2: 0,
		},
		common.ETH: map[int]float64{
			1: 0,
			2: 0,
		},
		common.USDT_ERC20: map[int]float64{
			1: 0,
			2: 0,
		},
		common.USDT_OMNI: map[int]float64{
			1: 0,
			2: 0,
		},
	}

	for _, v := range infos {
		mInfo[v.Coin][v.Type] = v.Amount
	}

	cnyPrice := common.GetCurrencyPrice(common.CURRENCY_CNY)
	for _, v := range utils.UpperCoins {

		usdPrice := common.GetCoinUsdPrice(v)
		accountAssets, err := wallet.Account{}.AssetStatistics(v)
		if err != nil {
			return assetComparisonInfo, err
		}
		//热钱包
		hotWallet.Coin = v
		hotWallet.HoldWalletAmount = mInfo[v][1]
		hotWallet.AccountAssets = accountAssets
		hotWallet.Diff = xfloat64.Float64Truncate(xfloat64.Sub(hotWallet.HoldWalletAmount, hotWallet.AccountAssets), 2)
		hotWallet.DiffCny = xfloat64.Float64Truncate(xfloat64.Mul(hotWallet.Diff, xfloat64.Mul(usdPrice, cnyPrice)), 2)

		//地址资产合计
		addressAsset.Coin = v
		addressAsset.AccountAssets = accountAssets
		addressAsset.AddressAssetsAmount = xfloat64.Add(mInfo[v][1], mInfo[v][2])
		addressAsset.Diff = xfloat64.Float64Truncate(xfloat64.Sub(addressAsset.AddressAssetsAmount, addressAsset.AccountAssets), 2)
		hotWallet.DiffCny = xfloat64.Float64Truncate(xfloat64.Mul(addressAsset.Diff, xfloat64.Mul(usdPrice, cnyPrice)), 2)

		assetComparisonInfo.HotWallet = append(assetComparisonInfo.HotWallet, hotWallet)
		assetComparisonInfo.AddressAssets = append(assetComparisonInfo.AddressAssets, addressAsset)
	}

	return
}

//同步公司地址
func SyncAddress(req *monitor.MonitorAddressReq) error {
	data := monitor.MonitorAddress{
		Coin: req.Coin,
		Type: req.Type,
	}

	addresses, err := data.Get()
	if err != nil {
		return err
	}

	rpcClient := client.NewClient()
	coin := utils.GetCoin(req.Coin)

	for _, addr := range addresses {
		f, _ := rpcClient.GetBalance(&client.BalanceReq{
			Coin:    coin,
			Address: addr.Address,
		})
		address := monitor.MonitorAddress{
			Coin:    addr.Coin,
			Type:    addr.Type,
			Address: addr.Address,
			Amount:  f,
		}
		err := address.UpdateAddress()
		if err != nil {
			return err
		}
	}
	return nil
}
