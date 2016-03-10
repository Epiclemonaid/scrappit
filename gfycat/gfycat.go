package gfycat

import (
  "reddit-scraper/http"
  "reddit-scraper/util"
  "regexp"
  "strings"
)


/******************************************
 *                                        *
 * Struct definitions                     *
 *                                        *
 ******************************************/

type GfyJson struct {
  GfyItem struct {
    GfyName string `json:"gfyName"`
    WebmUrl string `json:"webmUrl"`
    WebmSize string `json:"webmSize"`
  } `json:"gfyItem"`
  Error string `json:"error"`
}


/******************************************
 *                                        *
 * Exported functions                     *
 *                                        *
 ******************************************/

/*
 *  Get the JSON request URL from the public URL
 */
func GetAjaxUrl(url string) string {
  r, _ := regexp.Compile(`\/\w+$`)
  id := r.FindString(url)
  return "https://gfycat.com/cajax/get" + id + ".json"
}


/*
 *  Get the download webm URL from the public URL
 */
func GetDownloadUrl(url string) string {
  ajaxUrl := url
  if !strings.Contains(url, "cajax") {
    ajaxUrl = GetAjaxUrl(url)
  }
  item := GfyJson{}
  err := http.GetJson(ajaxUrl, &item)
  util.CheckWarn(err)
  return item.GfyItem.WebmUrl
}
