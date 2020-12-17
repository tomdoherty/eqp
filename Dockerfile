FROM golang:1.15.6 AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/run cmd/eqp.go
FROM scratch AS bin
COPY --from=build /out/run /

CMD ["/run"]
