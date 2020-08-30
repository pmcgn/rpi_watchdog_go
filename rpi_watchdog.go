package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	startDelayString := os.Getenv("WATCHDOG_START_DELAY")
	url := os.Getenv("HTTP_HEALTH_CHECK_URL")
	delayString := os.Getenv("HTTP_HEALTH_CHECK_DELAY")
	errorCountMaxString := os.Getenv("HTTP_HEALTH_CHECK_ERRORTHRESHOLD")

	var httpHealthCheckActive = false
	if url != "" {
		httpHealthCheckActive = true
	}

	var startDelay int = 120
	if startDelayString != "" {
		startDelay64, err := strconv.ParseInt(startDelayString, 10, 32)
		if err == nil {
			if startDelay64 < 60 {
				startDelay = 60
			} else {
				startDelay = int(startDelay64)
			}
		}
	}
	fmt.Printf("%v: Trigger ready to fire in %v seconds.\n", time.Now().Format(time.RFC850), startDelay)

	var delay int = 10
	if delayString != "" {
		delay64, err := strconv.ParseInt(delayString, 10, 32)
		if err != nil {
			delay = int(delay64)
		}

	}

	var errorCountMax int = 5
	if errorCountMaxString != "" {
		errorCountMax64, err := strconv.ParseInt(errorCountMaxString, 10, 32)
		if err != nil {
			errorCountMax = int(errorCountMax64)
		}
	}

	resultFromConnectionCheck := make(chan bool)

	/*
	* Ready to Go, just wait for the start delay to avoid infinite restarts in case that the watched URL is never available
	 */
	time.Sleep(time.Duration(startDelay) * time.Second)

	if httpHealthCheckActive {
		go checkConnectionIsAlive(url, delay, resultFromConnectionCheck)
	}

	wdog, err1 := os.Create("/dev/watchdog")
	if err1 != nil {
		fmt.Println("Could not access /dev/watchdog.")
		fmt.Println("Did you run the container with the option --device /dev/watchdog:/dev/watchdog ?")
		fmt.Println("Did you run 'apt install watchdog' on your host ?")
		panic("Exiting")
	}
	defer wdog.Close()

	var errorCount int = 0
	receivedAtLeastOneSuccess := false
	fmt.Printf("%v: Watchdog is now active!\n", time.Now().Format(time.RFC850))
	for {

		fmt.Println("Check write -> errorcount: ", errorCount)
		if errorCount <= errorCountMax {
			_, err2 := wdog.WriteString("WATCH, DOG!\n")

			if err2 != nil {
				log.Fatal(err2)
				panic(err2)
			}
		} else {
			//Quit application properly to stop docker container in case of a missing write access to /dev/watchdog
			//I don't stop the goroutine, as this is really not path, which should be reached
			if errorCount > errorCountMax+2 {
				fmt.Println("Exiting")
				break
			}

		}

		select {
		case msg := <-resultFromConnectionCheck:
			if msg == true {
				fmt.Println("Received true")
				errorCount = 0
				receivedAtLeastOneSuccess = true
			} else {
				fmt.Println("received false")
				if receivedAtLeastOneSuccess {
					errorCount++
				}
			}
		default:
			//do nothing....
		}
		time.Sleep(5 * time.Second) //the watchdog hardware will trigger a reset after 15s
	}
}

func checkConnectionIsAlive(url string, interval int, isOnlineNotification chan<- bool) {
	for {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		_, err := http.Get(url)
		if err != nil {
			isOnlineNotification <- false
		} else {
			isOnlineNotification <- true
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
