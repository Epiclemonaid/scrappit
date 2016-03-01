package main

import (
  "fmt"
  "reddit-scraper/gfycat"
  "reddit-scraper/http"
)

func main() {
  fmt.Printf("Scraper v1.0\n")

  item := gfycat.GfyJson{}
  http.GetJson("https://gfycat.com/cajax/get/ImpoliteImmenseCaracal", &item)
  fmt.Printf("Results: %s\n", item.GfyItem.GifUrl)

  http.DownloadFile("file", item.GfyItem.GifUrl)
}

