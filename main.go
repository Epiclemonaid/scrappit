package main

import (
  "fmt"
  "reddit-scraper/http"
  "reddit-scraper/reddit"
)



func main() {
  fmt.Printf("Scraper v0.1\n")
  fmt.Printf("Created by Curtis Li\n")

  subreddits := []string {"/r/funny", "/r/AskReddit", "/r/gfycats"}

  for _, subreddit := range subreddits {
    fmt.Println(subreddit)
    listing := reddit.ListingJson{}
    http.GetJson("http://www.reddit.com/" + subreddit + "/.json", &listing)
    fmt.Printf("%v\n", listing)
    downloadUrls := reddit.DownloadPosts(listing.Data.Children)
    fmt.Println(downloadUrls)
  }
}

