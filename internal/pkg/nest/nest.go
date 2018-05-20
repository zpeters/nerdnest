package nest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/spf13/viper"
)

type Thermostat struct {
	Humidity     int
	DeviceId     string `json:"device_id"`
	Name         string
	TargetTempF  int     `json:"target_temperature_f"`
	TargetTempC  float32 `json:"target_temperature_c"`
	AmbientTempF int     `json:"ambient_temperature_f"`
	AmbientTempC float32 `json:"ambient_temperature_c"`
	HVACState    string  `json:"hvac_state"`
	StructureID  string  `json:"structure_id"`
}

type JResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
}

func (t *Thermostat) Refresh() {
	resp, err := http.Get("https://developer-api.nest.com/devices/thermostats/" + t.DeviceId + "?auth=" + viper.GetString("accesstoken"))
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalf("Error: %v\n", err2)
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	for {
		if err := decoder.Decode(&t); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func (t Thermostat) SetAway(status string) {
	client := &http.Client{}

	bodystring := "{\"away\":\"" + status + "\"}"
	url := "https://developer-api.nest.com/structures/" + t.StructureID + "?auth=" + viper.GetString("accesstoken")

	req, err := http.NewRequest("PUT", url, strings.NewReader(bodystring))
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if resp.StatusCode == 307 {
		client := &http.Client{}
		req, err := http.NewRequest("PUT", resp.Header["Location"][0], strings.NewReader("{\"away\":\""+status+"\"}"))
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		_, err = client.Do(req)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	}
}

func (t Thermostat) SetTemp(temperature float64) {
	t.SetAway("home")

	client := &http.Client{}

	units := viper.GetString("units")

	var body string

	if units == "c" || units == "C" {
		body = fmt.Sprintf("{\"target_temperature_c\": %f}", temperature)
	} else {
		temperature_f := int(temperature)
		body = fmt.Sprintf("{\"target_temperature_f\": %d}", temperature_f)
	}
	req, err := http.NewRequest("PUT", "https://developer-api.nest.com/devices/thermostats/"+t.DeviceId+"?auth="+viper.GetString("accesstoken"), strings.NewReader(body))
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if resp.StatusCode == 307 {
		client := &http.Client{}
		url := resp.Header["Location"][0]
		req, err := http.NewRequest("PUT", url, strings.NewReader(body))
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		_, err = client.Do(req)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

	}
}

func (t Thermostat) String() string {
	units := viper.GetString("units")

	if units == "c" || units == "C" {
		return fmt.Sprintf("Name: %s\nCurrent Temp: %.1fC\nTarget Temp: %.1fC\nHumidity: %d\nState: %s\nDevice ID: %s\nStructure ID: %s",
			t.Name,
			t.AmbientTempC,
			t.TargetTempC,
			t.Humidity,
			t.HVACState,
			t.DeviceId,
			t.StructureID)
	}

	return fmt.Sprintf("Name: %s\nCurrent Temp: %dF\nTarget Temp: %dF\nHumidity: %d\nState: %s\nDevice ID: %s\nStructure ID: %s",
		t.Name,
		t.AmbientTempF,
		t.TargetTempF,
		t.Humidity,
		t.HVACState,
		t.DeviceId,
		t.StructureID)

}

// Generic get devices list from server
func GetDeviceList() (map[string]interface{}, error) {
	resp, err := http.Get("https://developer-api.nest.com/devices/" + "?auth=" + viper.GetString("accesstoken"))
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalf("Error: %v\n", err2)
	}

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq.Object("thermostats")
}
