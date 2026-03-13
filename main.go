package main

import (
	"os"
	"fmt"
	"encoding/json"
)

func ReadConfigFile() map[string]interface{} {
	f, err := os.ReadFile("config.conf")

	if err != nil {
		fmt.Println("Error reading config.conf:", err)
		fmt.Println("Please create or edit config.conf properly.")
		fmt.Println("Exiting...")
		return nil
	}
	fmt.Println("config.conf exists! Converting json to map..")
	var result map[string]interface{}
	err = json.Unmarshal(f, &result)
	if err != nil {
		fmt.Println("Error converting json to map:", err)
	}
	return result
}

func main() {
	var config map[string]interface{}

	fmt.Println("=======================")
	fmt.Println("  Goldilocks Shooter   ")
	fmt.Println("=======================")
	fmt.Println()
	config = ReadConfigFile()
	if config != nil {
		fmt.Println("ReadConfigFile() Done!")
	}
}
