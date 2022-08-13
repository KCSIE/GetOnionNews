package fetcher

import (
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type msg struct{
	urls []string
	err error
}

func GetAllPagesLinks(url string,start int, end int)([]string,error){
	var allLinks []string
	ch := make(chan msg)
	page1,page2 := start,start
	for page1 <= end {
		pageurl := url+strconv.Itoa(page1)
		go GetNewsLinks(pageurl,ch)
		page1 = page1 + 20
	}
	for page2 <= end {
		rmsg := <-ch
		if rmsg.err != nil {
			log.Fatal(rmsg.err)
			return nil,rmsg.err
		}
		allLinks = append(allLinks, rmsg.urls...)
		page2 = page2 + 20
	}
	return allLinks,nil
}


func GetNewsLinks(url string, ch chan msg){
	client := http.Client{}
	req,_ := http.NewRequest("GET", url, nil) // https://www.theonion.com/breaking-news/news-in-brief?startIndex={0-6820}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	
	resp,err := client.Do(req)	
	if err != nil {
		log.Fatal(err)
		ch <- msg{nil,err}
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		log.Fatalf("Status code error: %d %s", resp.StatusCode, resp.Status)
		ch <- msg{nil,err}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		ch <- msg{nil,err}
	}
	
	var pageLinks []string

	doc.Find(".sc-17uq8ex-0 article .cw4lnv-5").Each(func(i int, s *goquery.Selection) {
		link,_ := s.Find("a").Attr("href")
		// fmt.Printf("Link %d: %s\n", i, link)
		pageLinks = append(pageLinks, link)
	})
	
	rmsg := msg{
		urls:pageLinks,
		err:nil,
	}
	
	ch <- rmsg
}
