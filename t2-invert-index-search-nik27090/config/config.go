package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerAddress string
	Direct        string
	Addr          string
	Username      string
	Pass          string
	DB            string
}

func Load() Config {
	conf, _ := os.Open("config.json")
	defer conf.Close()
	decoder := json.NewDecoder(conf)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}
