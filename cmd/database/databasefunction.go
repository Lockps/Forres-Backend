package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

func ReadFirstFieldFromUsersDB() ([]string, error) {
	dbname := "Users.db"

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
			firstFields = append(firstFields, fields[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return firstFields, nil
}

func UserExists(datacome string) (bool, error) {
	samp, err := ReadFirstFieldFromUsersDB()
	exist := false

	if err != nil {
		return false, err
	}
	for i := 0; i < len(samp); i++ {
		if datacome == samp[i] {
			return true, nil
		}
	}
	return exist, nil
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

	userData := string(body)

	exists, err := UserExists(userData)
	if err != nil {
		w.Write([]byte("Error checking user existence: " + err.Error()))
		return
	}

	if exists {
		w.Write([]byte("User already exists in the database"))
		return
	}

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
		_, err = file.WriteString("\n1  " + uuid.NewString() + " ")
		if err != nil {
			w.Write(function.StrToByteSlice("Can't connect to the database"))
			return
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

	valid, err := validateUser(string(body))
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

func validateUser(dataFromFrontend string) (bool, error) {
	fields := strings.Fields(dataFromFrontend)
	if len(fields) != 4 {
		return false, fmt.Errorf("invalid data format")
	}

	username := fields[1]
	password := fields[2]
	// permission := fields[3]

	file, err := os.Open("Users.db")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dbFields := strings.Fields(line)
		if len(dbFields) >= 4 && dbFields[1] == username && dbFields[2] == password {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

//!==========================  TOKEN  ===============================

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
		_, err = file.WriteString("\n1  ")
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
