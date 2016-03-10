package main

import (
  "encoding/json"
  "fmt"
  "flag"
  "net/url"
  "os"
  "reddit-scraper/util"
  "reddit-scraper/http"
  "reddit-scraper/reddit"
  "regexp"
  "strconv"
  "sync"
  "time"
)


/******************************************
 *                                        *
 * Global variables and structs           *
 *                                        *
 ******************************************/

const defaultConfigFile = "conf.json"
const defaultMaxThreads = 10
const redditUrl = "http://www.reddit.com"

var config Configuration

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
  MinScore int `json:"minScore"`
  SearchFor string `json:"searchFor"`
  CustomFolderName string `json:"customFolderName"`
}


/******************************************
 *                                        *
 * Main function                          *
 *                                        *
 ******************************************/

func main() {
  // Set up the program (config file and command-line flags)
  config = setup()

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
    posts := []reddit.Post(listing.Data.Children)
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
    var wg sync.WaitGroup
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
      wg.Add(1)
      go downloadToFolder(outputPath, posts[currentStart: currentEnd], subreddit, &wg)
      currentStart = currentEnd
    }

    // Block and wait for go routines to complete
    wg.Wait()

    endTime := time.Now()
    fmt.Println("Total time taken:", endTime.Sub(startTime))
  }
}


/******************************************
 *                                        *
 * Helper functions                       *
 *                                        *
 ******************************************/

/*
 *  Loads a configuration file to the program
 *  If no configuration file exists, creates a new one
 */
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
    configuration.Subreddits = append(configuration.Subreddits, SubredditConfig{"/r/subreddit1", 50, "new", "all", 0, "", ""}, SubredditConfig{"/r/subreddit2", 20, "hot", "all", 0, "", ""})
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


/*
 *  Creates the JSON request URL for Reddit given a configuration
 */
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


/*
 *  Downloads a file from a URL to the given folder
 *  Go routine thread function
 *  Outputs success messages to main function
 */
func downloadToFolder(folder string, posts []reddit.Post, config SubredditConfig, wg *sync.WaitGroup) {
  fmt.Println("Go routine to download", len(posts), "posts")
  defer wg.Done()

  for _, post := range posts {

    // Get the post data
    downloadPost := reddit.GetDownloadPost(post)

    if downloadPost.Score < config.MinScore && config.MinScore > 0 {
      continue
    }

    // Determine output location
    outputFile := folder + downloadPost.Title + downloadPost.FileType

    // Download the file
    err := http.DownloadFile(outputFile, downloadPost.Url)
    util.CheckWarn(err)

    fmt.Println("Downloaded:", downloadPost.Title, "\n\tID:", downloadPost.Id, "\tScore:", downloadPost.Score)
  }
}


/*
 *  Set up function to be run when program is loading
 *  Handles command-line flags and configurations file
 *  Prints help statements
 */
func setup() Configuration {
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
  config = configSettings(*configFile)

  if *maxThreads > 0 {
    config.MaxThreads = *maxThreads
  } else if config.MaxThreads <= 0 {
    config.MaxThreads = defaultMaxThreads
  }

  return config
}

