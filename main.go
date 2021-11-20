package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	url                 string // url to make request to
	outputFile          string // path of the outputfile
	duration            int    // the time between each request in milliseconds
	n                   int    // amount of requests to make
	parallelism         int    // make use of multiple cores by increasing the amount of go routines
	currentRequestCount int    // The amount of requests sent at any point in time
)

func main() {
	setupFlags()
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go loop(&wg)
	}
	wg.Wait()
	elapsed := time.Since(start)

	w := getWriter()
	fmt.Fprintf(w, "%v:\nUrl: %s, Requests: %d, Elapsed: %.3fs, Parallelism: %d\n",
		time.Now().Format(time.RFC1123),
		url,
		currentRequestCount,
		elapsed.Seconds(),
		parallelism)
}

func setupFlags() {
	flag.StringVar(&url, "url", "http://localhost:8000/api/v1/encode", "Specify the url you want to send requests to.")
	flag.StringVar(&outputFile, "o", "", "Specify the path / file name of the output file. (keep empty to print to stdout)")
	flag.IntVar(&duration, "duration", 500, "Specify the duration between each request, in milliseconds")
	flag.IntVar(&n, "n", 5, "Specify the amount of requests to be sent.")
	flag.IntVar(&parallelism, "parallelism", 1, "Specify the amount of go routines to use")
	flag.Parse()

	if parallelism > n {
		parallelism = n
	}
}

func loop(wg *sync.WaitGroup) {
	defer wg.Done()

	for ; currentRequestCount <= n-parallelism; currentRequestCount++ {
		makeGetRequest(http.Client{})
		time.Sleep(time.Millisecond * time.Duration(duration))
	}
}

func makeGetRequest(client http.Client) {
	_, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
}

func getWriter() io.Writer {
	var w io.Writer = os.Stdout
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		w = file
	}
	return w
}
