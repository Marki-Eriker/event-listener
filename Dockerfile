FROM golang:1.19 AS builder

WORKDIR $GOPATH/test
COPY . ./
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /test/service

FROM alpine

RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
COPY --from=builder /test ./
COPY service-docker.toml /service-dev.toml

ENTRYPOINT ./service
