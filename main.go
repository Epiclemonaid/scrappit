package main

import (
  "fmt"
  "os"
  "regexp"
  "strconv"
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
  OutputPath string `json:"outputPath"`
  Stats bool `json:"stats"`
}


func main() {
  fmt.Println("Scraper v0.1")
  fmt.Println("Created by Curtis Li")

  // Get configuration settings
  config := configSettings(configFile)

  // Loop through subreddit list
  for _, subreddit := range config.Subreddits {
    fmt.Println("-----------------------------\n")
    fmt.Println(subreddit)

    // Get subreddit JSON
    listing := reddit.ListingJson{}
    http.GetJson("http://www.reddit.com/" + subreddit + "/.json", &listing)

    // Get download links
    downloadUrls := reddit.DownloadPosts(listing.Data.Children)
    fmt.Println(len(downloadUrls), "posts to download")

    // Get output directory path
    outputPath := config.OutputPath
    if outputPath == "" {
      outputPath = "output/"
    }
    r, _ := regexp.Compile(`.*/$`)
    if !r.MatchString(outputPath) {
      outputPath = outputPath + "/"
    }

    // Download to folder
    for index, downloadUrl := range downloadUrls {
      outputFile := outputPath + strconv.Itoa(index)
      fmt.Println("Downloading", downloadUrl, "to", outputFile)
      err := http.DownloadFile(outputFile, downloadUrl)
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
    configuration.Username = "Username"
    configuration.Subreddits = append(configuration.Subreddits, "/r/subreddit1", "/r/subreddit2", "/r/subreddit3")
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

