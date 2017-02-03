package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("nest")
	viper.SetConfigName("nerdnest")
	viper.AddConfigPath("$HOME/.nerdnest")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Please make sure you have created a config file")
		fmt.Println("See https://github.com/zpeters/nerdnest/ for examples")
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}

type Thermostat struct {
	Humidity     int
	DeviceId     string `json:"device_id"`
	Name         string
	TargetTempF  int    `json:"target_temperature_f"`
	AmbientTempF int    `json:"ambient_temperature_f"`
	HVACState    string `json:"hvac_state"`
	StructureID  string `json:"structure_id"`
}

type JResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
}

func (t *Thermostat) Refresh() {
	resp, err := http.Get("https://developer-api.nest.com/devices/thermostats/" + viper.GetString("mythermostat") + "?auth=" + viper.GetString("accesstoken"))
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

func (t Thermostat) SetTemp(temperature int) {
	t.SetAway("home")

	client := &http.Client{}
	body := fmt.Sprintf("{\"target_temperature_f\": %d}", temperature)

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
	return fmt.Sprintf("Name: %s\nCurrent Temp: %d\nTarget Temp: %d\nHumidity: %d\nState: %s\nDevice ID: %s\nStructure ID: %s",
		t.Name,
		t.AmbientTempF,
		t.TargetTempF,
		t.Humidity,
		t.HVACState,
		t.DeviceId,
		t.StructureID)
}

func register() {
	fmt.Printf("Registering...\n")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("1. Enter Product ID: ")
	productId, _ := reader.ReadString('\n')
	productId = strings.TrimSpace(productId)

	fmt.Print("2. Enter Client Secret: ")
	clientSecret, _ := reader.ReadString('\n')
	clientSecret = strings.TrimSpace(clientSecret)

	fmt.Print("3. Enter PIN: ")
	pin, _ := reader.ReadString('\n')
	pin = strings.TrimSpace(pin)

	fmt.Printf("4. Posting...\n")

	postUrl := "https://api.home.nest.com/oauth2/access_token"

	form := url.Values{}
	form.Add("client_id", productId)
	form.Add("client_secret", clientSecret)
	form.Add("code", pin)
	form.Add("grant_type", "authorization_code")

	resp, err3 := http.PostForm(postUrl, form)
	if err3 != nil {
		log.Fatal(err3)
	}
	defer resp.Body.Close()
	body, err4 := ioutil.ReadAll(resp.Body)
	if err4 != nil {
		log.Fatal(err4)
	}

	var jresp JResponse

	err5 := json.Unmarshal(body, &jresp)
	if err5 != nil {
		log.Fatal(err5)
	}

	if jresp.AccessToken == "" {
		log.Printf("Couldn't register\n")
		log.Fatalf("Last respones: %s", string(body))
	} else {
		fmt.Printf("Please set the following access code in your configuration\n")
		fmt.Printf("%s\n", jresp.AccessToken)
	}

}

func main() {
	var cmdAway = &cobra.Command{
		Use:   "away",
		Short: "'home' or 'away'",
		Run: func(cmd *cobra.Command, args []string) {
			t := Thermostat{}
			t.Refresh()
			t.SetAway(args[0])
		},
	}

	var cmdStatus = &cobra.Command{
		Use:   "status",
		Short: "Current Status",
		Run: func(cmd *cobra.Command, args []string) {
			t := Thermostat{}
			t.Refresh()
			fmt.Println(t)
		},
	}

	var cmdTemp = &cobra.Command{
		Use:   "temp",
		Short: "Set target temp",
		Run: func(cmd *cobra.Command, args []string) {
			temp, _ := strconv.Atoi(args[0])

			t := Thermostat{}
			t.Refresh()
			t.SetAway("home")
			t.SetTemp(temp)
		},
	}

	var cmdRegister = &cobra.Command{
		Use:   "register",
		Short: "Register with nest",
		Run: func(cmd *cobra.Command, args []string) {
			register()
		},
	}

	var rootCmd = &cobra.Command{Use: "nest"}
	rootCmd.AddCommand(cmdAway, cmdStatus, cmdTemp, cmdRegister)
	rootCmd.Execute()
}
