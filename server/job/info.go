package job

//归集结构
type MonitorInfo struct {
	Coin    string        //币种
	Address []AddressInfo //用户地址集合
	Ratio   []RatioInfo   //归庥占比集合
	Total   float64       //归集总数量
}

//地址结构
type AddressInfo struct {
	Address string  //用户地址
	Amount  float64 //数量
}

//归集比例
type RatioInfo struct {
	Address string  //归集到地址
	Ratio   float64 //归集比例
	Amount  float64 //归集占比数量
}

//时间信息
type SettingInfo struct {
	Coin       string  //币种
	Type       int     //类型
	Day        int     //日
	Week       int     //周
	Hour       int     //小时
	Max        float64 //最大
	Min        float64 //最小
	FeeAddress string
}
