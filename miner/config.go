package miner

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"runtime"
	"strings"
)

// Data is the structure that stores all arguments passed to the miner.
type Data struct {
	Crt        string
	Key        string
	Endpoint   string
	UserConfig UserConfig
	Workers    int
}

// TODO: Add in a UserConfig model folder.
type UserConfig struct {
	Name      string   `json:"Name"`
	Mails     []string `json:"Mails"`
	Skype     string   `json:"Skype"`
	BirthDate string   `json:"BirthDate"`
	Country   string   `json:"Country"`
	Addess    []string `json:"Addess"`
}

// GetUserConfigurationFile parshe json with user contact information file.
// TODO: Add Error in case it fails to process JSON.
func getUserConfigurationFile(configFile string) UserConfig {
	userConfig := UserConfig{}
	file, err := os.Open(configFile)
	if err != nil {
		log.Println("error:", err)
	} else {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&userConfig)
		if err != nil {
			log.Println("error:", err)
		}
	}
	return userConfig
}

// Get will parshe the args and return the configuration in the Data structure.
func Get() Data {
	var config Data

	flag.StringVar(&config.Endpoint, "connect", "localhost:4433", "who to connect to")
	flag.StringVar(&config.Crt, "crt", "./config/certs/public.crt", "certificate")
	flag.StringVar(&config.Key, "key", "./config/certs/private.key", "key")
	flag.IntVar(&config.Workers, "workers", runtime.NumCPU(), "number of workers to run in the pool")

	userConfigFilePath := flag.String("userConfigFile", "./config/config.json", "JSON config file to read.")
	flag.Parse()
	config.UserConfig = getUserConfigurationFile(*userConfigFilePath)

	if !strings.Contains(config.Endpoint, ":") {
		config.Endpoint += ":443"
	}

	return config
}
