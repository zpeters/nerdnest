# nerdnest
Control your Nest device from the command line

[![Build Status](https://travis-ci.org/zpeters/nerdnest.svg?branch=master)](https://travis-ci.org/zpeters/nerdnest)
[![Go Report Card](https://goreportcard.com/badge/github.com/zpeters/nerdnest)](https://goreportcard.com/report/github.com/zpeters/nerdnest)
[![Github All Releases](https://img.shields.io/github/downloads/zpeters/nerdnest/total.svg?style=plastic)](https://github.com/zpeters/nerdnest)

# Setup

1. Enable developer mode and allow access to your nest
  - https://github.com/zpeters/nerdnest/wiki/Nest-Developer-Account
2. Register nerdnest with your developer account
  - https://github.com/zpeters/nerdnest/wiki/Registering-nerdnest-with-your-account
3. Choose the device you want to control
  - https://github.com/zpeters/nerdnest/wiki/Select-your-Nest-Device
4. Make sure your settings from step 2 and 3 are saved into your nerdnest.toml file.  This file can be in:
  - $HOME/.nerdnest/nerdnest.toml
  - CURRENTDIRECTORY/nerdnest.toml
  - _Submit a new issue for additional paths_

# Usage
```
Usage:
  nerdnest [command]

Available Commands:
  away        'home' or 'away'
  status      Current Status
  temp        Set target temp
  register    Register with nest
  list        List devices

Flags:
  -h, --help   help for nest

Use "nerdnest [command] --help" for more information about a command.
```

# Current Status
```
./nerdnest status
Name: Nest
Current Temp: 69
Target Temp: 69
Humidity: 45
State: off
Device ID: KoTA9-raY9xdYrYY036u2rgaeP_lJ-mg
Structure ID: Suha_CVEVHdOreQFLWC-XlHaPXSRHcEwOb8dKkwYIjcVN0XCBSnKLQ
```

# Set away status
No output is sent after the command runs
```
./nerdnest away home
./nerdnest away home
 ```

# Set temperature
```
./nerdnest temp 70
./nerdnest status
Name: Nest
Current Temp: 69
Target Temp: 70
Humidity: 45
State: heating
Device ID: KoTA9-raY9xdYrYY036u2rgaeP_lJ-mg
Structure ID: Suha_CVEVHdOreQFLWC-XlHaPXSRHcEwOb8dKkwYIjcVN0XCBSnKLQ

./nerdnest temp 68
./nerdnest status
Name: Nest
Current Temp: 69
Target Temp: 68
Humidity: 45
State: off
Device ID: KoTA9-raY9xdYrYY036u2rgaeP_lJ-mg
Structure ID: Suha_CVEVHdOreQFLWC-XlHaPXSRHcEwOb8dKkwYIjcVN0XCBSnKLQ
```

# Set default device
```
./nerdnest setdefault
Nest: raY9xdYrYY036u2rgaeP_lJ
Kitchen: ku34h5kjefhkdsjfhsdf
Enter default device ID from above list: ku34h5kjefhkdsjfhsdf 

./nerdnest status
Name: Kitchen
Current Temp: 69
Target Temp: 70
Humidity: 45
State: heating
Device ID: ku34h5kjefhkdsjfhsdf
Structure ID: Suha_CVEVHdOreQFLWC-XlHaPXSRHcEwOb8dKkwYIjcVN0XCBSnKLQ

./nerdnest status ku34h5kjefhkdsjfhsdf
Name: Nest
Current Temp: 69
Target Temp: 68
Humidity: 45
State: off
Device ID: KoTA9-raY9xdYrYY036u2rgaeP_lJ-mg
Structure ID: Suha_CVEVHdOreQFLWC-XlHaPXSRHcEwOb8dKkwYIjcVN0XCBSnKLQ
```

# Multiple device support
All commands that interact with a Nest can either use the default device from your config file or you can specify a device ID
e.g.
```
./nerdnest temp 70
./nerdnest temp 75 ku34h5kjefhkdsjfhsdf
./nerdnest status
./nerdnset status ku34h5kjefhkdsjfhsdf
```
To set the default device, even you only have one, run the setdefault command:
```
./nerdnest setdefault
Nest: raY9xdYrYY036u2rgaeP_lJ
Kitchen: ku34h5kjefhkdsjfhsdf
Enter default device ID from above list: ku34h5kjefhkdsjfhsdf
```
# Configuration keys
accesstoken = "ACCESSTOKEN"

mythermostat = "MYDEVICEID"

units = "[cCfF]"

# Choosing units for temperature
By default nerdnest uses Farenheit for temperature both to display the status and when setting temperature. You can
override this behavior by adding a configuration key called units and setting it to either "c" or "C".

For Farenheit you must specify the temperature in whole numbers e.g. 70, 75
For Celcius you can specify half units as well e.g. 19, 20.5, 23.5
