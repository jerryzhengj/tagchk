package reader

import (
	"bytes"
	"encoding/binary"
	log "github.com/jeanphorn/log4go"
	"github.com/jerryzhengj/tagchk/etc"
	"github.com/tarm/goserial"
	"time"
)

func Open(options *Options) (session Reader) {

	log.LOGGER("App").Info("Open connection %s:%d", options.PortName, options.BaudRate)

	//serialOptions := serial.OpenOptions{
	//	PortName: options.PortName,
	//	BaudRate: options.BaudRate,
	//	DataBits: 8,
	//	StopBits: 1,
	//	MinimumReadSize: 4,
	//}
	//port, err := serial.Open(serialOptions)

	serialOptions := &serial.Config{
		Name: options.PortName,
		Baud: int(options.BaudRate),
	}
	port, err := serial.OpenPort(serialOptions)

	if err != nil {
		log.LOGGER("App").Error(err)
		panic(err)
	}
	session = Reader{}
	session.opts = *options
	session.readBuffer = make([]byte, 128)
	session.conn = port

	_ = setPower(session, etc.Config.Antenna.Power)

	return session
}

func (session Reader) Close () {
	_ = session.conn.Close()
}

func (session Reader) Inventory(callback Callback) () {

	for _, antenna := range session.opts.Antennas {
		inventoryAntenna(session, antenna, callback)
	}
}

func (session Reader) FastInventory(callback Callback) () {

	fastInventory(session)
	handleInventoryResp(session, callback)
}

func (session Reader) WriteEPC(epc string,callback WriteCallback) () {

	writeEpc([]byte(epc),session)
	handleWriteResp(session, callback)
}

func inventoryAntenna(session Reader, antenna uint, callback Callback) {

	log.LOGGER("App").Debug("inventory antenna:%d", antenna)
	_ = selectAntenna(session, antenna)
	inventoryStart(session)
	handleInventoryResp(session, callback)
}

func handleInventoryResp(session Reader, callback Callback) {
	defer func() {
		if err := recover(); err != nil {
			log.LOGGER("App").Error(err)
			res := InventoryInfo{}
			res.PacketType = TagSummaryType
			summary := TagSummary{}
			res.TagSummary = summary
			callback(res)
		}
	}()
	for {
		inventoryInfo := handleInventory(session)
		callback(inventoryInfo)
		if inventoryInfo.PacketType == TagSummaryType {
			log.LOGGER("App").Debug("[handleInventoryResp]inventory complete")
			break
		}
	}
}


func handleWriteResp(session Reader, callback WriteCallback) {
	finished := false
	revNum := 0
	defer func() {
		if err := recover(); err != nil {
			log.LOGGER("App").Error("handleWriteResp receive error:%s",err)
			resultInfo  := WriteResult{}
			resultInfo.ResultType = 2
			resultInfo.ErrCode = 0x33
			callback(resultInfo)
		}
	}()
	var tmp int16
	var begin  = time.Now().Unix()
	for !finished{
		revNum++
		resultInfo := handleWrite(session)
		callback(resultInfo)

		binary.Read(bytes.NewBuffer(resultInfo.TagCount), binary.BigEndian, &tmp)
		current := time.Now().Unix()
		if revNum == int(tmp) || current - begin >= 5 {
			finished = true
		}
	}

	resultInfo := WriteResult{}
	resultInfo.ResultType = 2
	resultInfo.ErrCode = Success
	callback(resultInfo)
}

