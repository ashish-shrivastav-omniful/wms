package main

import (
	"log"
	nethttp "net/http"
	"net/url"
	"time"

	"github.com/omniful/go_commons/http"
)
func main(){
// Create transport
transport := &nethttp.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 100,
}

// Initialize client with base URL
client, err := http.NewHTTPClient(
	"my-service",           // client service name
	"https://api.example.com", // base URL
	transport,
	http.WithTimeout(30 * time.Second), // optional timeout
)
if err != nil {
	// Handle error
}
request := &http.Request{
	Url: "/users",
	Body: map[string]interface{}{
		"name": "John Doe",
		"email": "john@example.com",
	},
	Headers: map[string][]string{
		"Content-Type": {"application/json"},
	},
	QueryParams: url.Values{
		"page": []string{"1"},
		"limit": []string{"10"},
	},
	Timeout: 5 * time.Second, // Optional request-specific timeout
}

// For GET request
var response struct{}
resp, err := client.Get(request, &response)
log.Println(resp) 
}