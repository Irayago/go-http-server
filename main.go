package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
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
	// define HTTP API server routes
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
		IdleTimeout:  120 * time.Second, // 2min idle timeout for keep-alive connections
	}

	// creating tcp keep-alive config to disable keep-alive connections
	tcpConfig := net.ListenConfig{
		KeepAlive: -1,
	}

	// create tcp listener with keep-alive config
	tcpListener, err := tcpConfig.Listen(context.Background(), "tcp", curListeningPort)
	if err != nil {
		log.Fatalf("Error thrown creating tcpListener: %v\n", err)
	}

	// start HTTP server with tcp listener
	err = httpServer.Serve(tcpListener)
	if err != nil {
		log.Fatalf("Error thrown from httpServer.Serve(): %v\n", err)
	}
}


// function to serve /httpTimer API endpoint
func ServeTimerOp(w http.ResponseWriter, r *http.Request) {

	LogRequestMetaData(r)

	// depending on the HTTP method, either get or set the HTTP timer value
	switch r.Method {
	case "PUT":
		// if request body empty, throw error
		if r.ContentLength <= 0 {
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

		fmt.Println(string(body))

		// unmarshall body into Go map
		var unmarshalledResp map[string]interface{}
		err = json.Unmarshal(body, &unmarshalledResp)
		if err != nil {
			fmt.Printf("Error unmarshalling request body: %v\n", err)
			w.WriteHeader(400) // return 400 bad request
			break
		}

		userTimerToSet, found := unmarshalledResp["httpTimer"]
		
		// if request body doesn't contain a HTTP timer value
		if !found {
			fmt.Printf("No httpTimer key and value found in payload.")
			w.WriteHeader(400) // return 400 bad request
			break
		}

		
		switch varType := userTimerToSet.(type) {
		case float64:
			valToCheck := int(userTimerToSet.(float64)) // set global httpTimer value
			if valToCheck < 0 { // check if value is negative
				w.WriteHeader(400)
				_, err := w.Write(fmt.Appendf(nil, "HTTP timer value cannot be negative. Value provided: %v\n", valToCheck))
				if err != nil {
					fmt.Printf("Error thrown writing error HTTP timer negative value response: %v\n", err)
				}
				break
			}
			httpTimer = valToCheck
			w.WriteHeader(200) // return 200 OK 
			_, err := w.Write(fmt.Appendf(nil, "HTTP timer value set to: %v seconds\n", httpTimer))
			if err != nil {
				fmt.Printf("Error thrown writing success response: %v\n", err)
			}
		default:	
			w.WriteHeader(400) // return 400 bad request	
			_, err := w.Write(fmt.Appendf(nil, "HTTP timer value is not an integer. Of type: %T\n", varType))	
			if err != nil {
				fmt.Printf("Error thrown writing error HTTP timer non-integer value response: %v\n", err)
			}
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
// function to serve / API endpoint
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

	// if the request path is /bigresponse, send a big response
	if r.URL.Path == "/bigresponse" {
		bigRes := make([]byte, 15*1024) // changed from 10Mb to 15Kb response
		_, err = w.Write(bigRes)
		if err != nil {
			fmt.Printf("Error thrown writing big response: %v\n", err)
		}
	}

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
