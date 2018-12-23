FROM golang
ADD . /go/src/ArbitrageFinder
RUN go install ArbitrageFinder
ENTRYPOINT /go/bin/ArbitrageFinder
EXPOSE 8080


