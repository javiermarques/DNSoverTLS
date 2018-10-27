FROM golang:alpine
# Install SSL ca certificates
RUN apk update && apk add ca-certificates

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main *.go

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/app/main /
EXPOSE 53/tcp 53/udp
ENTRYPOINT ["/main"]