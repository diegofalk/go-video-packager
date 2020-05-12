# go-video-packager
A video packager written in Go

This microservice provide a simple API that allows to:
- Upload MP4 content files
- Package uploaded content for MPEG-DASH
- Expose stream files for playback 

## Run
In order to run the app on docker through the use of docker-compose run the following:

<pre><code>$ cd go-video-packager   // go to the project directory
$ docker-compose build   // build all docker images
$ docker-compose up      // run all containers 
</code></pre>

This will make docker run two containers. The app container and mongodb container.

## Basic workflow
#### 1. Upload content
Upload a file using `/content` endpoint. The response will include a `content_id`.
#### 2. Package
Calling `/package` endpoint with the `content_id` and the encryption params will trigger a packaging job. 
The actual packaging will occur in a background process. 
The response will include a `stream_id`.
#### 3. Get stream
Calling `/streaminfo` with the required `stream_id` will check for the packaging job status. 
If the packaging job is done, the response will provide a stream `url` and the needed params to play the stream. 

## References
- [Mongodb](https://www.mongodb.com) to store content and stream info.
- [Shaka packager](https://github.com/google/shaka-packager) to package and encrypt content.

## Documentation
- [API](https://github.com/diegofalk/go-video-packager/blob/master/doc/API.md)
- [Design](https://github.com/diegofalk/go-video-packager/blob/master/doc/high_level_design.md)
- [How to validate playback](https://github.com/diegofalk/go-video-packager/blob/master/doc/validate_playback.md)
- [Multi bitrate plan](https://github.com/diegofalk/go-video-packager/blob/master/doc/multi_bitrate.md)