package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func handleWeatherRequest() (string, error) {

	resp, err := http.Get("http://localhost:8080")

	if err != nil {

		if os.IsTimeout(err) {
			return "", fmt.Errorf("request timed out")
		}
		if err == io.EOF {
			return "", fmt.Errorf("connection terminated unexpectedly")
		}

		return "", fmt.Errorf("failed to fetch server: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")

		if retryAfter == "" {
			fmt.Println("server is too busy, automatically retries after 2 seconds...")

			// 2 seconds is a good choice because it's short enough to minimize user delay while still giving the server a
			//  brief moment to recover, avoiding overloading it with frequent requests.
			time.Sleep(2 * time.Second)
			return handleWeatherRequest()
		} else {
			retryAfterInSecond, err := strconv.Atoi(retryAfter)
			if err != nil {

				retryTime,err:= time.Parse(time.RFC1123, retryAfter);

				if err!=nil{
					fmt.Println("Retry-After header is invalid, automatically retries after 2 seconds...")
					time.Sleep(2 * time.Second)
					return handleWeatherRequest()
				}

				waitingDuration:=time.Until(retryTime);
				if waitingDuration.Seconds() > 5 {
				return "", fmt.Errorf("server response takes %.0f seconds, process canceled", waitingDuration.Seconds())
				}else if waitingDuration.Seconds()>0{
					fmt.Printf("Waiting for %.0f seconds...\n", waitingDuration.Seconds());
					time.Sleep(waitingDuration);
					return handleWeatherRequest();
				}else{
					fmt.Println("retrying to get the weather... please wait")
				}

				return handleWeatherRequest();
			}

			if retryAfterInSecond > 5 {
				return "", fmt.Errorf("server response takes %d seconds, process canceled", retryAfterInSecond)
			} else if retryAfterInSecond >= 1 && retryAfterInSecond <= 5 {
				fmt.Printf("Waiting for %d seconds...\n", retryAfterInSecond)
				time.Sleep(time.Duration(retryAfterInSecond) * time.Second)
				return handleWeatherRequest()
			}
		}

	} else if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("couldn't read the body response: %w", err)
		}
		fmt.Printf("%s", body)
		return string(body), nil
	} else {
		return "", fmt.Errorf("unkown status code:  %d . ", resp.StatusCode)
	}
	return "", fmt.Errorf("unexpected error happened")
}

func main() {
	_, err := handleWeatherRequest()
	if err != nil {
		fmt.Println("error happened: ", err)
		os.Exit(1)
	}

}
