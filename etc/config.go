package etc

import (
	"container/list"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/jeanphorn/log4go"
	"path/filepath"
)

var Config Configuration

var configName string

var envPath = new(list.List)

func LoadEnvs() {
	//log.LOGGER("App").Info("Load config")
	//log.LOGGER("App").Info("env size %d", envPath.Len())
	for e := envPath.Front(); e != nil; e = e.Next() {
		tmp := e.Value
		var path = tmp.(string)
		file := fmt.Sprintf("%s/%s", path, configName)
		log.LOGGER("App").Info("load env from %s", file)
		f, err := filepath.Abs(file)
		if err != nil {
			log.LOGGER("App").Info("Config file is not in %s", path)
		} else {
			if _, err := toml.DecodeFile(f, &Config); err != nil {
				log.LOGGER("App").Error("Config file is not in %s\n%s", path, err)
			} else {
				log.LOGGER("App").Info("%s. Version[%s]", Config.Name, Config.Version)
				break
			}
		}

	}
}

func AddEnvPath(env string) () {
	envPath.PushBack(env)
}

func SetConfigName(name string) () {
	configName = name
}
