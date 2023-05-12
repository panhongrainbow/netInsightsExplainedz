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



https://juejin.cn/post/6921830629369708552