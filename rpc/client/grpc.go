package client

import (
	"context"
	"flag"

	"wikipay-admin/rpc"
	"wikipay-admin/utils"

	"google.golang.org/grpc"
)

var (
	serverAddr  = flag.String("server_addr", "18.163.191.5:6688", "The server address in the format of host:port")
	rpcUserName = flag.String("rpc_username", "wikipay", "rpc username")
	rpcPassword = flag.String("rpc_pwd", "Wiki08@Pay", "rpc password")
)

//
type GrpClient struct {
	conn *grpc.ClientConn
}

type BaseReq struct {
	Coin string `json:"coin" form:"coin" binding:"required"`
}

//提现
type WithDrawReq struct {
	Coin        string  `json:"coin" form:"coin" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	FromAddress string  `json:"fromAddress" binding:"required"`
	ToAddress   string  `json:"toAddress" binding:"required"`
}

//
type BalanceReq struct {
	Coin    string `json:"coin" form:"coin" binding:"required"`
	Address string `json:"address" form:"address"`
}

//
type NetWorkInfoResp struct {
	Version string `json:"version"`
	Active  bool   `json:"active"`
	Errors  string `json:"errors"`
}

//
type Balances struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

//
type FundSendReq struct {
	Coin        string
	FromAddress string
	ToAddress   string
	Amount      string
	FeeAddress  string
}

//
func NewGrpClient() *GrpClient {
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithPerRPCCredentials(&AuthParam{
		Username: *rpcUserName,
		Password: *rpcPassword,
	}))

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err == nil {
		return &GrpClient{
			conn: conn,
		}
	}
	return nil
}

//验证地址
func (g *GrpClient) ValidateAddress(req *BalanceReq) (bool, error) {
	client := rpc.NewWalletClient(g.conn)
	pbReq := rpc.ValidateAddressRequest{
		Coin:    utils.GetCoin(req.Coin),
		Address: req.Address,
	}

	resp, err := client.ValidateAddress(context.Background(), &pbReq)
	if err != nil || resp == nil {
		return false, err
	}

	return resp.Status, nil
}

//验证地址
func (g *GrpClient) Withdraw(req *WithDrawReq) (string, error) {
	client := rpc.NewWalletClient(g.conn)
	pbReq := rpc.WithdrawRequest{
		Coin:        utils.GetCoin(req.Coin),
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Amount:      req.Amount,
	}

	resp, err := client.Withdraw(context.Background(), &pbReq)
	if err != nil || resp == nil {
		return "", err
	}

	return resp.Txid, nil
}

//获取余额
func (g *GrpClient) GetBalance(req *BalanceReq) float64 {
	client := rpc.NewWalletClient(g.conn)

	pbReq := rpc.BalanceRequest{
		Coin:    req.Coin,
		Address: req.Address,
	}

	resp, err := client.GetBalance(context.Background(), &pbReq)
	if err != nil || resp == nil {
		return 0
	}
	return resp.Balance
}

//获取网络状态
func (g *GrpClient) GetNetworkInfo(req *BaseReq) (*NetWorkInfoResp, error) {
	client := rpc.NewWalletClient(g.conn)

	pbReq := rpc.NetworkInfoRequest{
		Coin: req.Coin,
	}

	resp, err := client.GetNetworkInfo(context.Background(), &pbReq)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return &NetWorkInfoResp{
			Version: resp.Version,
			Active:  resp.Active,
			Errors:  resp.Warnings,
		}, nil
	}
	return nil, nil
}

//获取币的美元价格
func (g *GrpClient) GetUsdPrice(req *BaseReq) float64 {
	client := rpc.NewWalletClient(g.conn)

	pbReq := rpc.UsdPriceRequest{
		Coin: req.Coin,
	}

	resp, err := client.GetUsdPrice(context.Background(), &pbReq)
	if err != nil {
		return 0
	}
	return resp.Price
}

//获取美元况其它法币的价格
func (g *GrpClient) GetCurrency(req *BaseReq) float64 {
	client := rpc.NewWalletClient(g.conn)

	pbReq := rpc.CurrencyRequest{
		Name: req.Coin,
	}

	resp, err := client.GetCurrency(context.Background(), &pbReq)
	if err != nil {
		return 0
	}
	return resp.Price
}

//获取账户和余额
func (g *GrpClient) GetAccountAndBalance(req *BaseReq) ([]Balances, error) {
	client := rpc.NewWalletClient(g.conn)
	pbReq := rpc.AccountAndBalanceRequest{
		Coin: req.Coin,
	}

	resp, err := client.GetAccountAndBalance(context.Background(), &pbReq)
	var balances []Balances
	for _, v := range resp.Accounts {
		balance := Balances{
			Address: v.Address,
			Amount:  v.Balance,
		}
		balances = append(balances, balance)
	}
	return balances, err
}

//
func (g *GrpClient) FundedSend(req *FundSendReq) (string, error) {
	client := rpc.NewWalletClient(g.conn)
	pbReq := rpc.FundedSendRequest{
		Coin:        req.Coin,
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Amount:      req.Amount,
		FeeAddress:  req.FeeAddress,
	}

	resp, err := client.FundedSend(context.Background(), &pbReq)
	if err != nil {
		return "", err
	}
	return resp.TxId, nil
}
