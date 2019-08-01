package application

import (
	"fmt"
	cache "github.com/UncleBig/goCache"
	log "github.com/jeanphorn/log4go"
	"github.com/jerryzhengj/tagchk/etc"
	"github.com/jerryzhengj/tagchk/gpio"
	"github.com/jerryzhengj/tagchk/reader"
	"github.com/jerryzhengj/tagchk/util"
	"io/ioutil"
	"time"
)


var epcCache cache.Cache
var readerWorkStatus = 0 //0:free, 1:working; 2:read finished
var session reader.Reader
var epc2write string
var epcOfRead string
var startChk bool = false
var isWritingAndReading = false
var writeTotalNum = 0
var writeSuccessNum = 0

func onData(invInfo reader.InventoryInfo) {

	if invInfo.PacketType == reader.TagSummaryType{
		log.LOGGER("App").Debug("inventory complete")
		log.LOGGER("App").Debug("inventory summary:%v",invInfo.TagSummary)
		go repeatWrite()
		return
	}

	data := invInfo.TagInfo


	epc := string(data.EPC[0:len(epc2write)])
	//epc := hex.EncodeToString(data.EPC)
	//_, epc :=util.ParseEPC(data.EPC)
	//if len(epc) < 2{
	//	return
	//}
	//epc := string(data.EPC)
	////log.LOGGER("App").Info("report epc:%s,antId:%s",epc,fmt.Sprintf("%d", data.AntID))

	antIDStr := fmt.Sprintf("%d", data.AntID)
	epcOfRead = epc
	go gpio.TurnOnYellowLight()
	log.LOGGER("App").Info("found epc:%s",epc)
	//if _, found := epcCache.Get(epc); found {
	//
	//	//_ = epcCache.Replace(epc, foundEpc(epc,antIDStr), time.Duration(etc.Config.Cache.Expir) * time.Second)
	//	//log.LOGGER("App").Info("old epc:%s",epc)
	//	//return
	//} else {
	//	log.LOGGER("App").Info("found new epc:%s",epc)
	//	//go gpio.RingOnDelayOff()
	//	//epcCache.Set(epc, foundEpc(epc,antIDStr), time.Duration(etc.Config.Cache.Expir) * time.Second)
	//	//blink()
	//
	//}

	//将发现的标签存入缓存，以便做过期处理
	epcCache.Set(epc, foundEpc(epc,int(data.RSSI),antIDStr), time.Duration(etc.Config.Cache.Expir) * time.Second)
}

func foundEpc(epc string,rssi int,antIDStr string) TagInfo{
	var tagInfo TagInfo
	now := time.Now().Unix()
	v,ok := epcCache.Get(epc)
	if !ok {
		tagInfo =TagInfo{
			AntID:       antIDStr,
			Epc:         epc,
			Frequency:   1,
			Rssi:       rssi,
			Ftime:       now,
			Ltime:       now,
		}
	}else{
		tagInfo =v.(TagInfo)
		tagInfo.Ltime = now
		tagInfo.Rssi += rssi
		tagInfo.Frequency += 1
	}

	return tagInfo
}


func Run(c cache.Cache,opts reader.Options,fileName string) {
	log.LOGGER("App").Info("Application start")
	epcCache = c
	defer func() {
		if err := recover(); err != nil {
			log.LOGGER("App").Error(err)
			log.LOGGER("App").Info("Application run failed:%s",util.GetStack())
			session.Close()

		}
	}()
	_run(opts,fileName)
}

func _run(opts reader.Options,fileName string) {

	log.LOGGER("App").Info("Start process")
	session = reader.Open(&opts)

	gpio.TurnOnLight()
	time.Sleep(300 * time.Millisecond)

	gpio.TurnOffLight()
	time.Sleep(300 * time.Millisecond)

	gpio.TurnOnLightNormal()

	for{
		doMonitor(fileName)
		time.Sleep(time.Millisecond * 50)
	}
}


func doMonitor(fileName string){
	defer func(){
		if err := recover(); err != nil {
			log.LOGGER("App").Error(err)
			log.LOGGER("App").Info("Application run failed:%s",util.GetStack())
		}
	}()

	datas, err := ioutil.ReadFile(fileName)
	if err == nil{
		chkGPIOStatus(datas)
	}
}

func chkGPIOStatus(datas []byte){
	//log.LOGGER("App").Debug("datas[0]=%v,datas[1]=%v",datas[0],datas[1])
	if len(datas) < 2{
		return
	}

	if datas[util.PIN_1] == util.LOW{
		//log.LOGGER("App").Info("datas:%v",datas)
		if readerWorkStatus == 0 {
			//readerWorkStatus = 1
			readerWorkStatus = 2
			go gpio.TurnOnLightNormal()
			epc2write = time.Now().Format(util.TIME_LAYOUT)
			writeTotalNum = 0
			writeSuccessNum = 0
			startChk = true
			repeatWrite()
		}
	}
	if datas[util.PIN_2] == util.LOW{
		//log.LOGGER("App").Info("datas:%v",datas)
		if readerWorkStatus == 2{
			readerWorkStatus = 0
			startChk = false
			log.LOGGER("App").Info("start to report")
			report()
		}
	}
}


func repeatWrite(){
	if startChk {
		writeEpc(epc2write)
	}
}

func report(){
	lightRed := false
	if epcOfRead != epc2write {
		lightRed = true
	}
	if v, found := epcCache.Get(epcOfRead); found {
		log.LOGGER("App").Info("report of %s:%v",epcOfRead,v)
	}else{
		lightRed = true
	}

	if lightRed{
		go gpio.TurnOnRedLight()
	}

	log.LOGGER("App").Info("total write times:%d, success times:%d",writeTotalNum,writeSuccessNum)

}


func inventory() {
	epcOfRead = ""
	//session.Inventory(onData)
	session.FastInventory(onData)
	//time.Sleep(time.Millisecond * 100)
}

func writeEpc(epc string){
	writeTotalNum++
	session.WriteEPC(epc,onWriteResult)
}

func onWriteResult(invInfo reader.WriteResult) {
	log.LOGGER("App").Debug("write result:%v",invInfo)
	if invInfo.ErrCode == reader.Success{
		writeSuccessNum++
	}
	inventory()
}

