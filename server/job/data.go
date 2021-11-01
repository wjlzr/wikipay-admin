package job

import (
	"wikipay-admin/models/monitor"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/rpc/client"
)

//获取币种计划任务时间
func getSetting(coin string, coinType int) *SettingInfo {
	setting := monitor.MonitorSetting{
		Coin: coin,
		Type: coinType,
	}
	result, err := setting.Get()
	if err == nil {
		info := SettingInfo{
			Coin:       result.Coin,
			Type:       result.Type,
			FeeAddress: result.FeeAddress,
			Day:        result.Day,
			Week:       result.Week,
			Hour:       result.Hour,
			Min:        result.Min,
			Max:        result.Max,
		}
		return &info
	}
	return nil
}

//同步区块交易数据
func syncData(coin string) bool {
	_, err := monitor.SyncAccountAndBalance(&client.BaseReq{
		Coin: coin,
	})
	if err != nil {
		return false
	}
	return true
}

//获取归集数据
func getDatas(coin string, coinType int, min, max float64) *MonitorInfo {
	var info MonitorInfo

	user := monitor.MonitorUserAddress{
		Coin: coin,
	}
	userAddress, err := user.GetDatas(min, max)
	if err == nil {
		addresses := make([]AddressInfo, 0)
		for _, u := range userAddress {
			a := AddressInfo{
				Address: u.Address,
				Amount:  u.Amount,
			}
			info.Total = xfloat64.Add(info.Total, u.Amount)
			addresses = append(addresses, a)
		}
		info.Address = addresses
	}

	mAddress := monitor.MonitorAddress{
		Coin:  coin,
		Type:  coinType,
		Ratio: 1,
	}
	maddrs, err := mAddress.Get()
	if err == nil {
		rs := make([]RatioInfo, 0)
		for _, m := range maddrs {
			ro := RatioInfo{
				Address: m.Address,
				Ratio:   m.Ratio,
				Amount:  xfloat64.Mul(info.Total, m.Ratio),
			}
			rs = append(rs, ro)
		}
		info.Ratio = rs
	}
	return &info
}
