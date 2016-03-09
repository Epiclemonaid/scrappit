package main

import (
  "fmt"
  "os"
  "time"
  "flag"
  "strconv"
  "regexp"
  "encoding/json"
  "net/url"
  "reddit-scraper/util"
  "reddit-scraper/http"
  "reddit-scraper/reddit"
)


const defaultConfigFile = "conf.json"
const redditUrl = "http://www.reddit.com"
const defaultMaxThreads = 10

type Configuration struct {
  Subreddits []SubredditConfig `json:"subreddits"`
  OutputPath string `json:"outputPath"`
  MaxThreads int `json:"maxThreads"`
  Stats bool `json:"stats"`
  FileTypes []string `json:"fileTypes"`
}

type SubredditConfig struct {
  Name string `json:"subredditName"`
  Limit int `json:"numberOfPosts"`
  SortBy string `json:"sortBy"`
  Time string `json:"time"`
  SearchFor string `json:"searchFor"`
  CustomFolderName string `json:"customFolderName"`
}


func main() {
  fmt.Println("Scraper v0.1")
  fmt.Println("Created by Curtis Li")

  // Command-line flags
  configFile := flag.String("c", defaultConfigFile, "Path to configuration file")
  maxThreads := flag.Int("t", 0, "Maximum number of concurrent downloads")
  getHelp := flag.Bool("h", false, "Help")
  flag.Parse()

  if *getHelp {
    flag.PrintDefaults()
    os.Exit(0)
  }

  // Get configuration settings
  config := configSettings(*configFile)

  if *maxThreads > 0 {
    config.MaxThreads = *maxThreads
  } else if config.MaxThreads <= 0 {
    config.MaxThreads = defaultMaxThreads
  }

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
    if subreddit.CustomFolderName != "" {
      outputPath = outputPath + subreddit.CustomFolderName + "/"
    } else {
      outputPath = outputPath + subreddit.Name[3:] + "/"
    }
    err := os.MkdirAll(outputPath, 0755)
    util.Check(err)

    // Download to folder
    // Concurrent Go channels
    ch := make(chan string)
    startTime := time.Now()

    postsPerThread := len(posts)/config.MaxThreads
    currentStart, currentEnd := 0, 0
    remainder := len(posts)%config.MaxThreads
    for i := 0; i < config.MaxThreads; i++ {
      currentEnd = currentStart + postsPerThread
      if remainder > 0 {
        currentEnd++
        remainder--
      }
      fmt.Println("Starting goroutine from index", currentStart, "to index", currentEnd - 1)
      go downloadToFolder(outputPath, posts[currentStart: currentEnd], ch)
      currentStart = currentEnd
    }

    for i := 0; i < len(posts); i++ {
      v, ok := <-ch
      if !ok {
        break
      }
      fmt.Println(v)
    }

    endTime := time.Now()
    fmt.Println("Total time taken:", endTime.Sub(startTime))
  }
}


func downloadToFolder(folder string, posts []reddit.DownloadPost, ch chan string) {
  fmt.Println("Go routine to download", len(posts), "posts")
  for _, post := range posts {
    outputFile := folder + post.Name
    //fmt.Println("Downloading", post.Url, "to", outputFile)
    err := http.DownloadFile(outputFile, post.Url)
    util.CheckWarn(err)
    ch <- "Downloaded " + post.Url
  }
  close(ch)
}


func configSettings(filename string) Configuration {
  // Open config file
  configuration := Configuration{}
  file, err := os.Open(filename)

  // File does not exist, create it
  if os.IsNotExist(err) {
    // Create file
    fmt.Println("No configuration file found at", filename)
    fmt.Println("Initiating new configuration file...")
    file, err = os.Create(filename)
    util.Check(err)

    // Setup and encode the JSON
    var b []byte
    configuration.Subreddits = append(configuration.Subreddits, SubredditConfig{"/r/subreddit1", 50, "new", "all", "", ""}, SubredditConfig{"/r/subreddit2", 20, "hot", "all", "", ""})
    configuration.OutputPath = "Path/To/Output/Folder"
    configuration.MaxThreads = defaultMaxThreads
    b, err = json.MarshalIndent(configuration, "", "    ")
    util.Check(err)

    // Write to the new file
    _, err = file.Write(b)
    util.Check(err)

    // Close the fd
    err = file.Close()
    util.Check(err)

    // Exit
    fmt.Println("Please edit", filename)
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
  }

  values.Set("t", "all")
  if subreddit.Time != "" {
    values.Set("t", subreddit.Time)
  }

  redditReq.RawQuery = values.Encode()

  return redditReq.String()
}
