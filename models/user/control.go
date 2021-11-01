package user

import (
	"fmt"
	"strconv"
	orm "wikipay-admin/database"
	"wikipay-admin/models/financial"
	"wikipay-admin/redis"
	"wikipay-admin/tools/config"
)

const (
	MaxEmailSendNum = 10
	MaxPaypwdErrNum = 5
	MaxSmsSendNum   = 10
)

var (
	smsNumKey       = "wikipay_phone_send_num_"
	emailNumKey     = "wikibank_email_send_num_"
	pwdErrKey       = "pwd_err_"
	lockUserKey     = "{user}_locked"
	unlockUserKey   = "{user}_unlocked"
	disableTradeKey = "disable_trade_"
)

//
type ControlReq struct {
	financial.Pagination
	Info string `json:"info" form:"info"`
}

//
type UserStatusReq struct {
	UserId int64 `json:"userId" binding:"required"`
}

//用户风控
type UserControlResp struct {
	Id           int64  `json:"userId"`
	AreaCode     string `json:"areaCode"`
	Phone        string `json:"phone"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	NickName     string `json:"nickName"`
	Email        string `json:"email"`
	EmailSendNum string `json:"emailSendNum"`
	PaypwdErrNum string `json:"payPwdErrNum"`
	SmsSendNum   string `json:"smsSendNum"`
	DisableTrade string `json:"disableTrade"`
	Ip           string `json:"ip"`
	Status       int    `json:"status"`
	CreateAt     int64  `json:"createAt"`
	Version      string `json:"version"`
}

//获取用户风控信息
func GetUserControlInfos(req *ControlReq) ([]UserControlResp, int, error) {
	var (
		infos, countInfos []UserControlResp
		where             string
	)

	sql := `SELECT id,area_code,phone,email,first_name,last_name,nick_name,locked as status,last_login_ip AS ip,version,create_at
			FROM user
			%s
			ORDER BY create_at DESC
			`
	if req.Info != "" {
		where = fmt.Sprintf("WHERE concat(phone,id) like '%s%s%s'", "%", req.Info, "%")
	}
	sql = fmt.Sprintf(sql, where)
	err := orm.Eloquent.Raw(sql).Find(&countInfos).Error
	count := len(countInfos)

	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d",
		sql,
		req.PageSize,
		req.PageNum)

	err = orm.Eloquent.Raw(sql).Find(&infos).Error
	if err != nil {
		return nil, 0, err
	}

	for k, v := range infos {
		infos[k].SmsSendNum = getValue(redis.ClusterClient().Get(fmt.Sprintf("%s+%s%s", smsNumKey, v.AreaCode, v.Phone)).Val())
		infos[k].EmailSendNum = getValue(redis.ClusterClient().Get(fmt.Sprintf("%s%s", emailNumKey, v.Email)).Val())
		infos[k].PaypwdErrNum = getValue(redis.ClusterClient().Get(fmt.Sprintf("%s%d", getRedisKey(pwdErrKey, true), v.Id)).Val())
		infos[k].DisableTrade = "N"
		disableTrade := fmt.Sprintf("%s%d", getRedisKey(disableTradeKey, true), v.Id)
		if redis.ClusterClient().TTL(disableTrade).Val() > 0 {
			infos[k].DisableTrade = "Y"
		}
	}

	return infos, count, nil
}

//重置
func Reset(req *UserStatusReq) error {
	var info UserControlResp
	err := orm.Eloquent.Raw("SELECT id, area_code, phone, email FROM user WHERE id = ?", req.UserId).Scan(&info).Error

	if info.Phone != "" {
		//手机短信
		smsKey := fmt.Sprintf("%s+%s%s", smsNumKey, info.AreaCode, info.Phone)
		smsSendNum := toInt(redis.ClusterClient().Get(smsKey).Val())
		if smsSendNum >= MaxSmsSendNum {
			redis.ClusterClient().Del(smsKey)
		}

		//邮箱
		emailKey := fmt.Sprintf("%s%s", emailNumKey, info.Email)
		emailSendNum := toInt(redis.ClusterClient().Get(emailKey).Val())
		if emailSendNum >= MaxEmailSendNum {
			redis.ClusterClient().Del(emailKey)
		}

		//支付密码
		payPwdKey := fmt.Sprintf("%s%d", getRedisKey(pwdErrKey, true), info.Id)
		paypwdErrNum := toInt(redis.ClusterClient().Get(getRedisKey(payPwdKey, true)).Val())
		if paypwdErrNum >= MaxPaypwdErrNum {
			redis.ClusterClient().Del(payPwdKey)
		}

		//禁止交易
		disableTrade := fmt.Sprintf("%s%d", getRedisKey(disableTradeKey, true), info.Id)
		if redis.ClusterClient().TTL(disableTrade).Val() > 0 {
			redis.ClusterClient().Del(disableTrade)
		}
	}
	return err
}

//更新用户状态
func UpdateStatus(req *UserStatusReq) error {
	var locked []int

	sql := fmt.Sprintf("SELECT locked FROM user WHERE id= %d ", req.UserId)
	err := orm.Eloquent.Raw(sql).Pluck("locked", &locked).Error
	if err != nil {
		return err
	}
	if len(locked) == 1 {
		if locked[0] == 0 {
			locked[0] = 1
			lockerUser(req.UserId, true)
		} else {
			locked[0] = 0
			lockerUser(req.UserId, false)
		}
		sql = fmt.Sprintf(`UPDATE user SET locked = %d WHERE id = %d `, locked[0], req.UserId)
		err = orm.Eloquent.Exec(sql).Error

		return err
	}
	return nil
}

//锁定用户
func lockerUser(userId int64, locked bool) {
	if locked {
		if !redis.ClusterClient().SMove(unlockUserKey, lockUserKey, userId).Val() {
			redis.ClusterClient().SAdd(lockUserKey, userId)
		}
		return
	}
	redis.ClusterClient().SMove(lockUserKey, unlockUserKey, userId).Val()
}

//将string转成int
func toInt(s string) int {
	i, err := strconv.Atoi(getValue(s))
	if err != nil {
		return 0
	}
	return i
}

//获取非空值，空为0
func getValue(s string) string {
	if len(s) == 0 {
		return "0"
	}
	return s
}

//
func getRedisKey(key string, isLine bool) string {
	if config.ApplicationConfig.Mode == "dev" {
		if isLine {
			return key + "_test_"
		}
		return key + "_test"
	}
	return key
}
