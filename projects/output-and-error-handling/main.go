package main

import (
	"fmt"
	"net/http"
)


func weatherErrorHandler() error{

	resp , err := http.Get("http://localhost:8080");

	if err!=nil {
		return fmt.Errorf("failed to fetch server: %v", err)
	}

	defer resp.Body.Close();

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter:= resp.Header.Get("Retry-After");

		if retryAfter=="" {
			// here ------------
			return nil
		}
	}
}

func main(){

}