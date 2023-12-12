package test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func SignUpTest() {
	baseURL := "http://localhost/8080/createuser"

	for i := 1; i <= 30; i++ {
		username := fmt.Sprintf("user%d", i)
		email := fmt.Sprintf("user%d@example.com", i)
		password := "password123"

		requestBody := []byte(fmt.Sprintf("username=%s&email=%s&password=%s", username, email, password))

		resp, err := http.Post(baseURL, "application/x-www-form-urlencoded", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Error creating user %d: %s\n", i, err.Error())
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error creating user %d. Status code: %d\n", i, resp.StatusCode)
			continue
		}

		log.Printf("User %d signed up successfully\n", i)
	}
}
