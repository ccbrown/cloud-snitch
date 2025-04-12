FROM golang:1.24.1-bookworm AS builder

WORKDIR /go/src/github.com/ccbrown/cloud-snitch

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go generate ./...
RUN go build .

FROM debian:bookworm-slim AS runner

RUN apt update && apt install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/backend/bin

COPY --from=0 /go/src/github.com/ccbrown/cloud-snitch/backend .

RUN ./backend -h

ENTRYPOINT ["/opt/backend/bin/backend"]
