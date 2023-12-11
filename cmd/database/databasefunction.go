package database

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Lockps/Forres-release-version/cmd/function"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func FetchPost(r *http.Request, permission int) string {
	if r.Method != http.MethodPost {
		return "Method Not Allowed"
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "Can't Read Data : " + err.Error()
	}
	defer r.Body.Close()

	dbname := GetLocation(permission)

	file, err := os.OpenFile(dbname+".db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "Can't Open DBMS System : " + err.Error()
	}
	defer file.Close()

	fileinfo, _ := file.Stat()
	if fileinfo.Size() != 0 {
		_, err = file.WriteString("\n1  ")
		if err != nil {
			return "Can't Connect with Database"
		}
	}
	_, err = file.Write(body)
	if err != nil {
		return "Can't Store Data,please Try Again"
	}

	return "Store Data Successful"
}

func FetchGet(w http.ResponseWriter, r *http.Request, permission, coll int) {
	filePath := GetLocation(permission)
	fmt.Println(filePath)

	file, err := os.Open("./database" + filePath + ".db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var matchingLines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, fmt.Sprintf("%v", coll)) {
			matchingLines = append(matchingLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(matchingLines)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//? ==================== LOG IN =============================

func ReadFirstFieldFromUsersDB(permission, field int) ([]string, error) {
	dbname := GetLocation(permission) + ".db"

	file, err := os.Open(dbname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var firstFields []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 {
			firstFields = append(firstFields, fields[field])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return firstFields, nil
}

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(function.StrToByteSlice(("Method Not Allowed")))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(function.StrToByteSlice("Can't Read Data : " + err.Error()))
		return
	}
	defer r.Body.Close()

	dbname := GetLocation(0)
	w.Write([]byte(dbname))

	filepath := dbname + ".db"
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	filedata, err := file.Stat()
	if err != nil {
		w.Write(function.StrToByteSlice("Permission denied!"))
		return
	}

	if filedata.Size() != 0 {
		_, err = file.WriteString("\n" + uuid.NewString() + " ")
		if err != nil {
			w.Write(function.StrToByteSlice("Can't connect to the database"))
			return
		}
	} else {
		_, err := file.WriteString(uuid.NewString() + " ")
		if err != nil {
			w.Write(function.StrToByteSlice("Can't Connect to Database"))
		}
	}

	_, err = file.Write(body)
	if err != nil {
		w.Write(function.StrToByteSlice("Can't Store your information,please try again "))
	}

	return
}

func ValidateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't Read Data", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	valid, _, err := ValidateUser(string(body))
	if err != nil {
		http.Error(w, "Error validating user", http.StatusInternalServerError)
		return
	}

	if valid {
		fmt.Fprintln(w, "User is valid")
	} else {
		fmt.Fprintln(w, "Invalid credentials")
	}
}

func ValidateUser(dataFromFrontend string) (bool, string, error) {
	fields := strings.Fields(dataFromFrontend)
	fmt.Println(len(fields))
	if len(fields) != 2 {
		return false, "", fmt.Errorf("invalid data format")
	}

	username := fields[0]
	password := fields[1]
	// permission := fields[3]
	fmt.Println(username)
	fmt.Println(password)

	file, err := os.Open("Users.db")
	if err != nil {
		return false, "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dbFields := strings.Fields(line)
		if len(dbFields) >= 4 && dbFields[1] == username && dbFields[2] == password {
			fmt.Println(dbFields)
			return true, dbFields[0], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, "", err
	}
	fmt.Print("UnExepected DBField")

	return false, "", nil
}

func GetUserName(w http.ResponseWriter, r *http.Request) []string {
	name, err := ReadFirstFieldFromUsersDB(0, 1)
	if err != nil {
		w.Write(function.StrToByteSlice(err.Error()))
	}

	return name
}

//! ==========================  TOKEN  ===============================

type APIError struct {
	Error string
}

func WirteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling jwt Auth!")

		tokenString := r.Header.Get("x-jwt-token")

		_, err := validateJWT(tokenString)
		if err != nil {
			WirteJson(w, http.StatusForbidden, APIError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := "hunter0123"
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

//** ======================= BOOKING ====================================

func AddBookingToDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(function.StrToByteSlice(("Method Not Allowed")))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(function.StrToByteSlice("Can't Read Data : " + err.Error()))
		return
	}
	defer r.Body.Close()

	dbname := GetLocation(2)

	file, err := os.OpenFile(dbname+".db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		w.Write(function.StrToByteSlice("Can't Open DBMS System : " + err.Error()))
		return
	}
	defer file.Close()

	fileinfo, _ := file.Stat()
	if fileinfo.Size() != 0 {
		_, err = file.WriteString("\n")
		if err != nil {
			w.Write(function.StrToByteSlice("Can't Connect with Database"))
			return
		}
	}
	_, err = file.Write(body)
	if err != nil {
		w.Write(function.StrToByteSlice("Can't Store Data,please Try Again"))
		return
	}

	w.Write(function.StrToByteSlice("Store Data Successful"))
}

func GetUnAvaliableSeat(w http.ResponseWriter, r *http.Request) {
	dbName := GetLocation(2)
	file, err := os.Open(dbName + ".db")
	if err != nil {
		w.Write(function.StrToByteSlice("Can't Open The Database"))
		return
	}
	defer file.Close()

	var firstFields []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 {
			firstFields = append(firstFields, fields[0]) // Append the first field from each line
		}
	}

	if err := scanner.Err(); err != nil {
		w.Write(function.StrToByteSlice("Error reading the file"))
		return
	}

	// Join the first fields obtained from each line and send as output
	output := strings.Join(firstFields, ",")
	op, err := json.Marshal(output)
	if err != nil {
		w.Write([]byte("Error to convert to json!"))
		return
	}
	w.Write(op)

}

//* ======================= DELETE , UPDATE ===========================

func DeleteData(w http.ResponseWriter, r *http.Request) {

}

func ReadDataHandler(w http.ResponseWriter, r *http.Request) {
	x, _ := ReadAllData(0, 1, 2)
	w.Write(x)
}

func ReadAllData(permission, nameField, passwordField int) ([]byte, error) {
	dbname := GetLocation(permission) + ".db"

	file, err := os.Open(dbname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []map[string]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > passwordField {
			user := make(map[string]string)
			user["name"] = fields[nameField]
			user["password"] = fields[passwordField]

			users = append(users, user)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func ReadAllline(permission int) ([]byte, error) {
	dbname := GetLocation(permission) + ".db"

	file, err := os.Open(dbname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []map[string]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		user := make(map[string]string)
		for i, field := range fields {
			user[fmt.Sprintf("field_%d", i)] = field
		}
		users = append(users, user)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func DeleteLinesContainingValue(filepath, value string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)
		if len(fields) > 1 && fields[1] != value {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func UpdateFieldByCondition(permission int, field2Value string, fieldToUpdate int, newData string) error {

	filepath := GetLocation(permission) + ".db"

	fmt.Println(filepath)
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) > 2 && fields[2] == field2Value {
			if len(fields) > fieldToUpdate {
				fields[fieldToUpdate] = newData
			}
			line = strings.Join(fields, " ")
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

type Booking struct {
	Table  int    `json:"table"`
	Name   string `json:"name"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Car    string `json:"car"`
	People string `json:"people"`
	Course string `json:"course"`
}

func ReadFieldsFromDB(permission, field int) ([]string, error) {
	dbname := GetLocation(permission) + ".db"

	file, err := os.Open(dbname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fields []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		var booking Booking
		if err := json.Unmarshal(line, &booking); err != nil {
			return nil, err
		}

		switch field {
		case 1:
			fields = append(fields, fmt.Sprintf("%v", booking.Table))
		case 2:
			fields = append(fields, booking.Name)
		case 3:
			fields = append(fields, booking.Date)
		case 4:
			fields = append(fields, booking.Time)
		case 5:
			fields = append(fields, booking.Car)
		case 6:
			fields = append(fields, booking.People)
		case 7:
			fields = append(fields, booking.Course)
		default:
			return nil, errors.New("invalid field index")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

func ReadUserTable(w http.ResponseWriter, r *http.Request) {
	data, err := ReadAllline(0)
	if err != nil {
		return
	}
	w.Write(data)
}

func ReadCustomerTable(w http.ResponseWriter, r *http.Request) {
	data, err := ReadAllline(1)
	if err != nil {
		return
	}
	w.Write(data)
}

func ReadStaffTable(w http.ResponseWriter, r *http.Request) {
	data, err := ReadAllline(2)
	if err != nil {
		return
	}
	w.Write(data)
}

func ReadAndReturnString(permission int) (string, error) {
	filepath := GetLocation(permission) + ".db"
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("Error opening the file: %s", err)
	}
	defer file.Close()

	var result string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			log.Printf("Error decoding JSON: %s", err)
			continue
		}

		// Convert the data map to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error encoding JSON: %s", err)
			continue
		}

		result += string(jsonData) + "\n"
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Error reading the file: %s", err)
	}

	return result, nil
}
