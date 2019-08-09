package etc

type Configuration struct {
	Name string
	Version string
	DeviceNo string

	// 监听GPIO
	Listen listen

	// 天线
	Antenna antenna
	// 串口
	Connector connector

	Cache cache
	Lights light
}

type connector struct {
	Name string
	Port uint
}


type antenna struct {
	Antennas []uint
	Power uint
	StayDuration int64
}

type cache struct {
	Expir int64
}

type listen struct {
	Pin []int
	File string
	Interval int
}

type light struct {
	DelayOff int64
	PuffOffAfter int64
}