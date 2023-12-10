package database

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Lockps/Forres-release-version/cmd/function"
)

func BasicTestPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method Not Allowed")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Can't Read Data : " + err.Error())
	}
	defer r.Body.Close()

	dbname := "Customer"

	file, err := os.OpenFile(dbname+".db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Can't Open DBMS System : " + err.Error())
	}
	defer file.Close()

	fileinfo, _ := file.Stat()
	if fileinfo.Size() != 0 {
		_, err = file.WriteString("\n1  ")
		if err != nil {
			fmt.Println("Can't Connect with Database")
		}
	}
	_, err = file.Write(body)
	if err != nil {
		fmt.Println("Can't Store Data,please Try Again")
	}

	w.Write(function.StrToByteSlice("Store Data in " + dbname + " Successful"))
}
