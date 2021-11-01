package monitor

import (
	"wikipay-admin/models/wallet"
	"wikipay-admin/rpc/client"
)

//同步链上数据
func SyncAccountAndBalance(req *client.BaseReq) ([]client.Balances, error) {
	client := client.NewClient()
	srcCoin := req.Coin

	balances, err := client.SyncAccountAndBalance(req)
	if err != nil {
		return nil, err
	}

	for _, v := range balances {
		account, _ := new(wallet.Account).GetAccountWithAddress(v.Address)
		userAddress := MonitorUserAddress{
			Coin:    srcCoin,
			Address: v.Address,
			Amount:  v.Amount,
			UserId:  account.UserId,
		}
		CreateAndUpdate(&userAddress)
	}
	return balances, nil
}
