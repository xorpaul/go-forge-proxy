package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// readConfigfile creates the ConfigSettings struct from the config file
func readConfigfile(configFile string) ConfigSettings {
	Debugf("Trying to read config file: " + configFile)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		Fatalf("readConfigfile(): There was an error parsing the config file " + configFile + ": " + err.Error())
	}

	var config ConfigSettings
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		Fatalf("YAML unmarshal error: " + err.Error())
	}

	// check if cachedir exists
	config.CacheDir = checkDirAndCreate(config.CacheDir, "cachedir")

	if len(config.ForgeUrl) == 0 {
		config.ForgeUrl = "https://forgeapi.puppetlabs.com"
	}

	// set default timeout to 5 seconds if no timeout setting found
	if config.Timeout == 0 {
		config.Timeout = 5
	}

	// set default max Go routines for Forge and Git module resolution if none is given
	if maxworker == 0 {
		config.Maxworker = 5
	} else {
		config.Maxworker = maxworker
	}

	return config
}
