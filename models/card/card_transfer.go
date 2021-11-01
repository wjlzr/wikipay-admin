package card

import (
	"bytes"
	"errors"
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"

	"github.com/jinzhu/gorm"
)

type CardTransfer struct {
	Id          int    `json:"id" gorm:"type:int(10) unsigned;primary_key"` //
	OrderNumber string `json:"orderNumber" gorm:"type:bigint(18)"`
	Phone       string `json:"phone"  gorm:"-"`
	UserId      int    `json:"userId" gorm:"type:bigint(10);"`   // 用户序列号
	Type        int    `json:"type" gorm:"type:tinyint(1);"`     // 1、转入 2、转出
	Money       string `json:"money" gorm:"type:decimal(10,2);"` // 金额
	Status      int    `json:"status" gorm:"type:tinyint(1);"`   // 1、成功 2、失败 3、审核中
	CreateAt    int    `json:"createAt" gorm:"type:bigint(13);"` // 创建时间
	DataScope   string `json:"dataScope" gorm:"-"`
	Params      string `json:"params"  gorm:"-"`
}

//
type CardTransferReq struct {
	Id   int `json:"id" binding:"required"`
	Type int `json:"type" binding:"required"`
}

//筛选条件
type SearchParams struct {
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}

//表名
func (CardTransfer) TableName() string {
	return "card_transfer"
}

//创建
func (e *CardTransfer) Create() (CardTransfer, error) {
	var doc CardTransfer
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *CardTransfer) Get() (CardTransfer, error) {
	var doc CardTransfer
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取带分页
func (e *CardTransfer) GetPage(pageSize, pageIndex, status, cardType int, info string, req SearchParams) ([]CardTransfer, int, error) {
	var doc []CardTransfer

	table := orm.Eloquent.Select("card_transfer.*,user.phone").Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var where bytes.Buffer
	where.WriteString(" 1=1 ")
	if cardType == 1 || cardType == 2 {
		where.WriteString(fmt.Sprintf(" AND card_transfer.type = %d ", cardType))
	}
	if status >= 1 && status <= 3 {
		where.WriteString(fmt.Sprintf(" AND card_transfer.status = %d ", status))
	}
	if info != "" {
		where.WriteString(fmt.Sprintf(" AND concat(card_transfer.user_id) like '%s%s%s'", "%", info, "%"))
	}

	if req.StartTime != "" {
		where.WriteString(fmt.Sprintf(" AND card_transfer.create_at >= '%s'", req.StartTime))
	}

	if req.EndTime != "" {
		where.WriteString(fmt.Sprintf(" AND card_transfer.create_at <= '%s'", req.EndTime))
	}

	if err := table.Joins("LEFT JOIN user ON user.id = card_transfer.user_id").
		Where(where.String()).
		Order("card_transfer.create_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	var count int
	table.Where(where.String()).Count(&count)
	return doc, count, nil
}

//更新
func (e *CardTransfer) Update(id int) (update CardTransfer, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除CardTransfer
func (e *CardTransfer) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&CardTransfer{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *CardTransfer) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&CardTransfer{}).Error; err != nil {
		return
	}
	Result = true
	return
}

//审核
func (e *CardTransfer) CardAudit(req *CardTransferReq) error {
	var info CardTransfer

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&info, "id = ?", req.Id).Error; err != nil {
		return err
	}

	//开通数据权限
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return err
	}

	if info.Type != 2 && info.Status != 3 {
		return errors.New("不是转出且状态不是审核中")
	}

	//事务处理
	switch req.Type {
	case 1: //通过
		orderId := models.MustGenerateOrderId(6)
		return orm.Eloquent.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("UPDATE card_transfer SET status = 1 WHERE id = ?", info.Id).Error; err != nil {
				return err
			}

			if err := tx.Exec("UPDATE account SET frozen = frozen - ?, available = available + ?  WHERE user_id = ? AND coin='USD'", info.Money, info.Money, info.UserId).Error; err != nil {
				return err
			}

			if err := tx.Exec("INSERT INTO transaction(`order_number`,`from_id`,`to_id`,`tx_hash`,`from_money`,`type`,`create_at`,`comment`,`to_money`,`from_coin`,`to_coin`,`usd`) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)",
				orderId,
				6,
				info.UserId,
				"",
				info.Money,
				6,
				models.MilliSecond(),
				"",
				info.Money,
				"VISA",
				"USD",
				1).Error; err != nil {
				return err
			}

			return nil
		})
	case 2: //拒绝
		return orm.Eloquent.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("UPDATE card_transfer SET status = 2 WHERE id = ?", info.Id).Error; err != nil {
				return err
			}

			if err := tx.Exec("UPDATE account SET frozen = frozen - ? WHERE user_id = ? AND coin='USD'", info.Money, info.UserId).Error; err != nil {
				return err
			}

			if err := tx.Exec("UPDATE card SET card_limit = card_limit + ? WHERE user_id = ? AND brand='VISA'", info.Money, info.UserId).Error; err != nil {
				return err
			}
			return nil
		})
	}
	return errors.New("no record")
}
