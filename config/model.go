package config

import (
	"encoding/xml"
)

// AppConfig app's config
type AppConfig struct {
	XMLName    xml.Name     `xml:"config"`
	MiningSet  MiningSet	`xml:"mining"`
	RpcSet     RpcSet	`xml:"rpc"`
	configFile string
}

// MiningSet mining config
type MiningSet struct {
	Enabled   bool `xml:"enabled,attr"`
}

// RpcSet rpc server config
type RpcSet struct {
	Enabled   bool `xml:"enabled,attr"`
	Port	string	`xml:"port,attr"`
}


