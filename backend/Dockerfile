FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/server /usr/local/bin/server
EXPOSE 8081
ENTRYPOINT ["server"]
