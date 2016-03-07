# Reddit Scraper

* Created by Curtis Li
* A scraper for www.reddit.com, which downloads posts, images and gifs for offline viewing.
* Written in Go.
* Uses Reddit API, Gfycat API

## Installation

#### Building from Source
1. Clone the repository into your $GOPATH/src directory
2. Run 'go get golang.org/x/net/http'
3. Run 'go install'
4. Run 'reddit-scraper' to generate a conf.json template
5. Fill out template and run 'reddit-scraper' again

## Features

* Scrapes given subreddits for links
* Supports reddit features such as: searching, post limiting, time limiting, subcategorizing (hot, new, top)
* Automatically downloads the content into an output directory
* Current supported post types: gfycat links, static files (.png, .jpg, etc)
* Concurrent downloading

## TODO

* Scheduled running
* Keep store of last item downloaded (i.e. Get all posts up to some other post)
* Imgur posts
* Self.Reddit posts
