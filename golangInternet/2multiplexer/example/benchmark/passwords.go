package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// define charset
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#@$(in Traditional Chinese)"

// generate password function
func generatePassword(length int) string {
	// use current time as random seed
	rand.NewSource(time.Now().UnixNano())

	// allocate memory for password
	password := make([]byte, length)

	// generate random password
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	// return password
	return string(password)
}

// http handler function
func handler(w http.ResponseWriter, r *http.Request) {
	// channel to collect passwords
	passwords := make(chan string, 100)

	// generate 100 passwords concurrently
	for i := 0; i < 100; i++ {
		go func() {
			passwords <- generatePassword(12)
		}()
	}

	// collect the passwords
	result := make([]string, 100)
	for i := 0; i < 100; i++ {
		result[i] = <-passwords
	}

	// write the passwords to the response
	fmt.Fprint(w, result)
}

// main function
func main() {
	// set router
	http.HandleFunc("/", handler)
	// start server
	http.ListenAndServe(":8080", nil)
}
