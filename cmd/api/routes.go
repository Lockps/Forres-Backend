package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Lockps/Forres-release-version/cmd/database"
	"github.com/Lockps/Forres-release-version/cmd/function"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": ""})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(app.enableCORS)

	r.Post("/signup", database.CreateUsers)

	r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Can't read Data", http.StatusInternalServerError)
		}

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		Valid, uuid, err := database.ValidateUser(string(body))

		if Valid {
			_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": uuid})
			if err != nil {
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
				return
			}
			fmt.Print(tokenString)

			w.Header().Set("Content-Type", "application/json")
			w.Header().Add("Authorization", tokenString)
			// fmt.Fprintf(w, `{"token": "%s"}`, tokenString)
			w.Write(function.StrToByteSlice("valid"))
			return
		} else {
			w.Write(function.StrToByteSlice("invalid"))
		}

		http.Error(w, "Invalid request body", http.StatusUnauthorized)
	})

	r.Post("/Booking/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		username := chi.URLParam(r, "username")

		// Retrieve the token from the Authorization header
		// tokenString := r.Header.Get("Authorization")
		// if tokenString == "" {
		// 	http.Error(w, "Unauthorized: Token not provided", http.StatusUnauthorized)
		// 	return
		// }

		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized: Invalid user", http.StatusUnauthorized)
			return
		}

		if userID != username {
			http.Error(w, "Unauthorized: User does not match token", http.StatusUnauthorized)
			return
		}

		database.AddBookingToDB(w, r)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Booking successful"))
	})

	r.Get("/getusername", database.ReadDataHandler)
	r.Get("/getuserstable", database.ReadUserTable)
	r.Get("/getcustomertable", database.ReadCustomerTable)
	r.Get("/getstafftable", func(w http.ResponseWriter, r *http.Request) {
		data, err := database.ReadAndReturnString(2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(data))
	})

	r.Delete("/delete/{db}/{name}", func(w http.ResponseWriter, r *http.Request) {
		dbStr := chi.URLParam(r, "db")
		name := chi.URLParam(r, "name")

		db, err := strconv.Atoi(dbStr)
		if err != nil {
			http.Error(w, "Invalid db parameter: not an integer", http.StatusBadRequest)
			return
		}

		fmt.Println(name)
		database.DeleteLinesContainingValue(db, name)

		w.Write(function.StrToByteSlice("Delete Successful!"))
	})

	r.Get("/update/{db}/{name}/{field}/{data}", func(w http.ResponseWriter, r *http.Request) {
		dbStr := chi.URLParam(r, "db")
		name := chi.URLParam(r, "name")
		field := chi.URLParam(r, "field")
		data := chi.URLParam(r, "data")

		db, err := strconv.Atoi(dbStr)
		if err != nil {
			http.Error(w, "Invalid db parameter: not an integer", http.StatusBadRequest)
			return
		}
		fieldchange, err := strconv.Atoi(field)
		if err != nil {
			http.Error(w, "Invalid db parameter: not an integer", http.StatusBadRequest)
			return
		}

		fmt.Println(db)
		fmt.Println(name)
		fmt.Println(fieldchange)
		fmt.Println(data)

		database.UpdateFieldByCondition(db, name, fieldchange, data)

	})

	r.Post("/booking", database.AddBookingToDB)

	r.Get("/gettable", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}

		table, err := database.ReadFieldsFromDB(2, 1)
		if err != nil {
			http.Error(w, "Test", http.StatusInternalServerError)
			return
		}

		outputdata, _ := json.Marshal(table)
		w.Header().Set("Content-Type", "application/json")
		w.Write(outputdata)
	})

	r.Get("/getbalance/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		balance, _ := database.GetBalanceByValueFromFile(0, name)

		w.Write([]byte(balance))
		fmt.Print(name)
	})

	r.Post("/balance", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		data, _ := database.ReadAllline(0)
		jsonStr := string(data)

		var op []map[string]string
		if err := json.Unmarshal([]byte(jsonStr), &op); err != nil {
			return
		}

		for _, obj := range op {
			if val, ok := obj["field_1"]; ok && val == string(body) {
				if valurfield7, ok := obj["field_5"]; ok {
					w.Write([]byte(valurfield7))
				}
			}
		}

	})

	r.Post("/deleteRecord", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(string(body))
		database.DeleteLineByTable(2, string(body))

		w.Write(function.StrToByteSlice("Delete Success"))
	})

	r.Post("/topup/{cost}", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		cost := chi.URLParam(r, "cost")

		database.UpdateFieldByCondition(0, string(body), 6, cost)

	})

	r.Get("/staff/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		role, _ := database.GetRoleByValueFromFile(0, name)

		w.Write([]byte(role))
	})

	return r
}
