# API
This document describes the API endpoints and usage

## Upload content
Use this endpoint to upload MP4 content
#### URL
`/publish/<content>.mp4`
#### Method
`POST`
#### Request Body
MP4 file
#### Response Body
200
```
{
    "content_id": "5eb8d6f6eebd404bf330e00e"
}
```

## Package
This endpoint triggers packaging process
#### URL
`/package`
#### Method
`POST`
#### Request Body
```
{
	"content_id": "5eb8d6f6eebd404bf330e00e",
	"key": "hyN9IKGfWKdAwFaE5pm0qg",
	"kid": "oW5AK5BW43HzbTSKpiu3SQ"
}
```
#### Response Body
200
```
{
    "stream_id": "5eb8d705eebd404bf330e00f"
}
```

## Stream info
Get packed stream info
#### URL
`/streaminfo/<stream_id>`
#### Method
`GET`
#### Request Body
None
#### Response Body
200
```
{
    "url":"http://localhost:8081/stream/5eb8d705eebd404bf330e00f/5eb8d705eebd404bf330e00f.mpd",
    "key":"hyN9IKGfWKdAwFaE5pm0qg",
    "kid":"oW5AK5BW43HzbTSKpiu3SQ"
}
```
202
```
Packaging in progress
```

## Stream files
Get stream files
#### URLs
`/stream/<stream_id>/<stream_id>.mpd`
`/stream/<stream_id>/<folder>/<init_file>.mp4`
`/stream/<stream_id>/<folder>/<chunk_file>.m4s`
#### Method
`GET`
#### Request Body
None
#### Response Body
Stream file