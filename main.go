package main

import (
  "fmt"
  "os"
  "strconv"
  "regexp"
  "encoding/json"
  "net/url"
  "reddit-scraper/util"
  "reddit-scraper/http"
  "reddit-scraper/reddit"
)


const configFile = "conf.json"
const redditUrl = "http://www.reddit.com"

type Configuration struct {
  Subreddits []SubredditConfig `json:"subreddits"`
  OutputPath string `json:"outputPath"`
  Stats bool `json:"stats"`
  MaxFileSize int `json:"maxFileSize"`
  MinFileSize int `json:"minFileSize"`
  FileTypes []string `json:"fileTypes"`
}

type SubredditConfig struct {
  Name string `json:"subredditName"`
  Limit int `json:"numberOfPosts"`
  SortBy string `json:"sortBy"`
  SearchFor string `json:"searchFor"`
}


func main() {
  fmt.Println("Scraper v0.1")
  fmt.Println("Created by Curtis Li")

  // Get configuration settings
  config := configSettings(configFile)

  // Loop through subreddit list
  for _, subreddit := range config.Subreddits {
    fmt.Println("-----------------------------\n")
    fmt.Println(subreddit.Name)

    // Get subreddit JSON
    listing := reddit.ListingJson{}
    redditReq := createRedditJsonReq(subreddit)
    fmt.Println("Requesting data from", redditReq)

    http.GetJson(redditReq, &listing)

    // Get download links
    posts := reddit.DownloadPosts(listing.Data.Children)
    fmt.Println(len(posts), "posts to download")

    // Get output directory path
    outputPath := config.OutputPath
    if outputPath == "" {
      outputPath = "output/"
    }
    r, _ := regexp.Compile(`.*/$`)
    if !r.MatchString(outputPath) {
      outputPath = outputPath + "/"
    }
    outputPath = outputPath + subreddit.Name[3:]
    err := os.MkdirAll(outputPath, 0755)
    util.Check(err)

    // Download to folder
    for _, post := range posts {
      outputFile := outputPath + post.Name + ".webm"
      fmt.Println("Downloading", post.Url, "to", outputFile)
      err := http.DownloadFile(outputFile, post.Url)
      util.Check(err)
    }
  }
}


func configSettings(filename string) Configuration {
  // Open config file
  configuration := Configuration{}
  file, err := os.Open(filename)

  // File does not exist, create it
  if os.IsNotExist(err) {
    // Create file
    fmt.Println("Initiating new configuration file...")
    file, err = os.Create(filename)
    util.Check(err)

    // Setup and encode the JSON
    var b []byte
    configuration.Subreddits = append(configuration.Subreddits, SubredditConfig{"/r/subreddit1", 50, "new", ""}, SubredditConfig{"/r/subreddit2", 20, "hot", ""})
    configuration.OutputPath = "Path/To/Output/Folder"
    b, err = json.MarshalIndent(configuration, "", "    ")
    util.Check(err)

    // Write to the new file
    _, err = file.Write(b)
    util.Check(err)

    // Close the fd
    err = file.Close()
    util.Check(err)

    // Exit
    fmt.Println("Please edit 'config.json'")
    os.Exit(0)
  }

  util.Check(err)

  // Parse JSON
  err = json.NewDecoder(file).Decode(&configuration)
  util.Check(err)

  return configuration
}

func createRedditJsonReq(subreddit SubredditConfig) string {
  // Base URL
  redditReq, err := url.Parse(redditUrl + subreddit.Name)
  util.Check(err)

  // Search vs Sort
  if subreddit.SearchFor != "" {
    redditReq.Path = redditReq.Path + "/search"
  } else if subreddit.SortBy != "" {
    redditReq.Path = redditReq.Path + "/" + subreddit.SortBy
  }

  // End URL
  redditReq.Path = redditReq.Path + "/.json"

  // Query parameters
  values := url.Values{}

  // Limiting
  values.Set("limit", "20")
  if subreddit.Limit != 0 {
    values.Set("limit", strconv.Itoa(subreddit.Limit))
  }

  // Searching
  if subreddit.SearchFor != "" {
    values.Set("q", subreddit.SearchFor)
    values.Set("restrict_sr", "on")
    values.Set("sort", "relevance")
    values.Set("t", "all")
  }

  redditReq.RawQuery = values.Encode()

  return redditReq.String()
}
