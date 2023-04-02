FROM golang:1.19-alpine AS server-builder

RUN mkdir -p /app
WORKDIR /app
COPY server/ .

RUN COOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/dist/app /app/main.go

FROM alpine:3.17

RUN mkdir -p /app
WORKDIR /app

COPY --from=server-builder /app/dist/app /app/app

EXPOSE 8080
CMD /app/app
