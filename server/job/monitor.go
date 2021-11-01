package job

import (
	"strconv"
	"wikipay-admin/models/monitor"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/rpc/client"
	"wikipay-admin/utils"
)

//omni归集
func omniMonitor(info SettingInfo) {
	if xfloat64.FromFloatCmp(info.Min, 0) > 0 && xfloat64.FromFloatCmp(info.Max, 0) > 0 {
		if xfloat64.FromFloatCmp(info.Min, info.Max) > 0 {
			return
		}
	}
	//同步数据
	if syncData(info.Coin) {
		//获取数据
		datas := getDatas(info.Coin, info.Type, info.Min, info.Max)
		if datas != nil {
			//检查是否有归集数据
			ratioLen := len(datas.Ratio)
			if len(datas.Address) == 0 || ratioLen == 0 {
				return
			}

			var (
				usedAmount = 0.0 //已归集数量
				index      = 0   //归集地址下标
				toAddress  = datas.Ratio[index].Address
				coin       = utils.GetCoin(info.Coin)
				rpcClient  = client.NewClient()
			)

			for _, v := range datas.Address {
				if usedAmount >= datas.Ratio[index].Amount {
					index = ratioLen - 1
					usedAmount = 0.0
					toAddress = datas.Ratio[index].Address
				}

				txId, err := rpcClient.FundedSend(&client.FundSendReq{
					Coin:        coin,
					FromAddress: v.Address,
					ToAddress:   toAddress,
					Amount:      strconv.FormatFloat(v.Amount, 'f', -1, 64),
					FeeAddress:  info.FeeAddress,
				})
				if err == nil && txId != "" {
					history := monitor.MonitorHistory{
						Coin:             info.Coin,
						FromAddress:      v.Address,
						ToAddress:        toAddress,
						TxId:             txId,
						CollectionAmount: v.Amount,
						Balance: func() float64 {
							f, _ := rpcClient.GetBalance(
								&client.BalanceReq{
									Coin:    coin,
									Address: v.Address,
								})
							return f
						}(),
						Status: 1,
					}
					err = history.Create()
					usedAmount = xfloat64.Add(usedAmount, v.Amount)
				}
			}
		}
	}
}

//eth归集
func ethMonitor(info SettingInfo) {
	if xfloat64.FromFloatCmp(info.Min, 0) > 0 && xfloat64.FromFloatCmp(info.Max, 0) > 0 {
		if xfloat64.FromFloatCmp(info.Min, info.Max) > 0 {
			return
		}
	}

	//同步数据
	if syncData(info.Coin) {
		//获取数据
		datas := getDatas(info.Coin, info.Type, info.Min, info.Max)
		if datas != nil {
			//检查是否有归集数据
			ratioLen := len(datas.Ratio)
			if len(datas.Address) == 0 || ratioLen == 0 {
				return
			}

			var (
				usedAmount = 0.0 //已归集数量
				index      = 0   //归集地址下标
				toAddress  = datas.Ratio[index].Address
				coin       = utils.GetCoin(info.Coin)
				rpcClient  = client.NewClient()
			)

			gasPrice := utils.GetGasPrice(info.Coin)
			for _, v := range datas.Address {
				if usedAmount >= datas.Ratio[index].Amount {
					index = ratioLen - 1
					usedAmount = 0.0
					toAddress = datas.Ratio[index].Address
				}
				amount := xfloat64.Sub(v.Amount, gasPrice)
				if xfloat64.FromFloatCmp(amount, 0) > 0 {
					txId, err := rpcClient.Withdraw(&client.WithDrawReq{
						Coin:        coin,
						Amount:      amount,
						FromAddress: v.Address,
						ToAddress:   toAddress,
					})

					if err == nil && txId != "" {
						history := monitor.MonitorHistory{
							Coin:             info.Coin,
							FromAddress:      v.Address,
							ToAddress:        toAddress,
							TxId:             txId,
							CollectionAmount: v.Amount,
							Balance: func() float64 {
								f, _ := rpcClient.GetBalance(
									&client.BalanceReq{
										Coin:    coin,
										Address: v.Address,
									})
								return f
							}(),
							Status: 1,
						}
						err = history.Create()
						usedAmount = xfloat64.Add(usedAmount, v.Amount)
					}
				}
			}
		}
	}
}
