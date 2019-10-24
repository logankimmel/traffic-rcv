FROM golang:1.13 as golang
COPY main.go /go/src/traffic-rcv/
WORKDIR /go/src/traffic-rcv
RUN cd /go/src/traffic-rcv && go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o traffic-rcv .

FROM scratch
COPY --from=golang /go/src/traffic-rcv/traffic-rcv /
CMD ["/traffic-rcv"]
EXPOSE 8080
