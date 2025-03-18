package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)


func weatherErrorHandler() (string, error){

	resp , err := http.Get("http://localhost:8080");

	if err!=nil {
		return "", fmt.Errorf("failed to fetch server: %v", err)
	}

	defer resp.Body.Close();
	
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")

		if retryAfter == "" {
			fmt.Println("server is too busy, automatically retries after 2 seconds...")
			time.Sleep(2 * time.Second)
			return weatherErrorHandler()
		} else {
			retryAfterInSecond, err := strconv.Atoi(retryAfter)
			if err != nil {
				fmt.Println("Retry-After header is invalid, automatically retries after 2 seconds...")
				time.Sleep(2 * time.Second)
				return weatherErrorHandler()
			}

			if retryAfterInSecond > 5 {
				return "", fmt.Errorf("server response takes %d seconds, process canceled", retryAfterInSecond)
			} else if retryAfterInSecond >= 1 && retryAfterInSecond <= 5 {
				fmt.Printf("Waiting for %d seconds...\n", retryAfterInSecond)
				time.Sleep(time.Duration(retryAfterInSecond) * time.Second)
				return weatherErrorHandler()
			}
		}

	}else if resp.StatusCode==http.StatusOK{
		body , err :=io.ReadAll(resp.Body)
		if err!=nil{
			return "" , fmt.Errorf("couldn't read the body response: %w",err)
		}
		fmt.Printf("%s" , body)
		return string(body) ,nil
	}else{
		return "" , fmt.Errorf("unkown status code:  %d . ",resp.StatusCode)
	}
	return "" , fmt.Errorf("unexpected error happened")
}

func main(){
	_,err:=weatherErrorHandler();
	if err!=nil {
		fmt.Println("error happened: " , err)
		os.Exit(1);
	}

}