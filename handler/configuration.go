package handler

import (
	"os"
	"encoding/json"
	"fmt"
)

type Configuration struct {
	Authentication             bool      `json:"authentication"`
	AuthenticationParameter    string    `json:"authenticationParameter"`
	AuthenticationToken        string    `json:"authenticationToken"`
}

func (c *Configuration) ParseConfig() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	*c = Configuration{}
	err := decoder.Decode(&c)
	if err != nil {
		fmt.Println("error:", err)
	}
}