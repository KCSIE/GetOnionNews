package main

import (
	"GetOnionNews/fetcher"
	"GetOnionNews/parser"
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	start_seed := "https://www.theonion.com/breaking-news/news-in-brief?startIndex="
	start_page := 0 // e.g. start_page=0&end_page=0 => 0(20)
	end_page := 6840 // e.g. start_page=0&end_page=40 => 0(20)+20(20)+40(20)

	links,err := fetcher.GetAllPagesLinks(start_seed,start_page,end_page)
	if err != nil {
		fmt.Printf("Error: %s", err)
		panic(err)
	}

	newsinfos,err := parser.GetAllNewsInfo(links)
	if err != nil {
		fmt.Printf("Error: %s", err)
		panic(err)
	}

    csvFile, err := os.Create("The Onion's Breaking News - News In Brief.csv")
    if err != nil {
        panic(err)
    }
    defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	clmnm := []string{
		"Title",
		"Published Time",
		"Content",
	}
	writer.Write(clmnm)
    for _, newsinfo := range newsinfos {
        line := []string{
            newsinfo.Title,
			newsinfo.PubTime,
			newsinfo.Content,
        }
        err := writer.Write(line)
        if err != nil {
            panic(err)
        }
    }
    writer.Flush()
}