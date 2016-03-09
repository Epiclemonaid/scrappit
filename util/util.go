package util

import (
  "log"
  "strings"
)

func Check(err error) {
  if err != nil {
    log.Panic("Error:", err)
  }
}

func CheckWarn(err error) {
  if err != nil {
    log.Println("Error:", err)
  }
}

func ReplaceSlashes(path string) string {
  return strings.Replace(path, "/", "\u2215", -1)
}
