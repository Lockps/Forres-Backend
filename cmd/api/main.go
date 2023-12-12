package main

import (
	"log"
	"net/http"
)

type application struct {
	Domain string
	port   string
}

var key = make([]byte, 32)

func main() {
	// Encoder(0, key)
	// DeCoder(0, key)

	// EncoderBase64(0)
	// DecoderBase64(0)

	var app application
	// app.Test()

	app.port = ":8080"
	log.Println("Starting application on port :", app.port)

	err := http.ListenAndServe(app.port, app.routes())

	if err != nil {
		log.Fatal(err.Error())
	}
	// x, _ := database.ReadFirstFieldFromUsersDB()
	// fmt.Println(x)

}
