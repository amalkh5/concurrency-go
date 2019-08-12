package main

import (
	"encoding/csv"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Amalkh5/concurrency-go/pool"
)

func readFromFile(filename string) []string {
	var urls []string
	//read csv file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		// log.Fatal is enough
		log.Fatal("File reading error", err)
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

	// early return is a cleaner way to write code,
	// it avoids nesting and lets the reader knows the errors
	// for more info https://blog.timoxley.com/post/47041269194/avoid-else-return-early
	if len(*filename) <= 0 {
		log.Fatal("filename must be more than zero")
		// log.Fatal exits the process so no need to "return", i added it just so you could learn the return-early pattern
		return
	}

	urls := readFromFile(*filename)
	p := pool.NewPool()
	p.StartTheWorker(urls)
}
