# test-rip.sh
#!/bin/bash

ffmpeg -ss 0 -t 10 -i https://prod-fastly-ap-southeast-1.video.periscope.tv/Transcoding/v1/lhls/aksUJRBAr4ZqFN2o0hwNU-zS_GSeSmQw-Rr1SGN-hVvuhJ56ms1BgGR_IOZWyO9ZgrgMMCs7G9vnMD88Pix5-g/non_transcode/ap-southeast-1/periscope-replay-direct-prod-ap-southeast-1-public/dynamic_lowlatency.m3u8 -strict -2 -ac 1 /Users/io/Desktop/tester.wav -nostdin -nostats