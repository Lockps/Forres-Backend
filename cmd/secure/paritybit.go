package secure

import (
	"log"
	"net/http"

	"github.com/Lockps/Forres-release-version/cmd/database"
	"github.com/go-chi/chi/v5"
)

type SecureData struct {
	Domain string
	Port   string
}

var key int

func init() {
	key = 12323
}

func Secure() {
	var sec SecureData

	sec.Port = ":1111"
	log.Println("Starting parity application in : ", sec.Port)

	err := http.ListenAndServe(sec.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux := chi.NewRouter()
		mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

		})
	}))
	if err != nil {
		log.Fatal(err)
	}
}

func Getparity() string {
	data, _ := database.ReadFirstFieldFromUsersDB(4, 0)
	return (data[0])
}
