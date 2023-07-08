package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SeoData struct {
	URL        string
	Title      string
	H1         string
	MetaDesc   string
	StatusCode int
}

type Parser interface {
	getSEOData(response *http.Response) (SeoData, error)
}

//it is an empty struct for implementing default parser
type DefaultParser struct {
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

//it helps our server looks like a browser when making any request to any website
func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents) //to randomly select a user agent
	return userAgents[randNum]
}

func isSitemap(url []string) ([]string, []string) {
	sitemapFiles := []string{}
	pages := []string{}
	for _, page := range url {
		foundSitemap := strings.Contains(page, "xml")
		if foundSitemap == true {
			fmt.Println("Found sitemap", page)
			sitemapFiles = append(sitemapFiles, page)
		} else {
			pages = append(pages, page)
		}
	}
	return sitemapFiles, pages
}

func ExtractSitemapURLs(startURL string) []string {
	WorkList := make(chan []string)
	toCrawl := []string{}
	var n int
	n++

	go func() { WorkList <- []string{startURL} }()

	for ; n > 0; n-- {
		list := <-WorkList
		for _, link := range list {
			n++
			go func(link string) {
				response, err := MakeRequest(link)
				if err != nil {
					panic(err)
				}
				urls, err := ExtractURLs(response)
				if err != nil {
					panic(err)
				}
				sitemapFiles, pages := isSitemap(urls)
				if sitemapFiles != nil {
					WorkList <- sitemapFiles
				}
				for _, page := range pages {
					toCrawl = append(toCrawl, page)
				}
			}(link)
		}
	}
	return toCrawl
}

func MakeRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second, //to delay between our requests so that the website does not get overloaded
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", randomUserAgent())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ScrapeURLs(urls []string, parser Parser, concurrency int) []SeoData {
	tokens := make(chan struct{}, concurrency)
	var n int
	WorkList := make(chan []string)
	results := []SeoData{}

	go func() { WorkList <- urls }()

	for ; n > 0; n-- {
		list := <-WorkList
		for _, url := range list {
			if url != "" {
				n++
				go func(url string, token chan struct{}) {
					log.Printf("Requesting URL: %s", url)
					res, err := scrapePage(url, tokens, parser)
					if err != nil {
						log.Printf("Error in url: %s", url)
					}else{
						results = append(results, res)
					}
					WorkList <- []string{}
				}(url, tokens)
			}
		}
	}
	return results
}

func ExtractURLs(response *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []string{}
	sel := doc.Find("loc")
	for i := range sel.Nodes {
		loca := sel.Eq(i)
		res := loca.Text()
		results = append(results, res)
	}
	return results, nil
}

func scrapePage(url string, token chan struct{}, parser Parser) (SeoData, error) {
	response, err := crawlPage(url, token)
	if err != nil {
		return SeoData{}, err
	}
	data, err := parser.getSEOData(response)
	if err != nil {
		return SeoData{}, err
	}
	return data, nil
}

func crawlPage(url string, token chan struct{}) (*http.Response, error) {
	token <- struct{}{}
	res, err := MakeRequest(url)
	if err != nil{
		return nil, err
	}
	return res, nil
}

func (d DefaultParser) getSEOData(res *http.Response) (SeoData, error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return SeoData{}, err
	}
	result := SeoData{}
	result.URL = res.Request.URL.String()
	result.StatusCode = res.StatusCode
	result.Title = doc.Find("title").First().Text()
	result.H1 = doc.Find("h1").First().Text()
	result.MetaDesc, _ = doc.Find("meta[name^=description]").Attr("content")
	return result, nil
}

func ScrapeSitemap(URL string, parser Parser, concurrency int) []SeoData {
	results := ExtractSitemapURLs(URL)
	res := ScrapeURLs(results, parser, concurrency)
	return res
}

func main() {
	fmt.Println("Siremap Crawler")
	p := DefaultParser{}
	results := ScrapeSitemap("https://www.quicksprout.com/sitemap.xml", p, 10)
	for _, res := range results {
		fmt.Println(res)
	}
}
