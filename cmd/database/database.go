package database

import (
	"encoding/json"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	x := GetUserName(w, r)
	json.Marshal(x)
	w.Write([]byte(x[2]))
}
