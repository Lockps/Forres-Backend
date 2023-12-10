package main

import (
	"log"
	"net/http"
)

type application struct {
	Domain string
	port   string
}

func main() {
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
