package main

import (
	log "github.com/jeanphorn/log4go"
	"github.com/jerryzhengj/tagchk/application"
	"github.com/jerryzhengj/tagchk/etc"
	"github.com/jerryzhengj/tagchk/gpio"
	"github.com/jerryzhengj/tagchk/reader"
	"github.com/jerryzhengj/tagchk/util"
	"time"
)

func main() {
	log.LoadConfiguration("log.json")
	defer log.Close()

	etc.SetConfigName("config.toml")
	etc.AddEnvPath(".")
	etc.AddEnvPath("/etc/server")
	etc.LoadEnvs()
	if len(etc.Config.Version) == 0 {
		return
	}

	opts := reader.Options{
		BaudRate: etc.Config.Connector.Port,
		PortName: etc.Config.Connector.Name,
		Antennas: etc.Config.Antenna.Antennas,
	}

	log.LOGGER("App").Debug("opts:%v",opts)

	lightStatus := make([]byte,3)
	for i:=0; i < len(lightStatus); i++ {
		lightStatus[i] = gpio.Off
	}
	gpio.InitLightStatus(&lightStatus)

	//初始化一个缓存，application.app读到标签后读到的标签写入到缓存。
	//worker.Monitor准备上报标签前会从此缓存读取标签判断是否过期
	c :=util.NewCache(time.Duration(etc.Config.Cache.Expir), 30 * time.Second)

	for {
		application.Run(*c,opts,etc.Config.Listen.File)
		time.Sleep(time.Second * 2)
	}
}
