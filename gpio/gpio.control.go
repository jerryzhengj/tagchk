package gpio

import (
	"bytes"
	"github.com/jerryzhengj/tagchk/etc"
	log "github.com/jeanphorn/log4go"
	"os"
	"time"
)

// red/yellow/green/bee~
const Off = 0x00
const On = 0x01
type Light struct {
	Red []byte
	Green []byte
	Yellow []byte
	On []byte
	Off []byte
	Normal []byte
	BlinkOn []byte
	BlinkOff []byte
}
var lightState = Light{
	Red: []byte{On, Off, Off},
	Green: []byte{Off, On, Off},
	Yellow: []byte{Off, Off, On},
	On: []byte{On, On, On},
	Off: []byte{Off, Off, Off},
	Normal: []byte{Off, On, Off},
	BlinkOn: []byte{Off, On, On},
	BlinkOff: []byte{Off, Off, On},
}

var lLightstatus []byte
var lRingstatus []byte

func InitLightStatus(status *[]byte){
	lLightstatus = *status
	lRingstatus =  make([]byte,1)
}

func updateLightStatus(status []byte) {
	file, err := os.OpenFile("cgpio.bin", os.O_RDWR | os.O_CREATE, 0755)
	if err != nil {
		return
	}

	defer file.Close()
	if bytes.Equal(lLightstatus,status) {
		return
	}
	lLightstatus = status
	//log.LOGGER("App").Info("UpdateLightStatus:", util.BytesToHex(lLightstatus))
	_, err = file.WriteAt(lLightstatus, 0)

	if err != nil {
		log.LOGGER("App").Info(err)
	}

	//log.LOGGER("App").Info("Write light size:", len(lstatus))
}

func updateRingStatus(status byte) {
	file, err := os.OpenFile("cgpio.bin", os.O_RDWR | os.O_CREATE, 0755)
	if err != nil {
		return
	}

	defer file.Close()
	if lRingstatus[0] == status {
		return
	}
	lRingstatus[0] = status
	_, err = file.WriteAt(lRingstatus, 3)

	if err != nil {
		log.LOGGER("App").Info(err)
	}

	//log.LOGGER("App").Info("Write light size:", len(lstatus))
}


func RingOnDelayOff(){
	log.LOGGER("App").Info("RingOnDelayOff:%v", lRingstatus)
	if lRingstatus[0] == On {
		return
	}

	updateRingStatus(On)
	go RingDelayOff(etc.Config.Lights.DelayOff)
}



func RingDelayOff(delaySeconds int64){
	//log.LOGGER("App").Info("RingDelayOff:delaySeconds=%d", delaySeconds)
	time.Sleep(time.Duration(delaySeconds) * time.Second)
	log.LOGGER("App").Info("RingDelayOff:%v", lRingstatus)
	if lRingstatus[0] == Off {
		return
	}

	updateRingStatus(Off)
}

func TurnOnLight(){
	if bytes.Equal(lLightstatus,lightState.On) {
		return
	}
	log.LOGGER("App").Info("TurnOnLight:from %v to %v", lLightstatus,lightState.On)
	updateLightStatus(lightState.On)
}

func TurnOffLight(){
	if bytes.Equal(lLightstatus,lightState.Off) {
		return
	}
	log.LOGGER("App").Info("TurnOffLight:from %v to %v", lLightstatus,lightState.Off)
	updateLightStatus(lightState.Off)
}

func TurnOnLightNormal(){
	if bytes.Equal(lLightstatus,lightState.Normal) {
		return
	}
	log.LOGGER("App").Info("TurnOnLightNormal:from %v to %v", lLightstatus,lightState.Normal)
	updateLightStatus(lightState.Normal)
}

func TurnOnYellowLight(){
	if bytes.Equal(lLightstatus,lightState.Yellow) {
		return
	}
	log.LOGGER("App").Info("TurnOnYellowLight:from %v to %v", lLightstatus,lightState.Yellow)
	updateLightStatus(lightState.Yellow)
}

func TurnOnRedLight(){
	if bytes.Equal(lLightstatus,lightState.Red) {
		return
	}
	log.LOGGER("App").Info("TurnOnYellowLight:from %v to %v", lLightstatus,lightState.Red)
	updateLightStatus(lightState.Red)
}

func DelayNormal(delaySeconds int64){
	time.Sleep(time.Duration(delaySeconds) * time.Second)
	if bytes.Equal(lLightstatus,lightState.Green) {
		return
	}
	log.LOGGER("App").Info("DelayNormal:from %v to %v", lLightstatus,lightState.Green)
	updateLightStatus(lightState.Green)
}


