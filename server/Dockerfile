FROM golang:1.21.5-alpine3.18 as builder

WORKDIR /usr/local/go/src/server

COPY . .

RUN go mod tidy
RUN go build -o main


FROM alpine:latest
COPY --from=builder /usr/local/go/src/server/static/appv2.html .
COPY --from=builder /usr/local/go/src/server/.env .
COPY --from=builder /usr/local/go/src/server/main .

EXPOSE 4321

CMD ["/main"]

