package util

import (
  "log"
  "strings"
)


/******************************************
 *                                        *
 * Exported functions                     *
 *                                        *
 ******************************************/

/*
 *  Print error message and exit, if there is one
 */
func Check(err error) {
  if err != nil {
    log.Panic("Error:", err)
  }
}


/*
 *  Print error message, if there is one
 */
func CheckWarn(err error) {
  if err != nil {
    log.Println("Error:", err)
  }
}


/*
 *  Replaces all slashes in a string with a slash-like character
 *  Used for filenames with slashes
 */
func ReplaceSlashes(path string) string {
  return strings.Replace(path, "/", "\u2215", -1)
}
