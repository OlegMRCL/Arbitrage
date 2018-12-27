FROM golang

RUN go get github.com/OlegMRCL/ArbitrageFinder
WORKDIR /go/src/github.com/OlegMRCL/ArbitrageFinder
RUN go install

ENTRYPOINT /go/bin/ArbitrageFinder
EXPOSE 8181


