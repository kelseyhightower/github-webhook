FROM golang:1.10.0
WORKDIR /go/src/github.com/kelseyhightower/github-webhook
COPY . .
RUN go build .

FROM debian:9.3
RUN apt-get update && apt-get install -y git
COPY --from=0 /go/src/github.com/kelseyhightower/github-webhook .
CMD ["/github-webhook"]
