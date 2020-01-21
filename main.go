package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

// TwitterCreds : struct for Twitter credentials
type TwitterCreds struct {
	CustomerKey       string
	CustomerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// Periscope : Struct of json, meant for extracting the "LhlsURL" variable
type Periscope struct {
	LhlsURL   string `json:"lhls_url"`
	ReplayURL string `json:"replay_url"`
}

func main() {
	twitterCreds := TwitterCreds{ // Twitter crendentials
		CustomerKey:       os.Getenv("TWITTER_CUSTOMERKEY"),
		CustomerSecret:    os.Getenv("TWITTER_CUSTOMERSECRET"),
		AccessToken:       os.Getenv("TWITTER_ACCESSTOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESSTOKENSECRET"),
	}
	anaconda.SetConsumerKey(twitterCreds.CustomerKey)                                       // Apply Twitter Consumer Key
	anaconda.SetConsumerSecret(twitterCreds.CustomerSecret)                                 // Apply Twitter Consumer Secret
	api := anaconda.NewTwitterApi(twitterCreds.AccessToken, twitterCreds.AccessTokenSecret) // Twitter API

	v := url.Values{}
	v.Set("count", "500")                              // Set the "count" key with a value
	searchResult, err := api.GetSearch("periscope", v) // Search Twitter for 250 "periscope" results
	if err != nil {
		panic(err)
	}

	var collectedURLs []string
	for _, status := range searchResult.Statuses {
		for _, entities := range status.Entities.Urls {
			switch strings.Contains(entities.Expanded_url, "https://www.pscp.tv/w/") {
			case true:
				broadcastID := strings.SplitAfter(entities.Expanded_url, "https://www.pscp.tv/w/")
				apiURL := "https://api.periscope.tv/api/v2/accessVideoPublic?broadcast_id=" + broadcastID[1]
				switch strings.Contains(apiURL, "?t") {
				case true:
					removeTimeRef := strings.SplitAfter(apiURL, "?t")
					requestURLs := removeTimeRef[0][:len(removeTimeRef[0])-2]
					getStreamURL := acquireStreamURLs(requestURLs)
					collectedURLs = append(collectedURLs, getStreamURL)
				case false:
					getStreamURL := acquireStreamURLs(apiURL)
					collectedURLs = append(collectedURLs, getStreamURL)
				}
			}
		}
	}

	encounteredFiles := map[string]bool{} // incoming files
	nodupsCollectedURLs := []string{}     // duplicates removed
	for v := range collectedURLs {
		if encounteredFiles[collectedURLs[v]] == true {
		} else {
			encounteredFiles[collectedURLs[v]] = true
			nodupsCollectedURLs = append(nodupsCollectedURLs, collectedURLs[v])
		}
	}

	dir := os.Args[1]                                          // User-defined directory to search
	currentTime := time.Now()                                  // Get current time
	formattedTime := currentTime.Format("2006-01-02-15:04:05") // Format time with y-m-d-h-m-s

	for num := 0; num < len(nodupsCollectedURLs); num++ {
		newWavFile := dir + "/" + formattedTime + "-" + strconv.Itoa(num+1) + ".wav"
		cmd := exec.Command("ffmpeg", "-ss", "0", "-t", "10", "-i", nodupsCollectedURLs[num], "-strict", "-2", "-ac", "1", newWavFile, "-nostdin", "-nostats")
		cmd.Run()
	}
}

func acquireStreamURLs(url string) string {
	var outputURL string       // Variable for return string
	resp, err := http.Get(url) // GET request
	if err != nil {
		fmt.Println("FAIL")
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("FAIL")
	}

	var record Periscope
	json.Unmarshal(body, &record)

	switch strings.Contains(record.ReplayURL, "?type=replay") { // Switch between replay and live streams
	case true:
		if record.ReplayURL != "" {
			replayURL := record.ReplayURL[:len(record.ReplayURL)-12]
			outputURL = replayURL
		}
	case false:
		if record.LhlsURL != "" {
			outputURL = record.LhlsURL
		}
	}
	return outputURL
}
