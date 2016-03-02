package gfycat

type GfyItem struct {
  GfyId     string  `json: "gfyId"`
  GfyName   string  `json: "gfyName"`
  GfyNumber string  `json: "gfyNumber"`
  WebmUrl   string  `json: "webmUrl"`
  GifUrl    string  `json: "gifUrl"`
  GifSize   string  `json: "gifSize"`
  WebmSize  string  `json: "webmSize"`
  Title     string  `json: "title"`
  Url       string  `json: "url"`
}

type GfyJson struct {
  GfyItem GfyItem `json: ""`
  Error   string  `json ""`
}
