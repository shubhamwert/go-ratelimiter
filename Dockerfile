
FROM golang:alpine as builder
RUN apk update 
WORKDIR /go/RateLimiter/src/
COPY . .
RUN go get github.com/gorilla/mux    
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/rateLimiter .
# RUN go build -o /go/bin/rateLimiter
RUN apk --no-cache add ca-certificates

FROM busybox
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
WORKDIR /go/bin/
COPY --from=builder /go/bin/rateLimiter .
RUN chmod +x rateLimiter
EXPOSE 8080
ENTRYPOINT ["./rateLimiter"]