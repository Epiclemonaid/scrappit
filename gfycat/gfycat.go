package gfycat

type GfyItem struct {
  GfyId string
  GfyName string
  GfyNumber string
  UserName string
  Wdith string
  Height string
  FrameRate string
  NumFrames string
  Mp4Url string
  WebmUrl string
  WebpUrl string
  MobileUrl string
  MobilePosterUrl string
  PosterUrl string
  Thumb360Url string
  Thumb360PosterUrl string
  Thumb100PosterUrl string
  Max5mbGif string
  Max2mbGif string
  MjpgUrl string
  GifUrl string
  GifSize string
  Mp4Size string
  WebmSize string
  CreateDate string
  Views int
  Title string
  ExtraLemmas []string
  Md5 string
  Tags []string
  Nsfw string
  Sar string
  Url string
  Source string
  Dynamo string
  Subreddit string
  RedditId string
  RedditIdText string
  Likes string
  Dislikes string
  Published string
  Description string
  CopyrightClaimant string
}

type GfyJson struct {
  GfyItem GfyItem
  Error string
}
