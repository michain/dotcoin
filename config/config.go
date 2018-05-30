package config

import (
	"encoding/xml"
	"io/ioutil"
	"github.com/michain/dotcoin/util/filex"
	"github.com/michain/dotcoin/logx"
)

var (
	CurrentConfig  *AppConfig
	CurrentBaseDir string
)

// InitConfig load and marshal config file
func InitConfig(configFile string) *AppConfig {
	CurrentBaseDir = filex.GetCurrentDirectory()
	logx.Info("AppConfig::InitConfig init config [" + configFile + "] begin...")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		logx.Warn("AppConfig::InitConfig read [" + configFile + "] error - " + err.Error())
		panic(err)
	}

	var result AppConfig
	err = xml.Unmarshal(content, &result)
	if err != nil {
		logx.Warn("AppConfig::InitConfig read [" + configFile + "] unmarshal error - " + err.Error())
		panic(err)
	}

	result.configFile = configFile

	//init config base
	CurrentConfig = &result
	logx.Info("AppConfig::InitConfig init config [" + configFile + "] success")
	return CurrentConfig
}



