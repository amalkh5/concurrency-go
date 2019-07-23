package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Amalkh5/concurrency-go/pool"
)

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

func main() {

	filename := flag.String("filename", "url.csv", "file name of the URLs link")
	flag.Parse() //parse flags from command line.

	if len(*filename) != 0 {

		urls := readFromFile(*filename)
		p := pool.NewPool()
		p.StartTheWorker(urls)

	}

}
