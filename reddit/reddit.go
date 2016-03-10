package reddit

import (
  "fmt"
  "reddit-scraper/http"
  "reddit-scraper/gfycat"
  "reddit-scraper/util"
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

type DownloadPost struct {
  Name string
  Id string
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
func DownloadPosts(posts []Post) []DownloadPost {
  var newPosts []DownloadPost
  for _, post := range posts {
    // Parse the URL
    newUrl, err := url.Parse(post.Data.Url)
    util.Check(err)

    // Generate the new post data
    newPost := DownloadPost{}
    newPost.Subreddit = post.Data.Subreddit
    newPost.Name = util.ReplaceSlashes(post.Data.Title)
    newPost.Url = post.Data.Url

    // Regex
    staticRegex, _ := regexp.Compile(`\.(jpeg|jpg|gif|webm|png)$`)

    // Find the URL type
    switch {
    case staticRegex.MatchString(newUrl.Path):
      // Static file
      newPost.Url = post.Data.Url
      fileType :=  staticRegex.FindString(newUrl.Path)

      newPost.Name = newPost.Name + fileType
    case strings.Contains(post.Data.Domain, "gfycat"):
      // Gfycat
      newPost.Name = newPost.Name + ".webm"
      rawUrl := newUrl.Scheme + "://" + newUrl.Host + "/" + newUrl.Path
      newPost.Url = gfycat.GetDownloadUrl(rawUrl)

    case strings.Contains(post.Data.Domain, "imgur"):
      // Non-static Imgur
      fmt.Println("Imgur:", post.Data.Url)

    default:
      fmt.Println("Unsupported URL:", post.Data.Url)
    }

    // Add to URL list
    newPosts = append(newPosts, newPost)
  }
  return newPosts
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
