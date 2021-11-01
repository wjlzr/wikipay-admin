package client

import (
	"sort"
	"strings"

	"wikipay-admin/utils"
)

//rpc client
type RpcCient struct {
	client *GrpClient
}

//
func NewClient() *RpcCient {
	return &RpcCient{
		client: NewGrpClient(),
	}
}

//获取指定地址余额
func (r *RpcCient) GetBalance(req *BalanceReq) (float64, error) {
	req.Coin = utils.GetCoin(req.Coin)

	balance := r.client.GetBalance(req)
	return balance, nil
}

//获取网络状态
func (r *RpcCient) GetNetworkInfo(req *BaseReq) (*NetWorkInfoResp, error) {
	req.Coin = utils.GetCoin(req.Coin)

	resp, err := r.client.GetNetworkInfo(req)
	return resp, err
}

//获取币的美元价格
func (r *RpcCient) GetUsdPrice(req *BaseReq) float64 {
	req.Coin = utils.GetCoinName(req.Coin, false)

	price := r.client.GetUsdPrice(req)
	return price
}

//获取币的美元价格
func (r *RpcCient) GetCurrency(req *BaseReq) float64 {
	req.Coin = strings.ToUpper(req.Coin)

	price := r.client.GetCurrency(req)
	return price
}

//获取地址和余额
func (r *RpcCient) SyncAccountAndBalance(req *BaseReq) ([]Balances, error) {
	req.Coin = utils.GetCoin(req.Coin)

	balances, err := r.client.GetAccountAndBalance(req)
	if err != nil {
		return nil, err
	}
	//从大到小排序
	sort.Slice(balances, func(i, j int) bool {
		if balances[i].Amount > balances[j].Amount {
			return true
		}
		return false
	})
	return balances, err
}

//验证地址
func (r *RpcCient) ValidateAddress(req *BalanceReq) bool {
	isVaild, _ := r.client.ValidateAddress(req)
	return isVaild
}

//拽定手续费归集
func (r *RpcCient) FundedSend(req *FundSendReq) (string, error) {
	return r.client.FundedSend(req)
}

//发送交易
func (r *RpcCient) Withdraw(req *WithDrawReq) (string, error) {
	return r.client.Withdraw(req)
}
