package main

import (
	"os"
	"fmt"
	"encoding/json"
	"database/sql"
	_ "github.com/alexbrainman/odbc"
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
	var version string
	fmt.Println("=======================")
	fmt.Println("  Goldilocks Shooter   ")
	fmt.Println("=======================")
	fmt.Println()
	config = ReadConfigFile()
	if config != nil {
		fmt.Println("ReadConfigFile() Done!")
	}
	connString := fmt.Sprintf("DSN=%s",config["dsn"])
	db, err := sql.Open("odbc", connString)
	if err != nil {
		fmt.Println("Error connecting to DB:", err)
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connection success!")

	err = db.QueryRow("SELECT version FROM X$INSTANCE").Scan(&version)

	if err != nil {
		fmt.Println("Error querying:", err)
		panic(err)
	}
	fmt.Println("version:", version)

}
