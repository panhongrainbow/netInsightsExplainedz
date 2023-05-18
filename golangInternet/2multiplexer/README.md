# golang multiplexer

## Introduction

`multiplexer` is divided into two processes: `waiting for events to occur` and `handling events`.
Specifically:

1. `Waiting for events to occur`
   The program will monitor multiple events (usually IO events) and wait for any of them to occur.
2. `Handling events`
   When an event occurs, the program will suspend waiting for other events and handle the event.
3. `After handling, continue to wait for other events.` 
4. By looping through waiting and handling, one process can monitor and `respond to multiple events` at the same time.
   This is `the basic model of multiplexing`.

Example:

````go
package main

import (
    "fmt"
    "net/http"
)

func handleConn(w http.ResponseWriter, r *http.Request) {
    _, err := w.Write([]byte("Hello"))
    if err != nil {
        fmt.Println(err)
    } 
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handleConn)

    server := &http.Server{
        Addr:    ":8000",
        Handler: mux,
    }

    err := server.ListenAndServe()
    if err != nil {
        fmt.Println(err)
    }
}
````

AB test

```bash
sudo apt-get install apache2-utils

ab -n 1000 -c 100 http://localhost:8000/
```

Go routine leak

In the implementation of `golang's multiplexer`, the problem of goroutine leaks must be considered.

The main reasons are:

1. Multiple goroutines operate `concurrently`. If a certain goroutine `exits abnormally without releasing resources`, it will cause leaks.
2. When `a network connection is disconnected`, the corresponding goroutine does not exit and `still holds the resources of the connection`, which causes leaks. 
3. Goroutines `forget or incorrectly release resources when exiting`, which can also cause leaks.

Solutions

1. Use `context to manage the lifecycle of goroutines`.
   When the parent goroutine exits, call the cancel() method of the context.
   Child goroutines can receive the cancel signal through the Done() channel of the context and exit.
2. Goroutines must call `runtime.Goexit()` to release resources when exiting.
3. Avoid `creating goroutines directly in a for dead loop`, which can generate a large number of useless goroutines and waste resources.
4. Use `atomic operations` in the sync/atomic package to control access to and `release of global resources`.
5. Avoid starting `new goroutines in closures`, which makes it difficult for new goroutines to release resources in closures and prone to leaks.

## Ip6tables problem

The following is the `/etc/hosts` data:

````bash
$ cat /etc/hosts
# 127.0.0.1       localhost
# 127.0.1.1       debian5

# The following lines are desirable for IPv6 capable hosts
# ::1     localhost ip6-localhost ip6-loopback
# ff02::1 ip6-allnodes
# ff02::2 ip6-allrouters
````

I feel very strange. Why does the connection to `localhost` `fail` sometimes?

The reason is that I have the following settings in `ip6tables`:

```bash
$ ip6tables -nL -v | grep policy
# Chain INPUT (policy DROP 0 packets, 0 bytes)
# Chain FORWARD (policy DROP 0 packets, 0 bytes)
# Chain OUTPUT (policy DROP 113 packets, 11924 bytes)
```

Because `I blocked all IPv6 network packets` in the firewall, network testing will fail.

The solution is to `use the hostname debian5 or 127.0.0.1` for testing.

But use `the hostname debian5 first`.

## Benchmark

Use python to build a password generation website, but `do not enable coroutines`

```bash
$ vim netInsightsExplainedz/golangInternet/2multiplexer/example/benchmark/passwords.py
```

The content is as follows:

```python
import random
import string
from flask import Flask

app = Flask(__name__)

@app.route('/password')
def generate_password():
    password_list = []
    for i in range(100):
        password = ''.join(random.choices(string.ascii_letters + string.digits + '#@$', k=12))
        password_list.append(password)
    return str(password_list)

if __name__ == '__main__':
    app.run(port=8080)

```

Use golang to build a password generation website, enable `coroutines` and `race-safe channels`

```bash
$ vim netInsightsExplainedz/golangInternet/2multiplexer/example/benchmark/passwords.go
```

The content is as follows:

```go
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
```

Using Apache HTTP server benchmarking tool

```bash
# go to the test folder
$ cd netInsightsExplainedz/golangInternet/2multiplexer/example/benchmark/

# run the golang program
$ go run passwords.go

# perform stress testing
$ ab -n 500000 -c 1000 http://127.0.0.1:8080/password
# This is ApacheBench, Version 2.3 <$Revision: 1903618 $>
# Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
# Licensed to The Apache Software Foundation, http://www.apache.org/

# Benchmarking debian5 (be patient)
# Completed 50000 requests
# Completed 100000 requests
# Completed 150000 requests
# Completed 200000 requests
# Completed 250000 requests
# Completed 300000 requests
# Completed 350000 requests
# Completed 40000 requests
# Completed 45000 requests
# Completed 50000 requests
# Finished 50000 requests

# Server Software:        
# Server Hostname:        debian5
# Server Port:            8080

# Document Path:          /
# Document Length:        1301 bytes

# Concurrency Level:      1000
# Time taken for tests:   16.292 seconds
# Complete requests:      50000
# Failed requests:        0
# Total transferred:      71000016 bytes
# HTML transferred:       65050000 bytes
# Requests per second:    3068.98 [#/sec] (mean)
# Time per request:       325.841 [ms] (mean)
# Time per request:       0.326 [ms] (mean, across all concurrent requests)
# Transfer rate:          4255.82 [Kbytes/sec] received

# Connection Times (ms)
#               min  mean[+/-sd] median   max
# Connect:        0   11  17.0      2      87
# Processing:    10  312 193.1    278    1483
# Waiting:        2  307 191.5    274    1483
# Total:         10  323 192.3    286    1534

# Percentage of the requests served within a certain time (ms)
#   50%    286
#   66%    388
#   75%    451
#   80%    487
#   90%    579
#   95%    665
#   98%    800
#   99%    907
#  100%   1534 (longest request)

# run the python program
$ python passwords.py

# perform stress testing
$ ab -n 500000 -c 1000 http://127.0.0.1:8080/password
# This is ApacheBench, Version 2.3 <$Revision: 1903618 $>
# Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
# Licensed to The Apache Software Foundation, http://www.apache.org/

# Benchmarking 127.0.0.1 (be patient)
# apr_socket_recv: Connection reset by peer (104)
# Total of 10606 requests completed
```

Just when the network is normal, take this opportunity to run some benchmark tests.
As can be known from the above, if `goroutines are enabled`, the website can `increase the number of connections` that can be supported.
And the code will first `store the generated data in the channel`, which can solve `race problems`.
The data shows that the website will respond within `1.5 seconds (1534 ms)`.
I won't change the Python code, I don't know how to use Python coroutines.



https://juejin.cn/post/6921830629369708552