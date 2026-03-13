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
//읽어오고 1차 validation을 하는 게 맞을까? (dsn, threads 등의 key가 있어야 한다 같은.)
	fmt.Println("[Reading file config.conf..]")
	f, err := os.ReadFile("config.conf")

	if err != nil {
		fmt.Println("Error reading config.conf:", err)
		fmt.Println("Please create or edit config.conf correctly.")
		fmt.Println("Exiting...")
		return nil
	}
	fmt.Println("config.conf exists! Converting json to map..")
	var result map[string]interface{}
	err = json.Unmarshal(f, &result)
	if err != nil {
		fmt.Println("Error converting json to map:", err)
	}
	fmt.Println("ReadConfigFile() Done!")
	return result
}

func CheckStmtAndParams(tx map[string]interface{}) bool {
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

func ConnectToDB(config map[string]interface{}) (*sql.DB, error) {
	connString := fmt.Sprintf("DSN=%s",config["dsn"])
	db, err := sql.Open("odbc", connString)
	if err != nil {
		fmt.Println("Error connecting to DB:", err)
		panic(err)
	}
	return db, err
}

type stmtType int
const (
	stmtINSERT stmtType = iota
	stmtSELECT
	stmtUPDATE
	stmtDELETE
)
func GetStmtType(stmt string) int {
    extractedType := stmt[0:6]
	switch extractedType {
	case "INSERT":
		return 0
	case "SELECT":
		return 1
	case "UPDATE":
		return 2
	case "DELETE":
		return 3
	default:
		return -1
	}
}

func main() {
	var config map[string]interface{}
	var txSlice []map[string]interface{}
	//intro
	fmt.Println("=======================")
	fmt.Println("  Goldilocks Shooter   ")
	fmt.Println("=======================")
	fmt.Println()

	//read config file
	config = ReadConfigFile()

	//connect to DB
	db, _ := ConnectToDB(config)
	defer db.Close()

	fmt.Println("Connection success!")

	//validate transactions in config
    transactions := config["transactions"].(map[string]interface{})
	for key, rawTx := range transactions {
		convertedTx := rawTx.(map[string]interface{})
		//check statement and get the number of bind parameters
		isValid := CheckStmtAndParams(convertedTx)
		if isValid != true {
			fmt.Printf("[%s] The number of params and the number of bind parameters in statement is different.\n", key)
			fmt.Printf("Please edit content of %s in config.conf correctly.\n", key)
			fmt.Println("Exiting..")
			return
		}
		txSlice = append(txSlice, convertedTx)
	}

	//execute transactions
	// 나중에는 여러 thread를 설정했을 떄 어떻게 할지 생각하기.
	/* 
	  각 transaction을 수행할 때마다 아래 절차로 하자.
	  - statement 추출(preparedStatement 형태로)
	  - bind parameter 추출
	   - CONST이면 해당 literal 1개로만
	   - SERIAL이면 초기값, 공차(d)
	   - RANDOM이면 ??

	*/
	fmt.Println()
	fmt.Println("Executing transactions..")
	for i := 0; i < len(txSlice); i++ {
		stmt := txSlice[i]["stmt"].(string)
		typeOfStmt := GetStmtType(stmt)
		fmt.Println("stmt:", stmt)
		fmt.Println("stmtINSERT:", stmtINSERT)
		fmt.Println("typeOfStmt:", typeOfStmt)

		if stmtType(typeOfStmt) == stmtINSERT {
			fmt.Println("This statement is INSERT.")
		} else if stmtType(typeOfStmt) == stmtSELECT {
			fmt.Println("This statement is SELECT.")
		} else if stmtType(typeOfStmt) == stmtUPDATE {
			fmt.Println("This statement is UPDATE.")
		} else if stmtType(typeOfStmt) == stmtDELETE {
			fmt.Println("This statement is DELETE.")
		} else {
			fmt.Println("This statement is invalid.")
		}
	}

	fmt.Println("All processes has been completed. Exiting..")
}
