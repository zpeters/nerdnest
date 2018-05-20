package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/zpeters/nerdnest/internal/pkg/nest"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// injected via makefile
var VERSION string
var COMMIT string
var BRANCH string

// some config settings
var configPath = os.Getenv("HOME") + "/.nerdnest"
var configName = "nerdnest"
var fullConfigPath = configPath + "/" + configName + ".toml"

func init() {
	viper.SetEnvPrefix("nest")
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.SetDefault("units", "F")
	err := viper.ReadInConfig()

	if err != nil && os.Args[1] != "register" {
		fmt.Println("Please make sure you have created a config file")
		fmt.Println("See https://github.com/zpeters/nerdnest/ for examples")
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}

// List devices to user
func ListDevices() {

	obj, err := nest.GetDeviceList()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range obj {
		fmt.Printf("%s: %s\n", v.(map[string]interface{})["where_name"], k)
	}

}

// Choose default device
func SetDefaultDevice() {
	obj, err := nest.GetDeviceList()
	if err != nil {
		log.Fatal(err)
	}

	var defaultDeviceId string

	if len(obj) == 1 {
		fmt.Printf("Only one device, setting this to the default\n")
		for k, v := range obj {
			fmt.Printf("%s - %s\n", v.(map[string]interface{})["where_name"], k)
			defaultDeviceId = k
		}
	} else {
		fmt.Printf("Found multiple devices:\n")
		for k, v := range obj {
			fmt.Printf("%s - %s\n", v.(map[string]interface{})["where_name"], k)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter default device ID from above list:")
		defaultDeviceId, _ := reader.ReadString('\n')
		defaultDeviceId = strings.TrimSpace(defaultDeviceId)
	}

	// Now write config file
	var configString string
	var configKeys = []string{"accesstoken", "units"}
	for _, k := range configKeys {
		configString = configString + k + " = \"" + viper.GetString(k) + "\"\n"
	}
	configString = configString + "mythermostat = \"" + defaultDeviceId + "\"\n"

	err1 := ioutil.WriteFile(fullConfigPath, []byte(configString), 0640)
	if err1 != nil {
		log.Fatal(err1)
	}
}

func Register() {
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

	var jresp nest.JResponse

	err5 := json.Unmarshal(body, &jresp)
	if err5 != nil {
		log.Fatal(err5)
	}

	if jresp.AccessToken == "" {
		log.Printf("Couldn't register\n")
		log.Fatalf("Last responses: %s", string(body))
	} else {
		if _, err := os.Stat(fullConfigPath); err == nil {
			fmt.Printf("Configuration file already exists. Please set the following access code in your configuration\n")
			fmt.Printf("%s\n", jresp.AccessToken)

		} else {
			// file does not exist
			fmt.Printf("Configuration file did not exist, creating\n")
			err1 := os.MkdirAll(configPath, os.FileMode(0750))
			if err1 != nil {
				log.Fatal(err1)
			}

			fmt.Print("5. Enter temperature units (C or F): ")
			units, _ := reader.ReadString('\n')
			units = strings.TrimSpace(units)

			err := ioutil.WriteFile(fullConfigPath, []byte("accesstoken = \""+jresp.AccessToken+"\"\nunits = \""+units+"\"\n"), 0640)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	var cmdAway = &cobra.Command{
		Use:   "away",
		Short: "'home' or 'away'",
		Run: func(cmd *cobra.Command, args []string) {
			t := nest.Thermostat{}
			if len(args) == 2 {
				t.DeviceId = args[1]
			} else {
				t.DeviceId = viper.GetString("mythermostat")
			}
			t.Refresh()
			t.SetAway(args[0])
		},
	}

	var cmdStatus = &cobra.Command{
		Use:   "status",
		Short: "Current Status",
		Run: func(cmd *cobra.Command, args []string) {
			t := nest.Thermostat{}

			if len(args) == 1 {
				t.DeviceId = args[0]
			} else {
				t.DeviceId = viper.GetString("mythermostat")
			}

			t.Refresh()
			fmt.Println(t)
		},
	}

	var cmdTemp = &cobra.Command{
		Use:   "temp",
		Short: "Set target temp",
		Run: func(cmd *cobra.Command, args []string) {
			temp, _ := strconv.ParseFloat(args[0], 32)

			t := nest.Thermostat{}

			if len(args) == 2 {
				t.DeviceId = args[1]
			} else {
				t.DeviceId = viper.GetString("mythermostat")
			}

			t.Refresh()
			t.SetAway("home")
			t.SetTemp(temp)
		},
	}

	var cmdRegister = &cobra.Command{
		Use:   "register",
		Short: "Register with nest",
		Run: func(cmd *cobra.Command, args []string) {
			Register()
		},
	}

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "List devices",
		Run: func(cmd *cobra.Command, args []string) {
			ListDevices()
		},
	}

	var cmdSetDefault = &cobra.Command{
		Use:   "setdefault",
		Short: "Set default thermostat",
		Run: func(cmd *cobra.Command, args []string) {
			SetDefaultDevice()
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Show build information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", VERSION)
			fmt.Printf("Branch: %s\n", BRANCH)
			fmt.Printf("Commit: %s\n", COMMIT)
		},
	}

	var rootCmd = &cobra.Command{Use: "nerdnest"}
	rootCmd.AddCommand(cmdAway, cmdStatus, cmdTemp, cmdRegister, cmdList, cmdSetDefault, cmdVersion)
	rootCmd.Execute()
}
