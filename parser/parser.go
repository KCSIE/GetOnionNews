package parser

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type News struct {
	Title         string    
	PubTime       string    
	Content       string    
}

type msg struct{
	news News
	err error
}

func GetAllNewsInfo(links []string)([]News,error){
	var allNewsInfo []News
	ch := make(chan msg)
	for i, link := range links{
		fmt.Printf("Handling Link %d: %s\n", i, link)
		go GetNewsInfo(link,ch)
	}
	for i := range links{
		rmsg := <-ch
		if rmsg.err != nil {
			log.Fatal(rmsg.err)
			return nil,rmsg.err
		}
		allNewsInfo = append(allNewsInfo, rmsg.news)
		fmt.Printf("Link %d Done\n", i)
	}
	return allNewsInfo,nil
}

func GetNewsInfo(link string, ch chan msg){
	client := http.Client{}

	req, _ := http.NewRequest("GET", link, nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond) 
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		ch <- msg{News{},err}
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
		ch <- msg{News{},err}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		ch <- msg{News{},err}
	}

	title := doc.Find(".sc-157agsr-1 header .sc-1efpnfq-0").Text()
	// fmt.Printf("Title: %s\n", title)
	pubtime, _ := doc.Find(".sc-157agsr-0 .uhd9ir-0").Attr("datetime")
	// fmt.Printf("Published Time: %s\n", pubtime)
	content := doc.Find(".xs32fe-0 p").Text()
	// fmt.Printf("Content: %s\n", content)

	if title == "" {
		log.Printf("Something Went Wrong")
	}

	newsinfo := News{
		Title:   title,
		PubTime: pubtime,
		Content: content,
	}

	rmsg := msg{
		news:newsinfo,
		err:nil,
	}

	ch <- rmsg
}