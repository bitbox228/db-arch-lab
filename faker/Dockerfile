FROM golang:1.22-alpine as builder

WORKDIR /build

COPY .. .

RUN go mod download

RUN go build -o app ./cmd/main.go

FROM alpine

WORKDIR /root/

COPY --from=builder /build/app .

CMD [ "./app" ]