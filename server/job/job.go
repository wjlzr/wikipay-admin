package job

import (
	"strconv"

	"github.com/jasonlvhit/gocron"
)

var (
	coins = []string{"ETH", "USDT-ERC20", "USDT-OMNI"}
)

//开始服务
func MonitorJobsStart() {
	for k, coin := range coins {
		for i := 1; i < 3; i++ {
			info := getSetting(coin, i)
			switch k {
			case 0, 1:
				jobDo(info.Day, info.Week, info.Hour, ethMonitor, info)
			case 2:
				jobDo(info.Day, info.Week, info.Hour, omniMonitor, info)
			}
		}
	}
}

//启动宣时任务
func jobDo(day, week, hour int, jobFun interface{}, params ...interface{}) {
	if hour <= 0 {
		return
	}
	strHour := strconv.Itoa(hour) + ":00"
	if day > 0 {
		gocron.Every(uint64(day)).Days().At(strHour).Do(jobFun, params...)
	} else if week > 0 {
		gocron.Every(uint64(day)).Weeks().At(strHour).Do(jobFun, params...)
	} else {
		gocron.Every(1).At(strHour).Do(jobFun, params...)
	}
}
