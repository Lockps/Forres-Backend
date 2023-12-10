package database

import (
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	FetchGet(w, r, 1, 1)
}
