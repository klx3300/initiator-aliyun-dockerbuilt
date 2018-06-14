package configrd

import (
	"encoding/json"
	"logger"
	"os"
)

// Config represents a config file.
type Config string

// This is the string format of config file name.
func (c Config) This() string {
	return string(c)
}

// ReadConfig read config content from config file.
func (c Config) ReadConfig() map[string]string {
	file, err := os.Open(c.This())
	if err != nil {
		panic("Unable to read config file" + c.This() + err.Error())
	}
	defer file.Close()
	jdecoder := json.NewDecoder(file)
	cfg := make(map[string]string)
	err = jdecoder.Decode(&cfg)
	if err != nil {
		panic("Unable to unmarshal json while reading config" + c.This() + err.Error())
	}
	return cfg
}

// WriteConfig write config out to config file.
func (c Config) WriteConfig(cfg map[string]string) {
	file, err := os.OpenFile(c.This(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to write out to ", c, err)
		return
	}
	defer file.Close()
	jencoder := json.NewEncoder(file)
	jencoder.Encode(cfg)
}
