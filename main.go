package main

import (
	"fmt"
	"net/http"
	"time"
)

const PORT = "9999"

func main() {
	RunHTTPServer()
	fmt.Println("Exiting Main()...")
}

func RunHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// store client request meta data
		user := r.UserAgent()
		requestMethod := r.Method
		userIp := r.RemoteAddr

		curTime := time.Now()

		fmt.Printf("Time: %v\nMessage Type: %v\nUser: %v\nIP: %v\n", curTime, requestMethod, user, userIp)
		fmt.Printf("Request Headers: %v\n", r.Header)
		fmt.Println()

		// Sleep for 15 seconds
		time.Sleep(15 * time.Second)

		afterSleepTime := time.Now()
		elapsed := afterSleepTime.Sub(curTime)
		response := fmt.Sprintf("Time: %v\nMessage Type: Response to client\nUser: %v\nIP: %v\nTime Elapsed: %v\n", afterSleepTime, user, userIp, elapsed)
		_, err := fmt.Print(response)
		if err != nil {
			fmt.Printf("Error thrown from printing response: %v\n", err)
		}

		fmt.Printf("Response Headers: %v\n", w.Header())
		_, err = w.Write([]byte(response))
		if err != nil {
			fmt.Printf("Error thrown writing response: %v\n", err)
		}
		fmt.Println()

	})

	curListeningPort := ":" + PORT
	fmt.Printf("HTTP Server listening on port: %v\n", PORT)
	fmt.Println()

	// define server config
	httpServer := &http.Server{
		Addr:         curListeningPort,
		Handler:      nil,
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  0,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("Error thrown from httpServer.ListenAndServe(): %v\n", err)
	}
}
