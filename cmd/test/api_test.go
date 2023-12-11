package test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func SignUpTest() {
	baseURL := "http://your-api-url/signup" // Replace this with your actual API URL

	for i := 1; i <= 30; i++ {
		// Generate sample user data (replace this with your actual data)
		username := fmt.Sprintf("user%d", i)
		email := fmt.Sprintf("user%d@example.com", i)
		password := "password123" // Replace this with a secure password generation logic

		// Create the request body
		requestBody := []byte(fmt.Sprintf("username=%s&email=%s&password=%s", username, email, password))

		// Send POST request to /signup
		resp, err := http.Post(baseURL, "application/x-www-form-urlencoded", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Error creating user %d: %s\n", i, err.Error())
			continue
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error creating user %d. Status code: %d\n", i, resp.StatusCode)
			continue
		}

		// Successful signup
		log.Printf("User %d signed up successfully\n", i)
	}
}
