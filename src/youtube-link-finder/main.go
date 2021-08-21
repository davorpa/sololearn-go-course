/*
YouTube Link Finder +50 XP

You and your friends like to share YouTube links all throughout the day. You want to keep track of all the videos you watch in your own personal notepad, but you find that keeping the entire link is unnecessary. 
Keep the video ID (the combination of letters and numbers at the end of the link) in your notepad to slim down the URL.

Task: 
Create a program that parses through a link, extracts and outputs the YouTube video ID.

Input Format: 
A string containing the URL to a YouTube video. The format of the string can be in "https://www.youtube.com/watch?v=kbxkq_w51PM" or the shortened "https://youtu.be/KMBBjzp5hdc" format.

Output Format: 
A string containing the extracted YouTube video id.

Sample Input: 
https://www.youtube.com/watch?v=RRW2aUSw5vU

Sample Output: 
RRW2aUSw5vU

Note that the input can be in two different formats.


@author davorpatech
@since  2021-08-21
*/
package main

import (
    "bufio"
    "fmt"
    "log"
    "net/url"
    "os"
    "regexp"
    "strings"
    )

func init() {
    log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
    log.SetOutput(os.Stdout) //os.Stderr
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanLines)
    
    fmt.Print("Enter video url: ")
    var input string
    if scanner.Scan() {
        input = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprint(os.Stderr, "ERR! ", err)
        os.Exit(1)
    }
    
    vid, err := ExtractYoutubeVideoId(input)
    if err != nil {
        fmt.Fprint(os.Stderr, "ERR! ", err)
        os.Exit(1)
    }
    fmt.Printf("Video Id: %s\n", vid)
}



func ExtractYoutubeVideoId(s string) (string, error) {
    var vid string
    
    u, err := url.Parse(s)
    if err != nil {
        return vid, fmt.Errorf("\"%s\": unparseable URL", s)
    }
    if !u.IsAbs() {
        return vid, fmt.Errorf("\"%s\": non absolute URL", s)
    }
    
    switch host := u.Hostname(); true {
        
        // normal
        case matchre("(?i)(www\\.)?youtube\\.(com)", host):
            qs, err := url.ParseQuery(u.RawQuery)
            if err != nil {
                return vid, err
            }
            // get first "v" qs param
            vid = qs.Get("v")
        
        // shorten
        case matchre("(?i)youtu\\.(be)", host):
            // remove slash from path, if any
            vid = strings.TrimLeft(
                u.EscapedPath(),
                "/")
        
        // unrecognized
        default:
            err = fmt.Errorf("\"%s\"; invalid YouTube url", s)
    }
    
    // ckeck missing
    if err == nil && len(vid) == 0 {
        err = fmt.Errorf("\"%s\": video id not found", s)
    }
    return vid, err
}


func matchre(re, s string) bool {
    // @TODO: cache compile??
    return regexp.MustCompile(
        re).MatchString(
            s)
}
