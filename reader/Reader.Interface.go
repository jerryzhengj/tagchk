package reader

import "io"

const (
	Success byte = 0x10

)

type TagType uint
const (
	TagInfoType TagType = 0x01
	TagSummaryType TagType = 0x02
)

type Reader struct {
	// 配置
	opts Options

	// 连接
	conn io.ReadWriteCloser
	readBuffer []byte

	readBufferToken int
}

type TagInfo struct {
	AntID byte
	FreqAnt byte
	PC []byte
	EPC []byte
	RSSI byte
}

type TagSummary struct {
	AntID byte
	ReadRate []byte
	TotalRead []byte
	TotalDuration []byte
}

type InventoryInfo struct {
	TagInfo TagInfo
	TagSummary TagSummary
	PacketType TagType
}

type WriteResult struct{
	TagCount []byte
	DataLen byte
	PC []byte
	CRC []byte
	EPC []byte
	ErrCode byte
	Freq byte  //高6位是第一次读取的频点参数，低2位是天线号
	AntID byte //从Freq中解析出来的天线号
	WriteCount byte
}

type Callback func(data InventoryInfo)
type WriteCallback func(data WriteResult)

type Options struct {

	// 码率
	BaudRate uint

	// 串口
	PortName string

	// 天线列表
	Antennas []uint
}

type SerialHandler interface {

	readPacketHead() int

	readPacketBody(packetSize int) []byte

	// 发送数据参数为数据包body部分. 会自动添加head len check
	writeSerialBody([]byte) int

	// 发送数据
	writeSerial([]byte) int
}

type RFIDHandler interface {

	Open(*Options)

	Close()

	Inventory(Callback)

	Write(WriteCallback)
}