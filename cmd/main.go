package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/codeskyblue/go-sh"
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

	// Twitter crendentials
	twitterCreds := TwitterCreds{
		CustomerKey:       os.Getenv("TWITTER_CUSTOMERKEY"),
		CustomerSecret:    os.Getenv("TWITTER_CUSTOMERSECRET"),
		AccessToken:       os.Getenv("TWITTER_ACCESSTOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESSTOKENSECRET"),
	}

	// Apply credentials to anaconda creds
	anaconda.SetConsumerKey(twitterCreds.CustomerKey)
	anaconda.SetConsumerSecret(twitterCreds.CustomerSecret)

	// Twitter API
	api := anaconda.NewTwitterApi(twitterCreds.AccessToken, twitterCreds.AccessTokenSecret)

	// Search Twitter for "periscope" results
	v := url.Values{}
	v.Set("count", "250")
	searchResult, err := api.GetSearch("periscope", v)
	if err != nil {
		panic(err)
	}

	// Find links that contain the substring "https://www.pscp.tv/w/"
	var selectedfiles []string
	for _, tweet := range searchResult.Statuses {
		for _, nono := range tweet.Entities.Urls {
			if strings.Contains(nono.Expanded_url, "https://www.pscp.tv/w/") {
				selectedfiles = append(selectedfiles, nono.Expanded_url)
			}
		}
	}

	// Split strings and add the broadcast_id to the new url
	var newURLS []string
	for i := 0; i < len(selectedfiles); i++ {
		// fmt.Println(selectedfiles[i])
		mo := strings.SplitAfter(selectedfiles[i], "https://www.pscp.tv/w/")
		newURLS = append(newURLS, "https://api.periscope.tv/api/v2/accessVideoPublic?broadcast_id="+mo[1])
	}

	// Get periscope videos
	var collectedHLS []string
	for ko := 0; ko < len(newURLS); ko++ {

		// If the URL contains a '?t'
		// https://api.periscope.tv/api/v2/accessVideoPublic?broadcast_id=cOqTLjFZTEVKeU1HVndqTm58MUJkR1llcldNcUFHWGaZiJ145fdQMGpchLOQbQDkUdFw1znXOLwDG0vT7_MF?t=22s
		switch strings.Contains(newURLS[ko], "?t") {
		case true:
			// Split at the time
			splitter := strings.SplitAfter(newURLS[ko], "?t")
			moa := splitter[0][:len(splitter[0])-2]

			// GET request
			resp, err := http.Get(moa)
			if err != nil {
			}
			defer resp.Body.Close()
			body, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				fmt.Println("FAIL")
			}

			var record Periscope
			json.Unmarshal(body, &record)

			// If the URL is a replay
			switch strings.Contains(record.ReplayURL, "?type=replay") {
			case true: // REPLAY
				// record.LhlsURL will return empty
				// IF THE RECORD IS NOT EMPTY
				if record.ReplayURL != "" {
					replay := record.ReplayURL[:len(record.ReplayURL)-12]
					collectedHLS = append(collectedHLS, replay)
				}
			case false:
				// record.ReplayURL will return empty
				// Collect LhlsURL links
				if record.LhlsURL != "" {
					collectedHLS = append(collectedHLS, record.LhlsURL)
				}
			}
		// If the URL DOES NOT CONTAIN a '?t'
		case false:
			// GET request
			resp, err := http.Get(newURLS[ko])
			if err != nil {
			}

			defer resp.Body.Close()
			body, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				fmt.Println("FAIL")
			}

			var record Periscope
			json.Unmarshal(body, &record)

			switch strings.Contains(record.ReplayURL, "?type=replay") {
			case true:
				if record.ReplayURL != "" {
					// REPLAY
					replay := record.ReplayURL[:len(record.ReplayURL)-12]
					collectedHLS = append(collectedHLS, replay)
				}
			case false:
				// dynamic_lowlatency.m3u8
				if record.LhlsURL != "" {
					collectedHLS = append(collectedHLS, record.LhlsURL)
				}
			}
		}
	}

	// Output directory
	dir := os.Args[1]

	// Remove duplicates from the collection
	cleanDir := RemoveDuplicates(collectedHLS)

	// rip content for each HLS that was collected
	for mm := 0; mm < len(cleanDir); mm++ {
		fmt.Println(cleanDir[mm])
		newWavFile := dir + "/file" + strconv.Itoa(mm) + ".wav"
		// fmt.Println(newWavFile)
		RipPeriscopeContent(10, cleanDir[mm], newWavFile)
	}
}

// RipPeriscopeContent : Maintains the ffmpeg script for extracting .wav files from the Periscope stream
func RipPeriscopeContent(durationOfRip int, periscopeLink string, outputFile string) {
	sh.Command("ffmpeg", "-ss", "0", "-t", "10", "-i", periscopeLink, "-strict", "-2", "-ac", "1", outputFile, "-nostdin", "-nostats").Run()
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}
