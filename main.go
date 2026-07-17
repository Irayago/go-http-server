package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const PORT = "9999"

// global var for HTTP response delay
var httpTimer = 25 //initialized to 25

func main() {
	RunHttPServer()
	fmt.Println("Exiting Main()...")
}

func RunHttPServer() {
	http.HandleFunc("/", DelayHttpResponse)
	http.HandleFunc("/httptimer", ServeTimerOp)

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


// function to serve /httpTimer API endpoint
func ServeTimerOp(w http.ResponseWriter, r *http.Request) {

	LogRequestMetaData(r)

	// depending on the HTTP method, either get or set the HTTP timer value
	switch r.Method {
	case "PUT":
		// if request body empty, throw error
		if r.ContentLength == 0 {
			_, err := w.Write([]byte("Request body empty for PUT\n"))
			if err != nil {
				fmt.Printf("Error thrown writing zero body error: %v\n", err)
			}
			w.WriteHeader(400) // return 400 bad request
			break
		}	
		defer r.Body.Close()

		// read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading request body: %v\n", err)
			w.WriteHeader(500) // return 500 internal server error
			break
		}

		// unmarshall body into Go map
		var unmarshalledResp map[string]interface{}
		err = json.Unmarshal(body, &unmarshalledResp)
		if err != nil {
			fmt.Printf("Error unmarshalling request body: %v\n", err)
			w.WriteHeader(500) // return 500 internal server error
			break
		}

		userTimerToSet, found := unmarshalledResp["httpTimer"]
		
		// if request body doesn't contain a HTTP timer value
		if !found {
			fmt.Printf("No httpTimer key and value found.")
			w.WriteHeader(400) // return 400 bad request
			break
		}

		
		switch varType := userTimerToSet.(type) {
		case int:
			httpTimer = userTimerToSet.(int)
			w.WriteHeader(200) // return 200 OK 
		default:
			w.Write(fmt.Appendf(nil, "HTTP timer value is not an integer. Of type: %v\n", varType))
			w.WriteHeader(400) // return 400 bad request
		}



	case "GET":
		// GET /httptimer will return httpTimer value
		_, err := w.Write([]byte(strconv.Itoa(httpTimer)))
		if err != nil {
			fmt.Printf("Error thrown writing ServeTimerOp response for GET: %v\n", err)
			w.WriteHeader(500) // return 500 internal server error
		}
	// if HTTP method is other than PUT or GET
	default:
		_, err := w.Write([]byte("Path only supports PUT and GET"))
		if err != nil {
			fmt.Printf("Error thrown writing ServeTimerOp response: %v\n", err)
			w.WriteHeader(500) // return 500 internal server error
		}
	}

}

func DelayHttpResponse(w http.ResponseWriter, r *http.Request) {

	LogRequestMetaData(r)
	curTime := time.Now()

	// Sleep for httpTimer seconds
	time.Sleep(time.Duration(httpTimer) * time.Second)

	afterSleepTime := time.Now()
	elapsed := afterSleepTime.Sub(curTime)
	response := fmt.Sprintf("Time: %v\nMessage Type: Response to client\nUser: %v\nIP: %v\nTime Elapsed: %v\n", afterSleepTime, r.UserAgent(), r.RemoteAddr, elapsed)
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
}

func LogRequestMetaData(r *http.Request) {
	// store client request meta data and print to console
	user := r.UserAgent()
	requestMethod := r.Method
	userIp := r.RemoteAddr

	curTime := time.Now()
	fmt.Printf("Time: %v\nMessage Type: %v\nUser: %v\nIP: %v\n", curTime, requestMethod, user, userIp)
	fmt.Printf("Request Headers: %v\n", r.Header)
	fmt.Println()
}
