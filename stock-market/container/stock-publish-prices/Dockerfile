# STEP 1 
FROM golang:alpine as builder
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR $GOPATH/bin/
COPY ./exe/stock-publish-prices .

# STEP 2 build a small image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/stock-publish-prices /go/bin/stock-publish-prices

ENTRYPOINT ["/go/bin/stock-publish-prices"]