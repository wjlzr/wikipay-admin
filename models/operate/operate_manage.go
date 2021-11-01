package operate

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"wikipay-admin/common"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/tools"
	"wikipay-admin/utils"

	"github.com/mitchellh/mapstructure"
)

// OperateManage
type OperateManage struct {
	Id        int64   `gorm:"column:id" json:"id" form:"id"`
	Outline   string  `gorm:"column:outline" json:"outline" form:"outline"`
	Condition Coin    `gorm:"type:json" json:"condition" form:"condition"`
	UserId    string  `gorm:"column:user_id" json:"userId" form:"userId"`
	Status    int64   `gorm:"column:status" json:"status" form:"status"`
	Remark    string  `gorm:"column:remark" json:"remark" form:"remark"`
	StartTime int64   `gorm:"column:start_time" json:"startTime" form:"startTime"`
	EndTime   int64   `gorm:"column:end_time" json:"endTime" form:"endTime"`
	CreateAt  int64   `gorm:"column:create_at" json:"createAt" form:"createAt"`
	UpdateAt  int64   `gorm:"column:update_at" json:"updateAt" form:"updateAt"`
	Cny       float64 `gorm:"-" json:"cny"`
	DataScope string  `gorm:"-" json:"-"`
}

// json类型字段
type Coin struct {
	Coin Sub `json:"coin"`
}

type Sub struct {
	USDT float64 `json:"USDT"`
	BTC  float64 `json:"BTC"`
	ETH  float64 `json:"ETH"`
	CNY  float64 `json:"CNY"`
}

//表名
func (o *OperateManage) TableName() string {
	return "operate_manage"
}

//
func (o Coin) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	if err != nil {
		log.Println("Error Value json.Marshal err", err)
	}
	return string(b), err
}

//
func (o *Coin) Scan(input interface{}) (err error) {
	if _, ok := input.(map[string]interface{}); ok {
		jsonStr, err := json.Marshal(input)
		if err != nil {
			log.Println("Error Scan json.Marshal err", err)
		}
		if err = json.Unmarshal([]byte(jsonStr), o); err != nil {
			log.Println("Error Scan json.Unmarshal1 err", err)
		}
	}
	if err = json.Unmarshal(input.([]byte), o); err != nil {
		log.Println("Error Scan json.Unmarshal2 err", err)
	}
	return
}

//获取列表带分页
func (o OperateManage) GetPage(pageSize int, pageIndex int) (operateManage []OperateManage, count int, err error) {

	table := orm.Eloquent.Select("*").Table(o.TableName())

	//数据权限控制
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(o.DataScope)
	table, err = dataPermission.GetDataScope("operate_manage", table)
	if err != nil {
		return nil, 0, err
	}

	if err := table.Order("create_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&operateManage).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return
}

//创建
func (o *OperateManage) Create() (operateManage OperateManage, err error) {

	if o.Condition.Coin, err = conversion(o.Cny); err != nil {
		return
	}
	tt := tools.MilliSecond()
	o.CreateAt = tt
	o.UpdateAt = tt

	if err = orm.Eloquent.Table(o.TableName()).Create(&o).Error; err != nil {
		return
	}
	operateManage = *o
	return
}

// 更新活动配置
func (o *OperateManage) Update(id int64) (operateManage OperateManage, err error) {

	if err = orm.Eloquent.Table(o.TableName()).Where("id = ?", id).First(&operateManage).Error; err != nil {
		return
	}
	if o.Condition.Coin, err = conversion(o.Cny); err != nil {
		return
	}
	o.CreateAt = operateManage.CreateAt
	o.UpdateAt = tools.MilliSecond()
	orm.Eloquent.Table(o.TableName()).Save(&o)

	operateManage = *o
	return
}

//获取详情
func (o *OperateManage) Get() (operateManage OperateManage, err error) {
	table := orm.Eloquent.Table(o.TableName())
	if err := table.First(&operateManage).Error; err != nil {
		return operateManage, err
	}
	return operateManage, nil
}

//批量删除
func (o *OperateManage) BatchDelete(id []int) (result bool, err error) {
	if err = orm.Eloquent.Table(o.TableName()).Where("id in (?)", id).Delete(&OperateManage{}).Error; err != nil {
		return
	}
	return true, err
}

//人民币换算指定货币
func conversion(cny float64) (sub Sub, err error) {
	//人民币美元汇率
	cnyPrice := common.GetCurrencyPrice(common.CURRENCY_CNY)
	var f = make(map[string]float64, 3)
	for _, v := range utils.OperateCoins {
		//美元和货币汇率
		usdPrice := common.GetCoinUsdPrice(v)
		f[v] = xfloat64.Float64Truncate(xfloat64.Div(cny, xfloat64.Mul(usdPrice, cnyPrice)), common.CNYCoinBit[v])
	}
	if err := mapstructure.Decode(f, &sub); err != nil {
		return Sub{}, err
	}

	sub.CNY = cny
	return sub, nil
}
