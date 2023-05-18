.PHONY: help

goserver:
	go run ./passwords.go
race:
	go test -race -v -run='^\QTest_Race_' .
pyserver:
	python ./passwords.py
ab:
	ab -n 500000 -c 1000 http://127.0.0.1:8080/password

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "This makefile compares the differences between dynamically linking glibc and statically linking glibc."
	@echo "It also provides test methods and proofs that the net package is linked to glibc."
	@echo ""
	@echo "Available targets:"
	@echo "  make goserver      Run golang server"
	@echo "  make race          Run golang race detector"
	@echo "  make pyserver      Run python server"
	@echo "  make ab            Run website benchmark"
	@echo ""