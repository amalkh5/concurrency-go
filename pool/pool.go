package pool

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Job Structure that wraps Jobs information
type Job struct {
	id  int
	url string
}

type Pool struct {
	jobs    chan Job
	results chan string
	done    chan bool
}

// NewPool create the job,result chan
func NewPool() *Pool {
	log.Println("Stating new pool.")
	p := &Pool{}
	p.jobs = make(chan Job, 10)
	p.results = make(chan string, 10)
	return p

}
func (p *Pool) StartTheWorker(urls []string) {
	log.Print("worker pool starting")

	startTime := time.Now()
	go p.allocate(urls)
	p.done = make(chan bool)
	go p.collectResults()
	go p.worker()
	<-p.done
	endTime := time.Now()
	diff := endTime.Sub(startTime)

	log.Printf("total time taken: [%f] seconds", diff.Seconds())
}

func (p *Pool) allocate(urls []string) {
	defer close(p.jobs)
	log.Printf("Allocating [%d] resources", len(urls))
	for i, url := range urls {
		if isUrl(url) {
			// log.Print(i)
			job := Job{id: i, url: url}
			p.jobs <- job
		}
	}
	log.Printf("Done Allocating.")
}

func (p *Pool) worker() {
	defer close(p.results)
	for job := range p.jobs {
		start := time.Now()
		fmt.Println("working on job ID ", job.id)
		nbytes := getBodySize(job.url)
		secs := time.Since(start).Seconds()
		p.results <- fmt.Sprintf("fetch time %.2fs page sizes %7d  url %s", secs, nbytes, job.url)
	}
}

func (p *Pool) collectResults() {
	// Wait for the results
	for result := range p.results {
		fmt.Println("Result received from worker: ", result)
	}
	p.done <- true
}

//////////////////// mv to another file /////////////////

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
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
