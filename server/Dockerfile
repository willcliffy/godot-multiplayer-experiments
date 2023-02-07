FROM golang:1.19-alpine AS builder

RUN mkdir -p /app
WORKDIR /app
COPY . /app

RUN COOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/dist/app /app/main.go

FROM alpine:3.17

RUN mkdir -p /app
WORKDIR /app

COPY --from=builder /app/dist/app .

EXPOSE 8080
CMD /app/app
