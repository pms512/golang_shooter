package main

import (
	"os"
	"fmt"
	"encoding/json"
	"database/sql"
	_ "github.com/alexbrainman/odbc"
	"strings"
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

func checkStmtAndParams(tx map[string]interface{}) bool {
	paramCountInTx := 0
	paramCountInStmt := 0
    for key, value := range tx {
        if strings.HasPrefix(key,"param") {
            paramCountInTx++
        } else if key == "stmt" {
            convertedStmt := value.(string)
            paramCountInStmt = strings.Count(convertedStmt, "?")
        }
    }
	if paramCountInTx == paramCountInStmt {
		return true
	} else {
		return false
	}
}

func main() {
	var config map[string]interface{}
//	var version string
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
/*
	err = db.QueryRow("").Scan(&version)

	if err != nil {
		fmt.Println("Error querying:", err)
		panic(err)
	}
	fmt.Println("version:", version)
*/

	//extract statement and parameter from config
    transactions := config["transactions"].(map[string]interface{})

	for key, tx := range transactions {
		convertedTx := tx.(map[string]interface{})
		fmt.Println("txName:", key)
		//check statement and get the number of bind parameters
		isValid := checkStmtAndParams(convertedTx)
		if isValid != true {
			fmt.Println("The number of params and the number of bind parameters in statement is different. Please edit config.conf correctly. Exiting..")
			return
		}
	}

	fmt.Println("All processes has been completed. Exiting..")
}
