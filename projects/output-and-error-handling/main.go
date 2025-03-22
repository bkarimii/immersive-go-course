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
			fmt.Fprintf(os.Stderr ,"request timed out.\n");
			return "", fmt.Errorf("request timed out")
		}
		if err == io.EOF {
			fmt.Fprintf(os.Stderr ,"connection terminated unexpectedly.\n")
			return "", fmt.Errorf("connection terminated unexpectedly")
		}
		fmt.Fprintf(os.Stderr ,"failed to fetch server: %v\n", err);
		return "", fmt.Errorf("failed to fetch server: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")

		if retryAfter == "" {
			fmt.Fprintf(os.Stderr,"server is too busy")

			// 2 seconds is a good choice because it's short enough to minimize user delay while still giving the server a
			//  brief moment to recover, avoiding overloading it with frequent requests.
			time.Sleep(2 * time.Second)
			return handleWeatherRequest()
		} else {
			retryAfterInSecond, err := strconv.Atoi(retryAfter)
			if err != nil {

				retryTime,err:= time.Parse(time.RFC1123, retryAfter);

				if err!=nil{
					fmt.Fprintf(os.Stderr, "error happened: %v\n" , err);
					time.Sleep(2 * time.Second)
					return handleWeatherRequest()
				}

				waitingDuration:=time.Until(retryTime);
				if waitingDuration.Seconds() > 5 {

					fmt.Fprintf(os.Stderr,"server response takes %.0f seconds, process canceled", waitingDuration.Seconds())

					return "", fmt.Errorf("server response takes %.0f seconds, process canceled", waitingDuration.Seconds())

				}else if waitingDuration.Seconds()>0{

					fmt.Fprintf(os.Stderr ,"Waiting for %.0f seconds...\n", waitingDuration.Seconds())

					time.Sleep(waitingDuration);
					return handleWeatherRequest();
				}else{
					fmt.Fprintf(os.Stderr,"retrying to get the weather... please wait")
				}

				return handleWeatherRequest();
			}

			if retryAfterInSecond > 5 {
				fmt.Fprintf(os.Stderr , "server response takes %d seconds, process canceled \n", retryAfterInSecond);

				return "", fmt.Errorf("server response takes %d seconds, process canceled", retryAfterInSecond)

			} else if retryAfterInSecond >= 1 && retryAfterInSecond <= 5 {
				fmt.Fprintf(os.Stderr , "Waiting for %d seconds...\n", retryAfterInSecond);
				time.Sleep(time.Duration(retryAfterInSecond) * time.Second)
				return handleWeatherRequest()
			}
		}

	} else if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr , "couldn't read the body response: %v\n", err)
			return "", fmt.Errorf("couldn't read the body response: %w", err)
		}
		return string(body), nil
	} else {
		fmt.Fprintf(os.Stderr ,"unkown status code:  %d .\n", resp.StatusCode)
		return "", fmt.Errorf("unkown status code:  %d . ", resp.StatusCode)
	}

	fmt.Fprintf(os.Stderr ,"unexpected error happened");
	return "", fmt.Errorf("unexpected error happened");
}

func main() {
	weather, err := handleWeatherRequest()
	if err != nil {
		fmt.Fprintf(os.Stderr,"error happened: %v\n",err);
		os.Exit(1)
	}

	fmt.Printf("%s", weather);

}
