<p align="center">
  <img width="35%" height="35%" src="https://github.com/cskonopka/scopingVoices/blob/master/img/scopingVoices-logo.png?raw=true"/>
</p>

<p align="center"><em>*scopingVoices* is a creative tool for extracting audio and video content from the Periscope using Go, FFmpeg and Twitter.</em></p> 

## Why? 
The goal behind the project was to create custom utility for extracting audio and video content from around the world using [Periscope](https://www.pscp.tv). Once collected, the sampling sources are available for esoteric music compositions and video art.


## Twitter Setup
1. Generate [Twitter Access Tokens](https://developer.twitter.com/en/docs/basics/authentication/guides/access-tokens.html)

2. Add the Twitter crendentials to your *~/.bash_profile* 
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

## Go setup
Build the application.
``` go
go build scopingVoices.go
```

Run the application and add an output directory.
```go
./scopingVoices.go location/of/output/directory
```

