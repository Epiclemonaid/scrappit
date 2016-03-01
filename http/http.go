package http

import (
  "os"
  "io"
  "net/http"
  "encoding/json"
)

func GetJson(url string, target interface{}) error {
  // Get the data
  r, err := http.Get(url)
  if err != nil {
    return err
  }
  defer r.Body.Close()

  // Decode the JSON
  return json.NewDecoder(r.Body).Decode(target)
}

func DownloadFile(filepath string, url string) error {
  // Create the file
  out, err := os.Create(filepath)
  if err != nil {
    return err
  }
  defer out.Close()

  // Get the data
  r, err := http.Get(url)
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
