package main

import "syscall"

// go build -ldflags "-linkmode external -extldflags -static" ./socket.go
func main() {
	// socket system call returns a file descriptor
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}

	// print the file descriptor, should be 3
	println(sockfd)
}
