package utils

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

// Configurations json
type Configurations struct {
	ServerURL    string `json:"serverUrl"`
	ServerPort   int    `json:"serverPort"`
	UploadFolder string `json:"uploadFolder"`
	TempFolder   string `json:"tempFolder"`
}

// GetConfigs read and parse configurations json file
func GetConfigs() (Configurations, error) {
	// read file
	fdata, err := ioutil.ReadFile(path.Join(".", "config.json"))
	if err != nil {
		return Configurations{}, err
	}
	// json data
	var config Configurations
	// unmarshall it
	err = json.Unmarshal(fdata, &config)
	if err != nil {
		return Configurations{}, err
	}
	return config, nil
}
