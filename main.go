package main

import "net/http"

func main() {
	// HandleFunc registers a handler function for the "/hello" endpoint.
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		// Write a response back to the client with the message "hello, this is server".
		_, _ = writer.Write([]byte("hello, this is server"))
	})

	// ListenAndServe listens on the TCP network address 127.0.0.1:1234 for incoming connections.
	// It serves incoming HTTP requests using the registered handler, which in this case is the "/hello" endpoint.
	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
