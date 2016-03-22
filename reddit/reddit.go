package reddit

import (
  "fmt"
  "reddit-scraper/http"
  "reddit-scraper/gfycat"
  "reddit-scraper/imgur"
  "reddit-scraper/util"
  "path/filepath"
  "net/url"
  "regexp"
  "strings"
)


/******************************************
 *                                        *
 * Struct definitions                     *
 *                                        *
 ******************************************/

type Post struct {
  Data struct {
    Domain string `json:"domain"`
    Name string `json:"name"`
    Permalink string `json:"permalink"`
    Score int `json:"score"`
    Stickied bool `json:"stickied"`
    Subreddit string `json:"subreddit"`
    Title string `json:"title"`
    Url string `json:"url"`
  } `json:"data"`
}

type ListingJson struct {
  Data struct {
    Children []Post `json:"children"`
  } `json:"data"`
}

type DownloadPost struct {
  FileType string
  Id string
  Title string
  Score int
  Stickied bool
  Subreddit string
  Url string
}


/******************************************
 *                                        *
 * Exported functions                     *
 *                                        *
 ******************************************/

/*
 *  Parses the Post data
 *  Retrieves the correct download URL depending on host
 */
func GetDownloadPost(post Post) DownloadPost {
  // Parse the URL
  newUrl, err := url.Parse(post.Data.Url)
  util.Check(err)

  // Generate the new post data
  newPost := DownloadPost{}
  newPost.Id = post.Data.Name
  newPost.Score = post.Data.Score
  newPost.Stickied = post.Data.Stickied
  newPost.Subreddit = post.Data.Subreddit
  newPost.Title = util.ReplaceSlashes(post.Data.Title)
  newPost.Url = post.Data.Url

  // Regex
  staticRegex, _ := regexp.Compile(`\.(jpeg|jpg|gif|webm|png)$`)

  fmt.Println(post.Data.Url)

  // Find the URL type
  switch {
  case strings.Contains(post.Data.Domain, "imgur"):
    // Imgur
    newPost.Url = imgur.GetDownloadUrl(post.Data.Url)
    newPost.FileType = filepath.Ext(newPost.Url)

  case staticRegex.MatchString(newUrl.Path):
    // Static file
    newPost.Url = post.Data.Url
    newPost.FileType =  staticRegex.FindString(newUrl.Path)

  case strings.Contains(post.Data.Domain, "gfycat"):
    // Gfycat
    newPost.FileType = ".webm"
    rawUrl := newUrl.Scheme + "://" + newUrl.Host + "/" + newUrl.Path
    newPost.Url = gfycat.GetDownloadUrl(rawUrl)

  default:
    fmt.Println("Unsupported URL:", post.Data.Url)
  }

  return newPost
}


/*
 *  UNUSED
 *  Filters a list of URLs for URLs with domains "reddit.com"
 */
func GetRedditUrls(urls []http.URL) []http.URL {
  return http.Filter(urls, func(u http.URL) bool {
    match, err := regexp.MatchString(`^(\w+\.)?reddit\.com$`, u.Host)
    if err != nil {
      return false
    }
    return match
  })
}
