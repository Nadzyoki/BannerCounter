# syntax=docker/dockerfile:1

FROM golang:1.25.0-alpine3.22 AS build

ENV GO111MODULE=on
ENV PATH=$PATH:$GOPATH/bin
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go


FROM alpine:3.22

WORKDIR /app

COPY --from=build /app/main .
COPY migrations/ ./migrations/

CMD ["/app/main"]