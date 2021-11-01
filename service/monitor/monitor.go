package monitor

import (
	"strconv"
	"wikipay-admin/common"
	"wikipay-admin/models/monitor"
	"wikipay-admin/rpc/client"
	"wikipay-admin/utils"
)

//
type MonitorReq struct {
	Coin        string
	FromAddress string
	ToAddress   string
	Amount      float64
	FeeAddress  string
}

//归集
func MonitorSend(req *MonitorReq) error {
	var (
		coin      = utils.GetCoin(req.Coin)
		rpcClient = client.NewClient()
		err       error
		txId      string
	)
	switch req.Coin {
	case common.USDT_OMNI:
		txId, err = rpcClient.FundedSend(&client.FundSendReq{
			Coin:        coin,
			FromAddress: req.FromAddress,
			ToAddress:   req.ToAddress,
			FeeAddress:  req.FeeAddress,
			Amount:      strconv.FormatFloat(req.Amount, 'f', -1, 64),
		})
	case common.USDT_ERC20, common.ETH:
		txId, err = rpcClient.Withdraw(&client.WithDrawReq{
			Coin:        coin,
			FromAddress: req.FromAddress,
			ToAddress:   req.ToAddress,
			Amount:      req.Amount,
		})
	}
	if err == nil && txId != "" {
		history := monitor.MonitorHistory{
			Coin:             req.Coin,
			FromAddress:      req.FromAddress,
			ToAddress:        req.ToAddress,
			TxId:             txId,
			CollectionAmount: req.Amount,
			Balance: func() float64 {
				f, _ := rpcClient.GetBalance(
					&client.BalanceReq{
						Coin:    coin,
						Address: req.FromAddress,
					})
				return f
			}(),
			Status: 1,
		}
		err = history.Create()
	}
	return err
}
