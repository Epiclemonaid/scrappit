package gfycat

import (
  "reddit-scraper/http"
  "regexp"
  "strings"
  "fmt"
)


type GfyJson struct {
  GfyItem struct {
    GfyName string `json:"gfyName"`
    WebmUrl string `json:"webmUrl"`
    WebmSize string `json:"webmSize"`
  } `json:"gfyItem"`
  Error string `json:"error"`
}


func GetAjaxUrl(url string) string {
  r, _ := regexp.Compile(`\/\w+$`)
  id := r.FindString(url)
  return "https://gfycat.com/cajax/get" + id + ".json"
}


func GetDownloadUrl(url string) string {
  ajaxUrl := url
  if !strings.Contains(url, "cajax") {
    ajaxUrl = GetAjaxUrl(url)
  }
  item := GfyJson{}
  err := http.GetJson(ajaxUrl, &item)
  if err != nil {
    fmt.Println(err)
  }
  return item.GfyItem.WebmUrl
}
