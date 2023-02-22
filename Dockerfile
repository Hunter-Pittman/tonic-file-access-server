FROM alpine
COPY . /opt/tonic
WORKDIR /opt/tonic
CMD ["./tonic-file-access-server", "--api-token='mysecrettoken'"]
