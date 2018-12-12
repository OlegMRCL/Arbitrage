FROM golang
ADD . /go/src/ArbitageFinder
RUN go install ArbitageFinder
ENTRYPOINT /go/bin/ArbitageFinder
EXPOSE 8080


