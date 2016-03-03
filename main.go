package main

import (
  "fmt"
  "os"
  "encoding/json"
  "reddit-scraper/util"
  "reddit-scraper/http"
  "reddit-scraper/reddit"
)


const configFile = "conf.json"

type Configuration struct {
  Username string `json:"username"`
  Password string `json:"password"`
  Subreddits []string `json:"subreddits"`
  Output string `json:"output"`
  Stats bool `json:"stats"`
}


func main() {
  fmt.Printf("Scraper v0.1\n")
  fmt.Printf("Created by Curtis Li\n")

  // Read configuration file
  config := readConfig(configFile)
  fmt.Printf("%v\n", config)

  // Loop through subreddit list
  for _, subreddit := range config.Subreddits {
    fmt.Println(subreddit)

    // Get subreddit JSON
    listing := reddit.ListingJson{}
    http.GetJson("http://www.reddit.com/" + subreddit + "/.json", &listing)

    // Get downloads
    downloadUrls := reddit.DownloadPosts(listing.Data.Children)
    fmt.Println(downloadUrls)
  }
}


func readConfig(filename string) Configuration {
  // Open config file
  file, err := os.Open(filename)
  util.Check(err)

  // Parse JSON
  configuration := Configuration{}
  err = json.NewDecoder(file).Decode(&configuration)
  util.Check(err)

  return configuration
}
