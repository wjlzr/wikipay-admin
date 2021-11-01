package financial

import (
	"bytes"
	"errors"
	"fmt"
	"wikipay-admin/common"
	orm "wikipay-admin/database"
	"wikipay-admin/models/monitor"
	"wikipay-admin/models/wallet"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/redis"
	"wikipay-admin/rpc/client"
	"wikipay-admin/tools"
	"wikipay-admin/utils"

	"github.com/jinzhu/gorm"
)

//
const (
	DBError          = 1 //数据库错误
	RecordNotExist   = 2 //记录不存在
	RpcError         = 3 //rpc连接错误
	WithdrawError    = 4 //提现错误
	NoTxid           = 5 //无交易哈希
	WithdrawSuccess  = 6 //提现成功
	TypeError        = 7 //类型错误
	NotFindHotWallet = 8 //没有热钱包地址
	NoMoney          = 9 //热钱包没有足够数量
	WithDrawing      = 10
	SystemError      = 11
)

var (
	//提示信息
	TipMessage = map[int]string{
		DBError:          "数据库错误",
		RecordNotExist:   "没有找到相关记录",
		RpcError:         "rpc连接错误",
		WithdrawError:    "提现错误",
		NoTxid:           "没有交易哈希",
		WithdrawSuccess:  "成功",
		TypeError:        "充值、提现类型不存在",
		NotFindHotWallet: "没有找到热钱包地址",
		NoMoney:          "热钱包数量不足",
		WithDrawing:      "提现正在处理中，请勿重复操作",
		SystemError:      "系统异常，请稍后再试",
	}
)

//提现审核
type WithdrawAudit struct {
	CreateAt  int64  `json:"createAt"`
	OrderId   string `json:"orderId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	NickName  string `json:"nickName"`
	Address   string `json:"address"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	Account   string `json:"account"`
	Available string `json:"available"`
	Money     string `json:"money"`
	Usd       string `json:"usd"`
	Fee       string `json:"fee"`
	Total     string `json:"total"`
	Ip        string `json:"ip"`
	Status    int    `json:"status"`
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

//分页
type Pagination struct {
	PageSize int32 `form:"pageSize" json:"pageSize"`
	PageNum  int32 `form:"pageNum" json:"pageNum"`
}

//提现审核请求
type WithdrawAuditReq struct {
	Pagination
	StartTime int64  `json:"startTime" form:"startTime"`
	EndTime   int64  `json:"endTime" form:"endTime"`
	Info      string `json:"info" form:"info"`
	Status    int    `json:"status" form:"status"`
}

//提现请求结构
type WithdrawReq struct {
	OrderId string `json:"orderId" form:"orderId" binding:"required"`
	Type    int    `json:"type" form:"type" binding:"required"` //1、通过  2、拒绝
	Content string `json:"content" form:"content"`
}

//提现审核列表
func FindWithdrawAuditInfo(req *WithdrawAuditReq) ([]WithdrawAudit, int, error) {
	var (
		infos, countInfo []WithdrawAudit
	)
	sql := fmt.Sprintf(`
			SELECT a.create_at,
				a.order_number AS order_id,
				CONCAT(b.first_name,b.last_name) AS name,
				CONCAT(b.area_code,b.phone) AS phone,
				a.to_address AS address,
				b.email,
				b.nick_name,
				a.coin,
				a.amount,
				(IF(a.account_type=2,'USD',IF(LEFT(a.coin,4)='USDT','USDT',a.coin))) AS account,
				c.available,
				IF(a.account_type=2,a.amount * a.usd,'') AS money,
				a.usd,
				b.last_login_ip AS ip,
				a.fee,
				a.status,
				IF(a.account_type=2,a.amount * a.usd + a.fee*a.usd,a.amount + a.fee) AS total	
			FROM withdraw_deposit a
			LEFT JOIN user b ON a.user_id = b.id 
			LEFT JOIN account c ON c.user_id = a.user_id AND c.type = a.account_type AND c.coin = (IF(a.account_type=2,'USD',IF(LEFT(a.coin,4)='USDT','USDT',a.coin)))
			WHERE a.type = 2  %s
			ORDER BY a.create_at DESC
		`,
		func() string {
			var buf bytes.Buffer
			if req.Status == 0 {
				req.Status = 1
			}
			switch req.Status {
			case 2, 3, 4, 5, 6:
				buf.WriteString(fmt.Sprintf(" AND  a.status = %d ", req.Status))
			}
			if req.StartTime > 0 {
				buf.WriteString(fmt.Sprintf(" AND a.create_at >= %d ", req.StartTime))
			}
			if req.EndTime > 0 {
				buf.WriteString(fmt.Sprintf(" AND a.create_at <= %d ", req.EndTime))
			}
			if req.Info != "" {
				buf.WriteString(fmt.Sprintf(" AND (a.order_number = '%s' OR b.phone = '%s')", req.Info, req.Info))
			}
			return buf.String()
		}())
	err := orm.Eloquent.Raw(sql).Find(&countInfo).Error

	err = orm.Eloquent.Raw(fmt.Sprintf(`%s LIMIT %d OFFSET %d `,
		sql,
		req.PageSize,
		req.PageNum,
	)).Find(&infos).Error

	return infos, len(countInfo), err
}

//提币审核
func WithdrawAuditOn(req *WithdrawReq) (int, error) {
	//
	var wd wallet.WithdrawDeposit
	//检测提现订单是否存在
	// sql := fmt.Sprintf(`
	// 		SELECT order_number,coin,user_id,to_address,type,amount,fee,status,account_type,usd FROM withdraw_deposit
	// 		WHERE order_number = %s AND status = 2 AND type = 2`,
	// 	req.OrderId)
	err := orm.Eloquent.Raw(`
			SELECT order_number,coin,user_id,to_address,type,amount,fee,status,account_type,usd FROM withdraw_deposit 
			WHERE order_number = ? AND status = 2 AND type = 2`,
		req.OrderId,
	).Scan(&wd).Error
	if err != nil {
		return DBError, err
	}

	if len(wd.OrderNumber) != 18 {
		return RecordNotExist, errors.New("no record")
	}

	//检测是否有重复操作数据
	if redis.ClusterClient().SIsMember(common.RedisWithDrawingOrders, req.OrderId).Val() {
		return WithDrawing, nil
	} else {
		err := redis.ClusterClient().SAdd(common.RedisWithDrawingOrders, req.OrderId).Err()
		if err != nil {
			return SystemError, err
		}
	}

	switch req.Type {
	case 1:
		var (
			coin        = utils.GetCoin(wd.Coin)
			rpcClient   = client.NewClient()
			fromAddress = ""
		)

		//获取热钱包地址
		if wd.Coin != common.BTC {
			monitorAddress := monitor.MonitorAddress{
				Coin: wd.Coin,
				Type: 1,
			}
			addresses, err := monitorAddress.Get()
			if err != nil || len(addresses) < 1 {
				return NotFindHotWallet, err
			}

			//获取地地数量
			txFee := utils.GetGasPrice(wd.Coin)
			feeCoin := func() string {
				if coin == common.RPCERC20 {
					return common.RpcETH
				}
				return coin
			}()

			for _, v := range addresses {
				amount, _ := rpcClient.GetBalance(&client.BalanceReq{
					Coin:    coin,
					Address: v.Address,
				})

				//手续费
				feeAmount := xfloat64.StrToFloat64(wd.Amount)
				if wd.Coin == common.USDT_ERC20 {
					feeAmount, _ = rpcClient.GetBalance(&client.BalanceReq{
						Coin:    feeCoin,
						Address: v.Address,
					})
				}

				//检查提现金额比热钱包地址大，比手续费小
				if xfloat64.FromStringCmpFloat(wd.Amount, amount) < 0 &&
					xfloat64.FromFloatCmp(feeAmount, txFee) > 0 {
					fromAddress = v.Address
					break
				}
			}
			//数量不足
			if fromAddress == "" {
				return NoMoney, nil
			}
		} else {
			amount, _ := rpcClient.GetBalance(&client.BalanceReq{
				Coin:    coin,
				Address: "",
			})
			//
			if xfloat64.FromStringCmpFloat(wd.Amount, amount) > 0 {
				return NoMoney, nil
			}
		}

		rpcReq := client.WithDrawReq{
			Coin:        coin,
			Amount:      xfloat64.StrToFloat64(wd.Amount),
			FromAddress: fromAddress,
			ToAddress:   wd.ToAddress,
		}

		txId, err := rpcClient.Withdraw(&rpcReq)
		if txId == "" {
			return NoTxid, errors.New("no txid ")
		}

		err = orm.Eloquent.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("UPDATE withdraw_deposit SET tx_hash = ?, status = 3, review_at = ? WHERE order_number = ? AND status = 2 AND type = 2 ",
				txId,
				tools.MilliSecond(),
				req.OrderId).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return DBError, err
		}
		//提现成功后删除订单号
		redis.ClusterClient().SRem(common.RedisWithDrawingOrders, req.OrderId)
		return WithdrawSuccess, nil
	case 2:
		var (
			coin   string
			amount float64
		)
		switch wd.AccountType {
		case 1:
			coin = utils.GetCoinName(wd.Coin, true)
			amount = xfloat64.Float64Truncate(xfloat64.FromStringAdd(wd.Amount, wd.Fee), common.CoinBits[coin])
		case 2:
			coin = common.USD
			amount = xfloat64.Float64Truncate(xfloat64.Add(xfloat64.FromStringMul(wd.Amount, wd.Usd), xfloat64.FromStringMul(wd.Fee, wd.Usd)), common.CoinBits[coin])
		}
		err = orm.Eloquent.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("UPDATE withdraw_deposit SET status = 6,content = ? WHERE order_number = ? AND status = 2", req.Content, wd.OrderNumber).Error; err != nil {
				return err
			}

			if err := tx.Exec("UPDATE account SET available = available + ?, frozen = frozen - ? WHERE user_id = ? AND type = ? AND coin = ?", amount, amount, wd.UserId, wd.AccountType, coin).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return DBError, err
		}
		//提现成功后删除订单号
		redis.ClusterClient().SRem(common.RedisWithDrawingOrders, req.OrderId)
		return WithdrawSuccess, nil
	}
	return TypeError, errors.New("no this type")
}
