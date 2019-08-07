# scopingVoices
#### *scopingVoices* is an artist tool used for extracting audio and video content from the Periscope API in real-time using Go and FFmpeg. The goal behind the project was to create custom utility that can extract audio and video content from around the world for artistic endeavors using [Periscope](https://www.pscp.tv). 

<p align="center">
  <img width="35%" height="35%" src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/023d0f8dfc16d75c30b7409a8bd9883a0fd678b7.png"/>
</p>

## Why? 
#### 

## Requirements

1. Install the following libraries below.
``` go
go run github.com/ChimeraCoder/anaconda
go run github.com/codeskyblue/go-sh
```
2. Generate [Twitter Access Tokens](https://developer.twitter.com/en/docs/basics/authentication/guides/access-tokens.html)

3. Add the Twitter crendentials to your *~/.bash_profile* 
	- Open the Terminal and open *~/.bash_profile* using *nano ~/.bash_profile*
	- Add the Twitter crendentials in the pattern specified below.

	``` bash
	export TWITTER_CUSTOMERKEY="the-twitter-credential"
	export TWITTER_CUSTOMERSECRET="the-twitter-credential"
	export TWITTER_ACCESSTOKEN="the-twitter-credential"
	export TWITTER_ACCESSTOKENSECRET="the-twitter-credential"
	```

	- Close and save the file.
	- Type ```source ~/.bash_profile``` to update the file's contents.

## Run the program
``` go
go run SV-v1.go
```

### How does the program work?

Main structs
``` go
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
```

Twitter credentials using *anaconda*
``` go
// Twitter crendentials object
twitterCreds := TwitterCreds{
	CustomerKey:       os.Getenv("TWITTER_CUSTOMERKEY"),
	CustomerSecret:    os.Getenv("TWITTER_CUSTOMERSECRET"),
	AccessToken:       os.Getenv("TWITTER_ACCESSTOKEN"),
	AccessTokenSecret: os.Getenv("TWITTER_ACCESSTOKENSECRET"),
}

// Credential setup
anaconda.SetConsumerKey(twitterCreds.CustomerKey)
anaconda.SetConsumerSecret(twitterCreds.CustomerSecret)

// Twitter API
api := anaconda.NewTwitterApi(twitterCreds.AccessToken, twitterCreds.AccessTokenSecret)
```

Search the Twitter API for new results using the keyword "periscope".
``` go
// "periscope" results
search_result, err := api.GetSearch("periscope", nil)
if err != nil {
	panic(err)
}
```

Extract the "expanded_url" of each result and find the substring "https://www.pscp.tv/w/". Once the substring is found, append to a string array named "selectedfiles".
``` go
// scoped variables
var selectedfiles []string

// Find links that contain the substring "https://www.pscp.tv/w/"
for _, tweet := range search_result.Statuses {
	for _, nono := range tweet.Entities.Urls {
		if strings.Contains(nono.Expanded_url, "https://www.pscp.tv/w/") {
			selectedfiles = append(selectedfiles, nono.Expanded_url)
		}
	}
}
```

Split the "selectedfiles" string array, extracting the Periscope "broadcast_id" from each element. Add the "broadcast_id" to a new string array named "newURLS".
``` go
// Split strings and add the broadcast_id to the new url
var newURLS []string
for i := 0; i < len(selectedfiles); i++ {
	mo := strings.SplitAfter(selectedfiles[i], "https://www.pscp.tv/w/")
	newURLS = append(newURLS, "https://api.periscope.tv/api/v2/accessVideoPublic?broadcast_id="+mo[1])
}
```

Create a GET Request for each element within the "newURLS" string array. Use the *Periscope* struct to decode the JSON to receive the *lhls_url* of live and replay streams. Add "active" streams to a new string array named "collectedHLS".
``` go
// Get active periscope videos via lhls_url links
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
```

Extract audio content using FFmpeg via the *RipPeriscopeContent* function. In this example, the function will extract 10s of the *collectedHLS[mm]* link and save it to a new file using *newWavFile*.
``` go
// rip content for each HLS that was collected
for mm := 0; mm < len(collectedHLS); mm++ {
	newWavFile := "file" + strconv.Itoa(mm) + ".wav"
	fmt.Println(collectedHLS[mm])
	fmt.Println()
	RipPeriscopeContent(10, collectedHLS[mm], newWavFile)
}
```

FFmpeg function for extracting Periscope content.
``` go
// RipPeriscopeContent : Maintains the ffmpeg script for extracting .wav files from the Periscope stream
func RipPeriscopeContent(durationOfRip int, periscopeLink string, outputFile string) {
	sh.Command("ffmpeg", "-ss", "0", "-t", "10", "-i", periscopeLink, "-strict", "-2", "-ac", "1", outputFile, "-nostdin", "-nostats").Run()
}
