package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/codeskyblue/go-sh"
)

// NEED TO INSTALL
// go run github.com/ChimeraCoder/anaconda
// go run github.com/codeskyblue/go-sh

// TwitterCreds : struct for Twitter credentials
type TwitterCreds struct {
	CustomerKey       string
	CustomerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// Periscope : Struct of json, meant for extracting the "LhlsURL" variable
type Periscope struct {
	LhlsURL string `json:"lhls_url"`
}

func main() {
	// Twitter crendentials
	twitterCreds := TwitterCreds{
		CustomerKey:       os.Getenv("TWITTER_CUSTOMERKEY"),
		CustomerSecret:    os.Getenv("TWITTER_CUSTOMERSECRET"),
		AccessToken:       os.Getenv("TWITTER_ACCESSTOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESSTOKENSECRET"),
	}

	// creds
	anaconda.SetConsumerKey(twitterCreds.CustomerKey)
	anaconda.SetConsumerSecret(twitterCreds.CustomerSecret)

	// Twitter API
	api := anaconda.NewTwitterApi(twitterCreds.AccessToken, twitterCreds.AccessTokenSecret)

	// "periscope" results
	search_result, err := api.GetSearch("periscope", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("metadata count, search result status count ---- ", search_result.Metadata.Count, len(search_result.Statuses))
	fmt.Println()

	// Find links that contain the substring "https://www.pscp.tv/w/"
	var selectedfiles []string
	for _, tweet := range search_result.Statuses {
		for _, nono := range tweet.Entities.Urls {
			if strings.Contains(nono.Expanded_url, "https://www.pscp.tv/w/") {
				selectedfiles = append(selectedfiles, nono.Expanded_url)
			}
		}
	}

	// Split strings and add the broadcast_id to the new url
	var newURLS []string
	for i := 0; i < len(selectedfiles); i++ {
		mo := strings.SplitAfter(selectedfiles[i], "https://www.pscp.tv/w/")
		newURLS = append(newURLS, "https://api.periscope.tv/api/v2/accessVideoPublic?broadcast_id="+mo[1])
	}

	// Get periscope videos
	var collectedHLS []string
	for ko := 0; ko < len(newURLS); ko++ {
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

		if record.LhlsURL != "" {
			// fmt.Println(record.LhlsURL)
			fmt.Println("---------------stable link")
			fmt.Println("lhls_url : ", record.LhlsURL)
			fmt.Println("--------------------------")
			collectedHLS = append(collectedHLS, record.LhlsURL)
		} else {

		}
	}

	// rip content for each HLS that was collected
	for mm := 0; mm < len(collectedHLS); mm++ {
		newWavFile := "file" + strconv.Itoa(mm) + ".wav"
		fmt.Println(collectedHLS[mm])
		fmt.Println()
		RipPeriscopeContent(10, collectedHLS[mm], newWavFile)
	}

}

// RipPeriscopeContent : Maintains the ffmpeg script for extracting .wav files from the Periscope stream
func RipPeriscopeContent(durationOfRip int, periscopeLink string, outputFile string) {
	sh.Command("ffmpeg", "-ss", "0", "-t", "10", "-i", periscopeLink, "-strict", "-2", "-ac", "1", outputFile, "-nostdin", "-nostats").Run()
	// ffmpeg -ss 0 -t 10 -i https://prod-fastly-ap-southeast-1.video.periscope.tv/Transcoding/v1/lhls/aksUJRBAr4ZqFN2o0hwNU-zS_GSeSmQw-Rr1SGN-hVvuhJ56ms1BgGR_IOZWyO9ZgrgMMCs7G9vnMD88Pix5-g/non_transcode/ap-southeast-1/periscope-replay-direct-prod-ap-southeast-1-public/dynamic_lowlatency.m3u8 -strict -2 -ac 1 /Users/io/Desktop/tester.wav -nostdin -nostats
}
