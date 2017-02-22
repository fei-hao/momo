package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const (
	//replace this with yours
	MOMO_SHARE_LINK = "http://www.maimemo.com/share/page/?uid=749481&pid=774"

	PROXY_LIST_URL  = "http://www.xicidaili.com/wt/"
	USER_AGENT      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
)

var (
	mu    sync.Mutex
	count int
)

func main() {

	list := getIpList()
	ch := make(chan string)

	for _, proxy := range list {
		go visit(proxy, ch)
	}

	for count != len(list) {
		fmt.Println(<-ch)
	}
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

func visit(proxy string, ch chan string) {

	proxyUrl, err := url.Parse("http://" + proxy)
	if err != nil {
		log.Println(err)
	}
	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", MOMO_SHARE_LINK, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("User-Agent", USER_AGENT)

	client.Do(req)
	ch <- fmt.Sprintf("No %d, visited by: %s \n", count, proxy)

	mu.Lock()
	count++
	mu.Unlock()

}
