# scopingVoices
*scopingVoices* is an artist tool used for extracting audio and video content from the Periscope API in real-time using Go and FFmpeg.

### Requirements
```
go run github.com/ChimeraCoder/anaconda
go run github.com/codeskyblue/go-sh
```

## How to run 
```
go run SV-v1.go
```

### How it works

Main structs
```
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
```
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
```
	// "periscope" results
	search_result, err := api.GetSearch("periscope", nil)
	if err != nil {
		panic(err)
	}
```

Extract the "expanded_url" of each result and find the substring "https://www.pscp.tv/w/". Once the substring is found, append to a string array named "selectedfiles".
```
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
```
	// Split strings and add the broadcast_id to the new url
	var newURLS []string
	for i := 0; i < len(selectedfiles); i++ {
		mo := strings.SplitAfter(selectedfiles[i], "https://www.pscp.tv/w/")
		newURLS = append(newURLS, "https://api.periscope.tv/api/v2/accessVideoPublic?broadcast_id="+mo[1])
	}
```


Create a GET Request for each element within the "newURLS" string array. Use the *Periscope* struct to decode the JSON to receive the *lhls_url* of live and replay streams. Add "active" streams to a new string array named "collectedHLS".
```
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
```
	// rip content for each HLS that was collected
	for mm := 0; mm < len(collectedHLS); mm++ {
		newWavFile := "file" + strconv.Itoa(mm) + ".wav"
		fmt.Println(collectedHLS[mm])
		fmt.Println()
		RipPeriscopeContent(10, collectedHLS[mm], newWavFile)
	}
```

FFmpeg function for extracting Periscope content.
``` 
// RipPeriscopeContent : Maintains the ffmpeg script for extracting .wav files from the Periscope stream
func RipPeriscopeContent(durationOfRip int, periscopeLink string, outputFile string) {
	sh.Command("ffmpeg", "-ss", "0", "-t", "10", "-i", periscopeLink, "-strict", "-2", "-ac", "1", outputFile, "-nostdin", "-nostats").Run()
}
