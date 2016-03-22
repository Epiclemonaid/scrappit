package imgur

import (
  "reddit-scraper/http"
  "reddit-scraper/util"
  "path/filepath"
  "regexp"
  "strings"
)


/******************************************
 *                                        *
 * Struct definitions                     *
 *                                        *
 ******************************************/

type AlbumData struct {
  Id string `json:"id"`
  Title string `json:"title"`
  Description string `json:"description"`
  Link string `json:"link"`
  ImagesCount int `json:"images_count"`
  Images []ImageData `json:"images"`
}

type ImageData struct {
  Id string `json:"id"`
  Title string `json:"title"`
  Description string `json:"description"`
  Type string `json:"type"`
  Animated bool `json:"animated"`
  Size int `json:"size"`
  Link string `json:"link"`
  Webm string `json:"webm"`
}

type ImgurJson struct {
  Data ImageData `json:"data"`
  Success bool `json:"success"`
}

var baseUrl = "https://api.imgur.com/3/";

/******************************************
 *                                        *
 * Exported functions                     *
 *                                        *
 ******************************************/

/*
 *  Get the JSON request URL from the public URL
 */
func GetAjaxUrl(url string) string {
  r, _ := regexp.Compile(`\/[a-zA-Z0-9.?#]+$`)
  id := r.FindString(url)
  id = strings.TrimSuffix(id, filepath.Ext(id))
  return baseUrl + "image" + id
}


/*
 *  Get the download URL from the public URL
 */
func GetDownloadUrl(url string) string {
  ajaxUrl := url
  if !strings.Contains(url, baseUrl) {
    ajaxUrl = GetAjaxUrl(url)
  }
  item := ImgurJson{}

  headers := make(map[string]string)
  headers["Authorization"] = "Client-ID b149599c942745b"

  err := http.GetJson(ajaxUrl, &item, headers)
  util.CheckWarn(err)

  if item.Data.Animated {
    return item.Data.Webm
  } else {
    return item.Data.Link
  }
}
