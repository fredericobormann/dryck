FROM golang:1.14-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/dryck

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build main.go

FROM alpine as final

WORKDIR /app/dryck

COPY --from=build_base /tmp/dryck/main .
COPY --from=build_base /tmp/dryck/static ./static
COPY --from=build_base /tmp/dryck/templates ./templates

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./main"]