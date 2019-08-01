package util

import (
	"bytes"
	"regexp"
	"runtime"
	"strconv"
)

const (
	LOW = iota
	HIGH
	UNKNOWN
)

const (
	PIN_1 = iota
	PIN_2
)

const DATE_LAYOUT = "20060102150405"
const TIME_LAYOUT = "150405"

var reg  *regexp.Regexp

func init(){
	//ASCII控制字符
	reg = regexp.MustCompile(`\p{C}`)
}


func Nagation(a uint8) uint8 {
	return ^a + 1
}

func AbsInt(a int64) int64 {
	if a > 0 {
		return a
	}
	return a * -1
}

func MaxInt(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

// 计算校验值
func CRC(a []byte) uint8 {

	var sum uint8 = 0
	for i := 0; i < len(a); i++ {
		value := uint8(a[i])
		sum = sum + value
	}

	return Nagation(sum)
}

func BytesToHex(data []byte) string {
	buffer := new(bytes.Buffer)
	for index, b := range data {

		s := strconv.FormatInt(int64(b&0xff), 16)
		if index > 0 {
			buffer.WriteString(" ")
		}
		if len(s) == 1 {
			buffer.WriteString("0x0")
		} else {
			buffer.WriteString("0x")
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

func ByteToHex(data byte) string {

	s := strconv.FormatInt(int64(data&0xff), 16)
	if len(s) == 1 {
		s = "0x0" + s
	} else {
		s = "0x" + s
	}

	return s
}

func GetStack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}

func ParseEPC(raw []byte)(version ,epc string){
	length := int(raw[0]&0x0f)
	//fmt.Printf("ParseEPC length:%d\n",length)
	//log.LOGGER("App").Info("ParseEPC length:%d",length)
	if length <= 0 || length > len(raw) - 1{
		return "",""
	}
	version = strconv.FormatInt(int64(raw[0]>>4), 10)
	//epc = string(raw[1:length+1])

	epc = reg.ReplaceAllString(string(raw[1:length+1]),"")

	return version,epc
}
