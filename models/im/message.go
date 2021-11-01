package im

import (
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"

	_ "time"
)

type Message struct {
	Id             int    `json:"id" gorm:"type:int(11) unsigned;primary_key"`                        // 消息序列号
	MsgType        int    `json:"msgType" gorm:"type:tinyint(1);"`                                    //
	MsgTitleZhCN   string `json:"msgTitleZhCN" gorm:"column:msg_title_zh-CN;type:varchar(256);"`      //
	MsgContentZhCN string `json:"msgContentZhCN" gorm:"column:msg_content_zh-CN;type:varchar(1024);"` //
	MsgTitleZhHK   string `json:"msgTitleZhHK" gorm:"column:msg_title_zh-HK;type:varchar(256);"`      //
	MsgContentZhHK string `json:"msgContentZhHK" gorm:"column:msg_content_zh-HK;type:varchar(1024);"`
	MsgTitleZhTW   string `json:"msgTitleZhTW" gorm:"column:msg_title_zh-TW;type:varchar(256);"` //
	MsgContentZhTW string `json:"msgContentZhTW" gorm:"column:msg_content_zh-TW;type:varchar(1024);"`
	MsgTitleEn     string `json:"msgTitleEn" gorm:"type:varchar(256);"`    //
	MsgTitleVi     string `json:"msgTitleVi" gorm:"type:varchar(256);"`    //
	MsgTitleTh     string `json:"msgTitleTh" gorm:"type:varchar(256);"`    //
	MsgContentEn   string `json:"msgContentEn" gorm:"type:varchar(1024);"` //
	MsgContentVi   string `json:"msgContentVi" gorm:"type:varchar(1024);"` //
	MsgContentTh   string `json:"msgContentTh" gorm:"type:varchar(1024);"` //
	Image          string `json:"image" gorm:"type:varchar(512);"`         //
	CreateAt       int    `json:"createAt" gorm:"type:bigint(13);"`        // 创建时间
	MsgUrl         string `json:"msgUrl" gorm:"type:varchar(256);"`        //
	DataScope      string `json:"dataScope" gorm:"-"`
	Params         string `json:"params"  gorm:"-"`
}

func (Message) TableName() string {
	return "message"
}

//创建
func (e *Message) Create() (Message, error) {
	var doc Message
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *Message) Get() (Message, error) {
	var doc Message

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取带分页
func (e *Message) GetPage(pageSize int, pageIndex int) ([]Message, int, error) {
	var doc []Message

	table := orm.Eloquent.Select("*").Table(e.TableName())

	// 数据权限控制(如果不需要数据权限请将此处去掉)

	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var count int
	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err

	}

	table.Count(&count)
	return doc, count, nil
}

// 更新
func (e *Message) Update(id int) (update Message, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除
func (e *Message) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&Message{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Message) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Message{}).Error; err != nil {
		return
	}
	Result = true
	return
}
