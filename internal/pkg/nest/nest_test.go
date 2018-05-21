package nest

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func setup() {
	configPath := os.Getenv("HOME") + "/.nerdnest"
	configName := "nerdnest"

	viper.SetEnvPrefix("nest")
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.SetDefault("units", "F")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}

func TestGetDeviceList(t *testing.T) {
	setup()
	device, err := GetDeviceList()
	if err != nil {
		t.Error(err)
	}
	if len(device) != 1 {
		t.Errorf("Should get one device")
	}
}

func TestRefresh(t *testing.T) {
	setup()
	device := Thermostat{}
	device.DeviceId = viper.GetString("mythermostat")
	device.Refresh()
	if device.DeviceId == "" {
		t.Errorf("DeviceId should not be nil")
	}
}
