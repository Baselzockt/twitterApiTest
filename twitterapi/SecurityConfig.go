package twitterapi

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

type SecurityConfig struct {
	ApiSecret    string
	ApiKey       string
	AccessToken  string
	AccessSecret string
}

func LoadSecurityConfig() *SecurityConfig {
	log.Debug("Loading SecurityConfig from conf.json")
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	securityConfig := &SecurityConfig{}
	err := decoder.Decode(&securityConfig)

	if err != nil {
		log.Error("Could not load config creating file and closing program")
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		err = encoder.Encode(securityConfig)

		if err != nil {
			log.Error("Could not encode empty security config")
		}

		err = os.WriteFile("conf.json", buffer.Bytes(), os.ModeExclusive)

		if err != nil {
			log.Error("Could not save default(empty) config")
		}

		log.Fatal("Config file could not be parsed")
	}
	err = file.Close()

	if err != nil {
		log.Fatal("Could not close config file")
	}

	return securityConfig
}
