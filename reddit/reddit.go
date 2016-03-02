package reddit

import (
  "reddit-scraper/http"
  "reddit-scraper/gfycat"
  "regexp"
  "strings"
)


type Post struct {
  Data struct {
    Domain string `json:"domain"`
    Subreddit string `json:"subreddit"`
    Title string `json:"title"`
    Permalink string `json:"permalink"`
    Url string `json:"url"`
  } `json:"data"`
}


type ListingJson struct {
  Data struct {
    Children []Post `json:"children"`
  } `json:"data"`
}


func DownloadPosts(posts []Post) []string {
  var urls []string
  for _, post := range posts {
    domain := post.Data.Domain
    switch {
    case strings.Contains(domain, "gfycat"):
      urls = append(urls, gfycat.GetDownloadUrl(post.Data.Url))
    }
  }
  return urls
}


func GetRedditUrls(urls []http.URL) []http.URL {
  return http.Filter(urls, func(u http.URL) bool {
    match, err := regexp.MatchString(`^(\w+\.)?reddit\.com$`, u.Host)
    if err != nil {
      return false
    }
    return match
  })
}


func GetSubreddits(urls []http.URL) []http.URL {
  return http.Filter(urls, func(u http.URL) bool {
    return strings.Contains(u.Host, "reddit")
  })
}
