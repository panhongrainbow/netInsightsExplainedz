.PHONY: help

include:
	go build -ldflags "-linkmode external -extldflags -static" ./socket.go
	@ls -lh | grep socket

normal:
	go build ./socket.go
	@ls -lh | grep socket

strace:
	go build -ldflags "-linkmode external -extldflags -static" ./socket.go
	strace ./socket

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "This makefile compares the differences between dynamically linking glibc and statically linking glibc."
	@echo "It also provides test methods and proofs that the net package is linked to glibc."
	@echo ""
	@echo "Available targets:"
	@echo "  include  - statically linking glibc"
	@echo "  normal   - dynamically linking glibc"
	@echo "  strace   - proof that the net package is linked to glibc"
	@echo ""