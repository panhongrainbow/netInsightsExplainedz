package main

import "testing"

func Test_Race_password(t *testing.T) {
	// channel to collect passwords
	passwords := make(chan string, 100) // ----- race ---->

	// generate 100 passwords concurrently
	for i := 0; i < 100; i++ {
		go func() {
			passwords <- generatePassword(12) // <----- race -----
		}()
	}

	// collect the passwords
	result := make([]string, 100)
	for i := 0; i < 100; i++ {
		result[i] = <-passwords
	}
}
