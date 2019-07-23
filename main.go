package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var urls []string

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func readFromFile(filename string) []string {
	var urls []string
	//read csv file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		os.Exit(0)
	}

	r := csv.NewReader(strings.NewReader(string(file)))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		urls = append(urls, record[0])

	}
	return urls
}

func getBodySize(url string) int {
	// Make a get request
	rs, err := http.Get(url)
	// Process response
	if err != nil {
		log.Println(err)
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		log.Println(err)
	}

	bodySize := len(bodyBytes)
	return bodySize
}

func doWork(ID int, url string, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Println("Worker ", ID, " is working on job ", job)
		size := getBodySize(url)
		duration := time.Duration(rand.Intn(1e3)) * time.Millisecond
		time.Sleep(duration)
		fmt.Println("Worker ", ID, " completed work on job ", job, " within ", duration, "and with html size of", size)
		results <- ID
	}
}

func main() {

	// add error handling for the channel

	// read the url
	filename := flag.String("filename", "url.csv", "file name of the URLs link")
	flag.Parse() //parse flags from command line.

	if *filename != " " {

		//read from file and print the result ✅
		// store the result ✅
		urls = readFromFile(*filename)

		// make a channel ✅.
		jobs := make(chan int, len(urls))
		results := make(chan int, len(urls))

		//  Workers
		for i, url := range urls {
			if isUrl(url) {
				go doWork(i, url, jobs, results)
			}
		}

		// Give them jobs
		for j := 1; j <= len(urls); j++ {
			jobs <- j
		}
		close(jobs)

		// Wait for the results
		for r := 1; r <= len(urls); r++ {
			fmt.Println("Result received from worker: ", <-results)
		}

	}

}
