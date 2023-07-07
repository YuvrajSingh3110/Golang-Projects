package main

import "fmt"

type SeoData struct {
	URL        string
	Title      string
	H1         string
	MetaDesc   string
	StatusCode int
}

type parser interface {
}

type DefaultParser struct {
}

var userAgents = ""

func randomUserAgent() {

}

func ExtractSitemapURLs(startURL string) []string {
	WorkList := make(chan []string)
	toCrawl := []string{}
	var n int
	n++

	go func(WorkList <- []string{startURL})()

	for ; n>0; n-- {
		
	}
	list := <- WorkList
	for _, link := list{
		n++
		go func(link string){
			response, err := MakeRequest(link)
			if err != nil{
				panic(err)
			}
			urls, err := ExtractURLs(response)
			if err != nil{
				panic(err)
			}
			sitemapFiles, pages := isSitemap(urls)
			if sitemapFiles != nil{
				WorkList <- sitemapFiles
			}
			for _, page := range pages{
				toCrawl = append(toCrawl, page)
			}
		}(link)
	}
	return toCrawl
}

func MakeRequest() {

}

func ScrapeURLs() {

}

func scrapePage() {

}

func crawlPage() {

}

func getSEOData() {

}

func ScrapeSitemap(URL string) []SeoData {
	results := ExtractSitemapURLs(URL)
	res := ScrapeURLs(results)
	return res
}

func main() {
	p := DefaultParser{}
	results := ScrapeSitemap("")
	for _, res := range results {
		fmt.Println(res)
	}
}
