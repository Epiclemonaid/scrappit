package util

import (
  "log"
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
