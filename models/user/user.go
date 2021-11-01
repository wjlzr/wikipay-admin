package user

import (
	"fmt"
	_ "time"
	"wikipay-admin/common"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type User struct {
	Id                int64  `json:"id" gorm:"type:bigint(10);primary_key" `                                //
	Password          string `json:"password" gorm:"type:varchar(200);" binding:"required"`                 // 密码
	PayPwd            string `json:"payPwd" gorm:"type:varchar(6);" binding:"required"`                     // 支付密码
	Phone             string `json:"phone" gorm:"type:varchar(20);" binding:"required"`                     // 手机号
	AreaFlag          string `json:"areaFlag" gorm:"type:varchar(100);"`                                    // 国家对应的旗帜
	AreaCode          string `json:"areaCode" gorm:"type:varchar(10);" binding:"required"`                  // 区号
	Email             string `json:"email" gorm:"type:varchar(60);" binding:"required"`                     // 邮箱
	EmailAuthStatus   int    `json:"emailAuthStatus" gorm:"type:tinyint(1);"`                               // 邮箱认证(0、未认证 1、已认证)
	FirstName         string `json:"firstName" gorm:"type:varchar(60);"`                                    // 姓
	LastName          string `json:"lastName" gorm:"type:varchar(60);"`                                     // 名
	NickName          string `json:"nickName" gorm:"type:varchar(50);"`                                     // 昵称
	Sex               int    `json:"sex" gorm:"type:tinyint(1) unsigned zerofill;"`                         // 性别（0-女,1-男）
	IdentityType      int    `json:"identityType" gorm:"type:tinyint(1) unsigned zerofill;"`                // 身份类型(1、身份证  2、护照  3、驾照)
	IdcardUrl         string `json:"idcardUrl" gorm:"type:varchar(200);"`                                   // 身份证正面
	FaceStatus        int    `json:"faceStatus" gorm:"type:tinyint(1);"`                                    // face++状态(0-未通过，1-已通过)
	AuthStatus        int    `json:"authStatus" gorm:"type:tinyint(1) unsigned zerofill;"`                  // 认证状态(0-未认证，1-认证中 2-已认证 3、认证失败)
	Avatar            string `json:"avatar" gorm:"type:varchar(200);"`                                      // 头像地址
	ChannelCode       string `json:"channelCode" gorm:"type:varchar(100);"`                                 // 渠道标识码
	GaSecretCode      string `json:"gaSecretCode" gorm:"type:varchar(40);"`                                 // 谷歌私钥
	Locked            int    `json:"locked" gorm:"type:tinyint(1);"`                                        // 是否冻结(0-未冻结，1-已冻结)
	Bank              string `json:"bank" gorm:"type:varchar(100);"`                                        // 银行名称
	BankNo            string `json:"bankNo" gorm:"type:varchar(50);"`                                       // 银行卡号
	LastLoginIp       string `json:"lastLoginIp" gorm:"type:varchar(50);"`                                  // 上次登录的ip
	CreateAt          int64  `json:"createAt" gorm:"type:bigint(13) unsigned zerofill;"`                    // 创建时间
	Platform          int    `json:"platform" gorm:"type:tinyint(1) unsigned zerofill;"`                    // 访问来源类型（0-PC端,1-Android端，2-IOS端）
	DeviceCode        string `json:"deviceCode" gorm:"type:varchar(50);"`                                   // 设备标识码
	DeviceInformation string `json:"deviceInformation" gorm:"type:varchar(100);"`                           // 设备信息
	InviteCode        string `json:"inviteCode" gorm:"type:varchar(100);"`                                  // 邀请码
	Version           string `json:"version" gorm:"type:varchar(30);"`                                      // 注册时版本
	LanguageCode      string `json:"languageCode" gorm:"type:varchar(20);"`                                 // 语言编码
	CountryCode       string `json:"countryCode" gorm:"type:varchar(20);"`                                  // 国家编码
	RoleId            int    `json:"roleId" gorm:"type:tinyint(1);"`                                        // 角色编号
	GroupId           int    `json:"groupId" gorm:"type:tinyint(2) unsigned zerofill;"`                     // 分组编号
	GradeId           int    `json:"gradeId" gorm:"type:tinyint(1) unsigned zerofill;"`                     // 等级编号
	UserType          int    `json:"userType" gorm:"type:tinyint(1) unsigned zerofill;" binding:"required"` // 用户类型(1-普通用户 2-商户)
	Content           string `json:"content" gorm:"type:varchar(200);"`                                     // 审核信息
	ImUsername        string `json:"imUsername" gorm:"type:varchar(200);"`                                  // 环信用户名
	ImPassword        string `json:"imPassword" gorm:"type:varchar(100);"`                                  // 环信密码
	Location          string `json:"location" gorm:"type:varchar(100);"`                                    // 地理位置
	Birthday          string `json:"birthday" gorm:"type:varchar(20);"`
	RealName          string `json:"realName" gorm:"type:varchar(40);"`
	DataScope         string `json:"dataScope" gorm:"-"`
	Params            string `json:"params"  gorm:"-"`
}

type SearchParams struct {
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}

//
func (User) TableName() string {
	return "user"
}

// 创建
func (e *User) Create() (User, error) {
	var doc User

	e.Id = tools.GenerateUserId()
	e.Password = tools.Encrypt(e.Password)
	e.CreateAt = tools.MilliSecond()
	if e.UserType == 0 {
		e.UserType = 1
	}
	if e.ChannelCode == "" {
		e.ChannelCode = "0"
	}

	if ims := tools.GetImInfo(); ims != nil {
		e.ImUsername = ims[0]
		e.ImPassword = ims[1]
	}

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e

	//coins := new(wallet.Coin).GetAccoutCoins()
	//for _, coin := range coins {
	for _, coin := range common.Coins {
		orm.Eloquent.Exec("INSERT INTO account(`user_id`,`coin`,`available`,`frozen`,`in_address`,`tag`,`create_at`,`type`) VALUES(?,?,?,?,?,?,?,?)",
			e.Id,
			coin,
			0,
			0,
			"",
			"",
			e.CreateAt,
			common.CoinAccountType[coin])
	}

	return doc, nil
}

//获取
func (e *User) Get() (User, error) {
	var doc User
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "id = ?", e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取User带分页
func (e *User) GetPage(pageSize, pageIndex int, info string, req SearchParams) ([]User, int, error) {
	var doc []User
	table := orm.Eloquent.Select("*").Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var where string
	if info != "" {
		where = fmt.Sprintf("concat(phone,email,id) like '%s%s%s'", "%", info, "%")
	}

	if req.StartTime != "" {
		table = table.Where("create_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		table = table.Where("create_at <= ?", req.EndTime)
	}

	if err := table.Where(where).
		Order("create_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	var count int
	table.Where(where).Count(&count)
	return doc, count, nil
}

//更新
func (e *User) Update(id int64) (update User, err error) {
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

//删除
func (e *User) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&User{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *User) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&User{}).Error; err != nil {
		return
	}
	Result = true
	return
}
