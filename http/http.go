package http

import (
  "os"
  "io"
  "io/ioutil"
  "strings"
  "regexp"
  "net/http"
  "net/url"
  "encoding/json"
  "golang.org/x/net/html"
)

type URL url.URL


/******************************************
 *                                        *
 * Exported functions                     *
 *                                        *
 ******************************************/

/*
 *  Filters a list of URLs for URLs with a matching domain
 */
func GetSiteUrls(urls []URL, site string) []URL {
  return Filter(urls, func(u URL) bool {
    match, err := regexp.MatchString(site, u.Host)
    if err != nil {
      return false
    }
    return match
  })
}


/*
 *  UNUSED
 *  Crawls a page to find all the links
 */
func Crawl(u string) ([]URL, error) {
  var urls []URL

  r, err := http.Get(u)
  if err != nil {
    return nil, err
  }
  defer r.Body.Close()

  tokenizer := html.NewTokenizer(r.Body)

  for {
    token := tokenizer.Next()

    switch {
      case token == html.ErrorToken:
        // End of the document, we're done
        return urls, nil
      case token == html.StartTagToken:
        next := tokenizer.Token()
        isAnchor := next.Data == "a"
        if !isAnchor {
          continue
        }

        found, u := getHref(next)
        if !found {
          continue
        }

        hasProtocol := strings.Index(u, "http") == 0
        if hasProtocol {
          uObj, err := url.Parse(u)
          if err != nil {
            return nil, err
          }
          urls = append(urls, URL(*uObj))
        }
    }
  }

  return urls, nil
}


/*
 *  UNUSED
 *  Get the 'href' attribute of an link <a> tag
 */
func getHref(token html.Token) (found bool, href string) {
  for _, attr := range token.Attr {
    if attr.Key == "href" {
      href = attr.Val
      found = true
    }
  }

  return
}


/*
 *  Filters an array of URLs using a given function
 *  Calls the function for each element in the array
 *  The element is kept if the function returns true
 */
func Filter(origUrls []URL, f func(URL) bool) []URL {
  newUrls := make([]URL, 0)
  for _, ele := range origUrls {
    if f(ele) {
      newUrls = append(newUrls, ele)
    }
  }
  return newUrls
}


/*
 *  UNUSED
 *  Get the body of a page given its URL
 */
func GetBody(u string) ([]byte, error) {
  r, err := http.Get(u)
  if err != nil {
    return nil, err
  }

  bytes, err := ioutil.ReadAll(r.Body)
  if err != nil {
    return nil, err
  }

  r.Body.Close()
  return bytes, nil
}


/*
 *  Parse a JSON object from a URL
 *  Places the values into a provided struct
 */
func GetJson(u string, target interface{}) error {
  // Get the data
  r, err := http.Get(u)
  if err != nil {
    return err
  }
  defer r.Body.Close()

  // Decode the JSON
  return json.NewDecoder(r.Body).Decode(target)
}


/*
 *  Downloads a file from a given URL
 *  The new file location is at 'filepath'
 */
func DownloadFile(filepath string, u string) error {
  // Create the file
  out, err := os.Create(filepath)
  if err != nil {
    return err
  }
  defer out.Close()

  // Get the data
  r, err := http.Get(u)
  if err != nil {
    return err
  }
  defer r.Body.Close()

  // Write the body to a file
  _, err = io.Copy(out, r.Body)
  if err != nil {
    return err
  }

  return nil
}
