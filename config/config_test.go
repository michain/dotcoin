package config

import (
	"testing"
	"fmt"
)

func TestInitConfig(t *testing.T) {
	configFile := "app.conf"
	fmt.Println(fmt.Sprintf("%+v", InitConfig(configFile)))
}
