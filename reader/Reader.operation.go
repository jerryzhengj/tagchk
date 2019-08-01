package reader

import (
	"encoding/hex"
	log "github.com/jeanphorn/log4go"
	"time"
)

//
func setPower(session Reader, pow uint) (err error) {

	//log.LOGGER("App").Debug("config power %d", pow)
	power := byte(pow)
	if power > 0x21 {
		power = 0x21
	}
	var pCmd byte= 0x76
	//p := make([]byte, 3)
	p := make([]byte, 6)
	p[0] = 0x01 // addr
	p[1] = pCmd
	p[2] = 0x00
	p[3] = 0x00
	p[4] = 0x00
	p[5] = 0x00
	for _, antenna := range session.opts.Antennas {
		if antenna == 0 {
			p[2] = power
		}
		if antenna == 1 {
			p[3] = power
		}
		if antenna == 2 {
			p[4] = power
		}
		if antenna == 3 {
			p[5] = power
		}
	}
	log.LOGGER("App").Debug("setPower command:%v", p)
	session.writeSerialBody(p)
	time.Sleep(time.Millisecond * 100)

	// receive
	packetSize := session.readPacketHead()
	log.LOGGER("App").Debug("setPower packetSize:%d", packetSize)

	packetBody := session.readPacketBody(packetSize)
	log.LOGGER("App").Debug("setPower packetBody:%v", packetBody)

	if packetBody[1] == pCmd && packetBody[2] != Success {
		log.LOGGER("App").Error("set power error " + hex.EncodeToString(packetBody[2:3]))
		panic("set power error " + hex.EncodeToString(packetBody[2:3]))
	}

	return err
}

func selectAntenna(session Reader, antenna uint) (err error) {

	log.LOGGER("App").Debug("select antenna")
	p := make([]byte, 3)
	p[0] = 0x01 // addr
	p[1] = 0x74
	p[2] = byte(antenna)
	session.writeSerialBody(p)
	time.Sleep(time.Millisecond * 50)

	// receive
	packetSize := session.readPacketHead()

	packetBody := session.readPacketBody(packetSize)

	if packetBody[2] != Success {
		log.LOGGER("App").Error("select antenna error " + hex.EncodeToString(packetBody[2:3]))
		panic("select antenna error " + hex.EncodeToString(packetBody[2:3]))
	}

	return err
}

func writeEpc(data []byte,session Reader){
	p := make([]byte, 9+len(data))
	p[0] = 0x01 // addr
	p[1] = 0x82 //cmd_write
	p[2] = 0x00
	p[3] = 0x00
	p[4] = 0x00
	p[5] = 0x00
	p[6] = 0x01 // epc
	p[7] = 0x02 // WordAdd
	p[8] = byte(len(data)/2) // WordCnt
	for i:=0; i<len(data); i++{
		p[9+i] = data[i]
	}

	session.writeSerialBody(p)
}

func inventoryStart(session Reader) {

	//log.LOGGER("App").Debug("invoke inventory")
	p := make([]byte, 3)
	p[0] = 0x01 // addr
	p[1] = 0x89 //real time inventory
	p[2] = 0xFF
	session.writeSerialBody(p)
	//time.Sleep(time.Millisecond * 100)
}

func fastInventory(session Reader) {
	//log.LOGGER("App").Debug("invoke fast inventory")

	var times byte = 0x03
	var antennaSleep byte = 0x10
	var repeat byte = 3
	p := make([]byte, 12)
	p[0] = 0x01 // addr
	p[1] = 0x8A //fast inventory
	p[2] = 0xFF// A
	p[3] = times
	p[4] = 0xFF // B
	p[5] = times
	p[6] = 0xFF// C
	p[7] = times
	p[8] = 0xFF // D
	p[9] = times
	p[10] = antennaSleep
	p[11] = repeat

	for _, antenna := range session.opts.Antennas {
		if antenna == 0 {
			p[2] = 0x00
		}
		if antenna == 1 {
			p[4] = 0x01
		}
		if antenna == 2 {
			p[6] = 0x02
		}
		if antenna == 3 {
			p[8] = 0x03
		}
	}
	session.writeSerialBody(p)
	//time.Sleep(time.Millisecond * 100)
}

func handleInventory(session Reader) (res InventoryInfo) {

	//log.LOGGER("App").Debug("handle inventory")
	res = InventoryInfo{}
	packetSize := session.readPacketHead()

	//log.LOGGER("App").Debug("packet size: %d", packetSize)

	if packetSize == 0x04 {// 异常

		packetBody := session.readPacketBody(packetSize)
		log.LOGGER("App").Error("Inventory error " + hex.EncodeToString(packetBody[2:3]))
		panic("Inventory error " + hex.EncodeToString(packetBody[2:3]))
	}else if packetSize == 0x05{//天线未连接错误
		packetBody := session.readPacketBody(packetSize)
		log.LOGGER("App").Error("Ant connection error " + hex.EncodeToString(packetBody[2:3]))
		panic("Ant connection error " + hex.EncodeToString(packetBody[2:3]))
	}

	if packetSize == 0x0A {// 结束包
		res.PacketType = TagSummaryType

		packetBody := session.readPacketBody(packetSize)

		summary := TagSummary{}
		if packetBody[2] == 0x8A {//fast inventory response
			summary.TotalRead = packetBody[2:5]
			summary.TotalDuration= packetBody[5:9]
		}else{//real time inventory response
			summary.AntID = packetBody[2]
			summary.ReadRate = packetBody[3:5]
			summary.TotalRead = packetBody[5:9]
		}

		res.TagSummary = summary
	} else {// 标签包
		res.PacketType = TagInfoType
		packetBody := session.readPacketBody(packetSize)

		tagInfo := TagInfo{}
		tagInfo.FreqAnt = packetBody[2]
		tagInfo.AntID = (tagInfo.FreqAnt<<6)>>6
		tagInfo.PC = packetBody[3:5]
		tagInfo.EPC = packetBody[5:packetSize - 2]
		tagInfo.RSSI = packetBody[packetSize - 2]

		res.TagInfo = tagInfo
	}
	//log.LOGGER("App").Debug("packet type: %d", res.PacketType)

	return res
}


func handleWrite(session Reader) (res WriteResult) {

	//log.LOGGER("App").Debug("handle handleWrite")
	res = WriteResult{}
	packetSize := session.readPacketHead()

	//log.LOGGER("App").Debug("packet size: %d", packetSize)

	if packetSize == 0x04 {// 异常

		packetBody := session.readPacketBody(packetSize)
		log.LOGGER("App").Error("Write error " + hex.EncodeToString(packetBody[2:3]))
		panic("Write error " + hex.EncodeToString(packetBody[2:3]))
	}


	packetBody := session.readPacketBody(packetSize)
	log.LOGGER("App").Debug("handleWrite packetBody: %v", packetBody)
	tagInfo := WriteResult{}

	tagInfo.TagCount = packetBody[2:4]
	tagInfo.DataLen = packetBody[4]
	tagInfo.PC = packetBody[5:7]
	tagInfo.EPC = packetBody[7:packetSize - 6]
	tagInfo.CRC = packetBody[packetSize - 6:packetSize - 4]
	tagInfo.ErrCode = packetBody[packetSize - 4]
	tagInfo.Freq = (packetBody[packetSize - 3]>>2)<<2
	tagInfo.AntID = (packetBody[packetSize - 3]<<6)>>6
	tagInfo.WriteCount = packetBody[packetSize - 2]

	return tagInfo
}