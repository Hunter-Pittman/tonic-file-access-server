# tonic-file-access-server
This repo contains the code for a quick and easy file access service built for the cloud. It is fast, reliable, and takes minimal resource utilization.

```bash
make
sudo docker build . -t "tonic"
sudo docker volume create srv
sudo docker run -d --mount type=volume,src=srv,target=/srv --network=host "tonic"
```
