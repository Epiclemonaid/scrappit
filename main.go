package main

import (
  "fmt"
  "reddit-scraper/http"
//  "reddit-scraper/gfycat"
  "reddit-scraper/reddit"
)



func main() {
  fmt.Printf("Scraper v0.1\n")
  fmt.Printf("Created by Curtis Li\n")

  subreddits := []string {"/r/funny", "/r/AskReddit"}


  /*
  item := gfycat.GfyJson{}
  http.GetJson("https://gfycat.com/cajax/get/ImpoliteImmenseCaracal", &item)
  fmt.Printf("Results: %s\n", item.GfyItem.GifUrl)

  http.DownloadFile("file", item.GfyItem.GifUrl)
  */

  //urls, _ := http.Crawl("http://www.reddit.com/r/funny")

  //urls = reddit.GetRedditUrls(urls)

  //for _, url := range urls {
    //fmt.Println(url)
  //}

  for _, subreddit := range subreddits {
    fmt.Println(subreddit)
    listing := reddit.ListingJson{}
    http.GetJson("http://www.reddit.com/" + subreddit + "/.json", &listing)
    fmt.Printf("Results:\n %v\n\n", listing)
    reddit.GetDownloadUrls(listing.Data.Children)
  }
}

