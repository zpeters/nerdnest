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

# Basic Usage
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
