FROM golang:1.11
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app
EXPOSE 10000 
RUN go build
CMD ["./app"]
