FROM golang:1.13-stretch as build-env
# For golang >= 1.11

RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum  .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/server

# Runtime images, also FROM scratch
FROM alpine:3.10
WORKDIR /app
COPY --from=build-env /go/bin/server .
COPY config.yaml config.yaml

ENTRYPOINT ["/app/server"]
