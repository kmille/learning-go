package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const urlPrefixHTTP string = "http://"
const urlPrefixHTTPS string = "https://"
const defaultUserAgent = "gurl/0.0.1"

type additionalHeaders []string

func (ah *additionalHeaders) Set(value string) error {
	*ah = append(*ah, value)
	return nil
}

func (ah *additionalHeaders) String() string {
	return "additionon headers"
}

var (
	url             string
	method          *string
	outputFile      *string
	userAgent       *string
	postData        *string
	followRedirects *bool
	headers         additionalHeaders
)

func shouldFollowRedirect(req *http.Request, via []*http.Request) error {
	if *followRedirects {
		return nil
	}
	return http.ErrUseLastResponse
}

func gurlUsage() {
	fmt.Printf("%s - curl implemenation in Go\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Printf("  url\n\turl to visit\n")
}

func parseCommandLineArguments() {
	method = flag.String("X", "GET", "http method to use")
	outputFile = flag.String("w", "-", "output file")
	userAgent = flag.String("A", defaultUserAgent, "user agent")
	postData = flag.String("d", "", "post data")
	followRedirects = flag.Bool("L", false, "follow redirects")
	flag.Var(&headers, "H", "additional headers")
	flag.Usage = gurlUsage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	url = flag.Arg(0)

	if !strings.HasPrefix(url, urlPrefixHTTP) && !strings.HasPrefix(url, urlPrefixHTTPS) {
		url = urlPrefixHTTP + url
	}
}

func main() {

	parseCommandLineArguments()
	/*
		requestURL, err := goUrl.Parse(url)
		if err != nil {
			fmt.Println("Error parsing url", err)
		}
	*/
	//fmt.Println(requestURL.Host)

	client := &http.Client{
		CheckRedirect: shouldFollowRedirect,
	}

	if len(*postData) != 0 && *method == "GET" {
		*method = "POST"
	}

	req, err := http.NewRequest(*method, url, bytes.NewBufferString(*postData))
	if err != nil {
		fmt.Println("Error in http.Get:", err)
		os.Exit(1)
	}

	for _, header := range headers {
		parsedHeader := strings.Split(header, ":")
		req.Header.Add(parsedHeader[0], parsedHeader[1])
	}

	req.Header.Add("User-Agent", *userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Doing the Request", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading http body response", err)
		os.Exit(1)
	}

	if *outputFile == "-" {
		fmt.Printf("%s", body)
	} else {
		ioutil.WriteFile(*outputFile, body, 0644)
	}

	os.Exit(0)

}
