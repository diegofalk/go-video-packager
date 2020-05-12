# High level Design
In this document I will cover some high level design decisions such us selected infrastructure, code modules and data models.

## Infrastructure
### Golang
The language suits perfectly for developing microservices. It is easy to use and efficient when it comes to parallelism and process synchronization. 
Also, I’m familiar with it.
### MongoDB
A NoSQL document based database to persist the file's information and the related packaging jobs. 
Mongodb is widely used and easy to manage.
### Shaka packager
This packager fulfills the packaging and encryption requirements. I could have also used bento4 tools for this task.
### Docker
To test and develop the app, we will create a docker image and run the containers using docker compose. 
That makes it portable, easy to run.

## Data models
How the information get organized and persisted
### Content
A `content` is an uploaded file ready to be packaged. Defined by the following structure
```
{
    id
    name
}
```
Where `id` is a unique content id and `name` is the file name.
The uploaded file will be saved locally under `content/` as `name`. 
Unix timestamp is appended to the file name to avoid file overloading in storage.
### Stream
A `stream` represents a packaging job. Defined by the following structure
```
{
    id
    content_id
    key
    kid
    status
    url
}
```
- `id`: Unique stream id 
- `content_id`: content id to package 
- `key` and `kid`: Keys required for encryption
- `status`: Packaging status. Can be `IN PROGRESS`, `DONE` or `FAILED`
- `url`: When `status` is `DONE` it contains the .mpd manifest url.

## Module design

### Main processes
The application is divided in two main parallel processes. `API` and `Packager`. 

#### API
REST API interface. Described [here](https://github.com/diegofalk/go-video-packager/blob/master/doc/API.md)

Responsibilities:
- Handling requests and responses
- Upload and save content files
- Trigger packaging jobs
- Expose stream files

#### Packager
Packager is the module responsible for the actual packaging process. In this case, by calling shaka packager.
It will normally be blocked waiting for packaging jobs. To be performed in packaging workers.

#### Communication
The communication between this two will be through native golang queues (channels). 

### Other modules
#### Storage
Handles local storage. Save and Load content and stream files. 
It can be easily replaced by cloud storage.
#### database
Encapsulates mongodb interaction. Define data models.

## References
- [Golang](https://golang.org/)
- [Mongodb](https://www.mongodb.com)
- [Shaka packager](https://github.com/google/shaka-packager) 
- [Docker](https://www.docker.com/)
