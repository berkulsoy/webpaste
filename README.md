![](https://travis-ci.org/berkulsoy/webpaste.svg?branch=master)
# Web Paste
Simple http server that lets you post any file and later get it via an easy to remember name (provided by [golang-petname](https://github.com/dustinkirkland/golang-petname)).

## How

* Run daemon or docker continer
  ```
  ./webpaste-linux-amd64
  ```

* Send something to via curl. This will return you a unique, easy url to download later from somewhere else
  ```
  curl -F "f=@your_file" http://address:port/
  ```
 
* Use the given link in another computer to download  

### Parameters
* -p : port to bind
* -d : directory to save files


## Why ?
1. I wanted to experiment with Go (Obvious from the code)

1. During daily work, we fetch text/binary files from servers in order to later paste them as evidences in tickets, emails or to use them somewhere else

   Even though rsync/scp is de facto for such purposes, sometimes they are not the most practical/easy to use. Because:
   
     * Sometimes people have jump over many servers, with different credentials, causing them to traverse all hops for single file     
     * Sometimes people access servers via Windows terminal sessions, vnc or even via conferencing tools     
     * Sometimes, the easiest thing to access from multiple locations in a complex network is just an http server

## Build Requirements
* Go
* Make
* Docker (If wanted)  

## Build Steps
* make

## Notes
* It does not clean uploaded files
