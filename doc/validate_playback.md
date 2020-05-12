# How to validate playback
In order to validate the packaged content, a MPEG-DASH player is required. 
Chrome, for instance, has an native embedded player that works for unencrypted streams.

In order to play encrypted streams like in this case, you will need a player in which you can set the clearkey.

Providing the `*.mpd` url and the `kid:key` set this player will be able to download the media chunks as required and play them.

## Checking the stream download
A good way to check if chunks are being downloaded correctly is using the network inspector available in every browser. 
There you will see all the related request to the stream files. 

## Checking the files
To check the individual file correctness, you can use `ffmpeg` to decrypt and play the MP4 files (both audio and video) providing the corresponding `key`.
