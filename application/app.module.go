package application


type TagInfo struct {
	AntID string      // 天线编号
	Epc string      //epc
	Frequency int //上报epc的次数
	Rssi int  //Rssi汇总值
	Ftime int64     //第一次上报epc的时间
	Ltime int64     //最后一次上报epc的时间
}