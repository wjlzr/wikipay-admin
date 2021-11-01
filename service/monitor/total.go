package monitor

import (
	"strings"
	"wikipay-admin/common"
	"wikipay-admin/models/monitor"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/utils"
)

//
func WithDrawWithNow() ([]monitor.WithDrawInfo, error) {
	wds, err := monitor.FindWithDrawWithNow()
	if err != nil {
		return nil, err
	}

	//	rpcClient := client.NewClient()
	for i, wd := range wds {
		coin := utils.GetCoinName(wd.Coin, false)
		// usdPrice := rpcClient.GetUsdPrice(&client.BaseReq{
		// 	Coin: coin,
		// })
		usdPrice := common.GetCoinUsdPrice(coin)
		wds[i].Usd = xfloat64.Mul(wd.Amount, usdPrice)
	}
	return wds, nil
}

//资产统计分类(公司资产、用户资产)
func AssetsGroup() ([]monitor.WithDrawInfo, error) {
	assets, err := monitor.FindAssetGroup()
	if err != nil {
		return nil, err
	}

	infos := make([]monitor.WithDrawInfo, 0)
	userAsset := monitor.WithDrawInfo{}
	companyAsset := monitor.WithDrawInfo{}

	//	rpcClient := client.NewClient()
	for i, asset := range assets {
		coin := utils.GetCoinName(asset.Coin, false)
		// usdPrice := rpcClient.GetUsdPrice(&client.BaseReq{
		// 	Coin: coin,
		// })
		usdPrice := common.GetCoinUsdPrice(coin)
		if strings.Index(asset.Coin, common.USDT) >= 0 {
			switch asset.Type {
			case 1: //公司
				companyAsset.Coin = common.USDT
				companyAsset.Type = 1
				companyAsset.Usd = xfloat64.Add(companyAsset.Usd, xfloat64.Mul(asset.Amount, usdPrice))
				companyAsset.Amount = xfloat64.Add(companyAsset.Amount, asset.Amount)
			case 2: //个人
				userAsset.Coin = common.USDT
				userAsset.Type = 2
				userAsset.Usd = xfloat64.Add(userAsset.Usd, xfloat64.Mul(asset.Amount, usdPrice))
				userAsset.Amount = xfloat64.Add(userAsset.Amount, asset.Amount)
			}
		} else {
			assets[i].Usd = xfloat64.Mul(asset.Amount, usdPrice)
			infos = append(infos, assets[i])
		}
	}
	infos = append(infos, userAsset)
	infos = append(infos, companyAsset)

	return infos, nil
}

//归集资 产
func CollectAssetsGroup() ([]monitor.WithDrawInfo, error) {
	assets, err := monitor.FindAssetGroup()
	if err != nil {
		return nil, err
	}

	//	rpcClient := client.NewClient()
	for i, asset := range assets {
		coin := utils.GetCoinName(asset.Coin, false)
		// usdPrice := rpcClient.GetUsdPrice(&client.BaseReq{
		// 	Coin: coin,
		// })
		usdPrice := common.GetCoinUsdPrice(coin)
		assets[i].Usd = xfloat64.Mul(asset.Amount, usdPrice)
	}
	return assets, nil
}
