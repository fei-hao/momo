package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	PROXY_LIST_URL = "http://www.xicidaili.com/wt/"
	USER_AGENT     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
)

var (
	MOMO_SHARE_LINK string
	mu              sync.Mutex
	wd              sync.WaitGroup
	count           int
)

func main() {
	text := os.Args[1]
	MOMO_SHARE_LINK = text
	fmt.Println("your share link: " + MOMO_SHARE_LINK)

	list := getIpList()

	for _, proxy := range list {
		wd.Add(1)
		go visit(proxy)
	}
	wd.Wait()
}

func getIpList() []string {

	var ipList []string
	// get ip list page
	client := &http.Client{}
	req, err := http.NewRequest("GET", PROXY_LIST_URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// parse html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {

		address := s.Find("td").First().Next().Text()
		port := s.Find("td").First().Next().Next().Text()
		if address != "" && port != "" {
			proxy := address + ":" + port
			ipList = append(ipList, proxy)
		}
	})

	return ipList
}

func visit(proxy string) {
	defer wd.Done()
	proxyUrl, err := url.Parse("http://" + proxy)
	if err != nil {
		log.Println(err)
	}
	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	timeout := time.Duration(10 * time.Second)
	client := &http.Client{Transport: tr, Timeout: timeout}

	req, err := http.NewRequest("GET", MOMO_SHARE_LINK, nil)
	if err != nil {
		log.Println("Request error:", err)
		return
	}
	req.Header.Add("User-Agent", USER_AGENT)

	client.Do(req)

	mu.Lock()
	fmt.Printf("No %d, visited by: %s \n", count, proxy)
	count++
	mu.Unlock()

}
