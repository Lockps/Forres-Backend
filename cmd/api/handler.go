package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Person struct {
	fname string
	lname string
}

type Data struct {
	Items []Person `json:"items"`
}

type datacheck struct {
	Permission int    `json:"Permission"`
	Status     string `json:"status"`
	Version    string `json:"version"`
}

func (app *application) FirstPage(w http.ResponseWriter, r *http.Request) {
	payload := *&datacheck{
		Permission: 1,
		Status:     "Active",
		Version:    "1.0.0",
	}

	out, err := json.Marshal(payload)
	if err != nil {
		buf := bytes.NewBufferString(err.Error())
		fmt.Fprint(buf)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) postDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	file, err := os.OpenFile("file.db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	if fileInfo.Size() != 0 {
		_, err = file.WriteString("\n")
		if err != nil {
			http.Error(w, "Error writing to file", http.StatusInternalServerError)
			return
		}
	}

	_, err = file.Write(body)
	if err != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Data appended successfully to file.txt\n")
}

func (app *application) ReadFile(w http.ResponseWriter, r *http.Request) {
	filePath := "./file.db"

	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	lines := strings.Split(string(data), "\n")

	var people []Person
	for _, line := range lines {
		names := strings.Split(line, " ") // Assuming data is space-separated fname lname
		if len(names) == 2 {
			person := Person{
				fname: names[0],
				lname: names[1],
			}
			people = append(people, person)
		}
	}

	result := Data{
		Items: people,
	}

	out, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) Test(w http.ResponseWriter, r *http.Request) {

}
