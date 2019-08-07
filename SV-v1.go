package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	Session                     string `json:"session"`
	HlsURL                      string `json:"hls_url"`
	LhlsURL                     string `json:"lhls_url"`
	LhlswebURL                  string `json:"lhlsweb_url"`
	HTTPSHlsURL                 string `json:"https_hls_url"`
	HlsIsEncrypted              bool   `json:"hls_is_encrypted"`
	LhlsIsEncrypted             bool   `json:"lhls_is_encrypted"`
	Type                        string `json:"type"`
	MediaConfiguration          string `json:"media_configuration"`
	DefaultPlaybackBufferLength int    `json:"default_playback_buffer_length"`
	MinPlaybackBufferLength     int    `json:"min_playback_buffer_length"`
	MaxPlaybackBufferLength     int    `json:"max_playback_buffer_length"`
	ChatToken                   string `json:"chat_token"`
	LifeCycleToken              string `json:"life_cycle_token"`
	Broadcast                   struct {
		ClassName          string    `json:"class_name"`
		ID                 string    `json:"id"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		UserID             string    `json:"user_id"`
		UserDisplayName    string    `json:"user_display_name"`
		Username           string    `json:"username"`
		TwitterID          string    `json:"twitter_id"`
		ProfileImageURL    string    `json:"profile_image_url"`
		State              string    `json:"state"`
		IsLocked           bool      `json:"is_locked"`
		FriendChat         bool      `json:"friend_chat"`
		PrivateChat        bool      `json:"private_chat"`
		Language           string    `json:"language"`
		Version            int       `json:"version"`
		Start              time.Time `json:"start"`
		Ping               time.Time `json:"ping"`
		HasModeration      bool      `json:"has_moderation"`
		EnabledSparkles    bool      `json:"enabled_sparkles"`
		HasLocation        bool      `json:"has_location"`
		City               string    `json:"city"`
		Country            string    `json:"country"`
		CountryState       string    `json:"country_state"`
		IsoCode            string    `json:"iso_code"`
		IPLat              float64   `json:"ip_lat"`
		IPLng              float64   `json:"ip_lng"`
		Width              int       `json:"width"`
		Height             int       `json:"height"`
		CameraRotation     int       `json:"camera_rotation"`
		ImageURL           string    `json:"image_url"`
		ImageURLSmall      string    `json:"image_url_small"`
		Status             string    `json:"status"`
		ContentType        string    `json:"content_type"`
		BroadcastSource    string    `json:"broadcast_source"`
		AvailableForReplay bool      `json:"available_for_replay"`
		Expiration         int       `json:"expiration"`
		TweetID            string    `json:"tweet_id"`
		MediaKey           string    `json:"media_key"`
	} `json:"broadcast"`
	ShareURL              string `json:"share_url"`
	AutoplayViewThreshold int    `json:"autoplay_view_threshold"`
}

func main() {
	// Twitter crendentials
	twitterCreds := TwitterCreds{
		CustomerKey:       os.Getenv("TWITTER_CUSTOMERKEY"),
		CustomerSecret:    os.Getenv("TWITTER_CUSTOMERSECRET"),
		AccessToken:       os.Getenv("TWITTER_ACCESSTOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESSTOKENSECRET"),
	}

	// scoped variables
	var selectedfiles, collectedHLS []string

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
			// fmt.Println("got ittt")
		}
		// CreateVlcLinks(record.LhlsURL)
		// RipPeriscopeContent(10, record.LhlsURL, "test.wav")
	}

	// rip content for each HLS that was collected
	for mm := 0; mm < len(collectedHLS); mm++ {
		newWavFile := "file" + strconv.Itoa(mm) + ".wav"
		fmt.Println(collectedHLS[mm])
		fmt.Println()
		RipPeriscopeContent(10, collectedHLS[mm], newWavFile)
	}

}

// CreateVlcLinks : Generate terminal app links for VLC for the Periscope content
func CreateVlcLinks(vidlink string) {
	appLinks := "/Applications/VLC.app/Contents/MacOS/VLC -vvv " + vidlink
	fmt.Println(appLinks)
	fmt.Println()
}

// RipPeriscopeContent : Maintains the ffmpeg script for extracting .wav files from the Periscope stream
func RipPeriscopeContent(durationOfRip int, periscopeLink string, outputFile string) {
	sh.Command("ffmpeg", "-ss", "0", "-t", "10", "-i", periscopeLink, "-strict", "-2", "-ac", "1", outputFile, "-nostdin", "-nostats").Run()
	// ffmpeg -ss 0 -t 10 -i https://prod-fastly-ap-southeast-1.video.periscope.tv/Transcoding/v1/lhls/aksUJRBAr4ZqFN2o0hwNU-zS_GSeSmQw-Rr1SGN-hVvuhJ56ms1BgGR_IOZWyO9ZgrgMMCs7G9vnMD88Pix5-g/non_transcode/ap-southeast-1/periscope-replay-direct-prod-ap-southeast-1-public/dynamic_lowlatency.m3u8 -strict -2 -ac 1 /Users/io/Desktop/tester.wav -nostdin -nostats
}
