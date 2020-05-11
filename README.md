# go-video-packager
A video packager written in Go

This microservice provide a simple API that allows to:
- Upload MP4 content files
- Package uploaded content for MPEG-DASH
- Expose stream files for playback 

## Run docker-compose
In order to run the app on docker through the use of docker-compose run the following:

<pre><code>$ cd go-video-packager   // go to the project directory
$ docker-compose build   // build all docker images
$ docker-compose up      // run all containers 
</code></pre>

This will make docker run two containers. The app container and mongodb container.

## References
- [Mongodb](https://www.mongodb.com) to store content and stream info.
- [Shaka packager](https://github.com/google/shaka-packager) to package and encrypt content.

## Documentation
- [API](https://github.com/diegofalk/go-video-packager/blob/master/doc/API.md)
- Design
- How to validate playback
- Multi resolution solution plan