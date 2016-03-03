package reddit

import (
  "reddit-scraper/http"
  "reddit-scraper/gfycat"
  "reddit-scraper/util"
  "net/url"
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
      // Parse the gfycat URL
      gfyUrl, err := url.Parse(post.Data.Url)
      util.Check(err)

      // Generate a raw URL
      rawUrl := gfyUrl.Scheme + "://" + gfyUrl.Host + "/" + gfyUrl.Path

      // Add to URL list
      urls = append(urls, gfycat.GetDownloadUrl(rawUrl))
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
