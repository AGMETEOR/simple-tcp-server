# syntax=docker/dockerfile:1

FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /simple-tcp-server

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /simple-tcp-server /simple-tcp-server

EXPOSE 8083/tcp

USER nonroot:nonroot

ENTRYPOINT ["/simple-tcp-server"]

