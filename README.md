# tonic-file-access-server



## Docker Container
This repo contains the code for a quick and easy file access service built for the cloud. It is fast, reliable, and takes minimal resource utilization, but what it does not have is proper security (on purpose) so don't use this for anything that requires security.

```bash
make
sudo docker build . -t "tonic"
sudo docker volume create srv
sudo docker run -d --mount type=volume,src=srv,target=/srv --network=host "tonic"
```
## Endpoints

### /ping
GET Request
Returns:
```
message {
  "pong"
 }
```

### /upload
POST Request
Takes an array of files in multipart/form format and saves it to the preselected server storage.

Returns:
```
STATUS 200 OK
```

### /download/:filename
GET Request
Downloads a file added to the last part of the URL path
Returns:
```
A file
```

### /listdirectory
GET Request
Returns the file directory
```
[]File
```
