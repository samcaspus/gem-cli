package iooperations

import (
	"fmt"

	"github.com/samcaspus/gem-cli/constants"
	"github.com/samcaspus/gem-cli/utils"
)

func TakeApiKeyInput() {
	// take input from user for api key and write it to confir_file_path
	fmt.Println("Config file does not exist. Please enter your API key:")
	var apiKey string
	fmt.Scanln(&apiKey)
	utils.WriteToFile(constants.CONFIG_FILE_PATH, apiKey)
}

func GetApiKey() string {
	return utils.ReadFile(constants.CONFIG_FILE_PATH)
}
