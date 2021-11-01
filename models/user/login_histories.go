package user

import (
	_ "time"
	orm "wikipay-admin/database"
)

//
type LoginHistories struct {
	Id        int    `json:"id" gorm:"type:int(10);primary_key"` // 序列号
	UserId    int64  `json:"userId" gorm:"type:int(10);"`        // 用户序列号
	CreateAt  int64  `json:"createAt" gorm:"type:bigint(13);"`   // 登陆时间
	Status    int    `json:"status" gorm:"type:tinyint(1);"`     // 状态(0-登陆成功，1-登陆失败)
	Equipment string `json:"equipment" gorm:"type:varchar(50);"` // 终端设备
	Ip        string `json:"ip" gorm:"type:varchar(32);"`        // ip地址
	Title     string `json:"title" gorm:"type:varchar(100);"`    // 标题
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

//表名
func (LoginHistories) TableName() string {
	return "login_histories"

}

// 创建LoginHistories
func (e *LoginHistories) Create() (LoginHistories, error) {
	var doc LoginHistories
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取LoginHistories
func (e *LoginHistories) Get() (LoginHistories, error) {
	var doc LoginHistories

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "id = ?", e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取LoginHistories带分页
func (e *LoginHistories) GetPage(pageSize int, pageIndex int) ([]LoginHistories, int, error) {
	var doc []LoginHistories
	table := orm.Eloquent.Select("*").Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	// dataPermission := new(models.DataPermission)
	// dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	// table, err := dataPermission.GetDataScope(e.TableName(), table)
	// if err != nil {
	// 	return nil, 0, err
	// }

	var count int
	if err := table.Order("create_at DESC").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

// 更新LoginHistories
func (e *LoginHistories) Update(id int) (update LoginHistories, err error) {
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

// 删除LoginHistories
func (e *LoginHistories) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&LoginHistories{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *LoginHistories) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&LoginHistories{}).Error; err != nil {
		return
	}
	Result = true
	return
}
