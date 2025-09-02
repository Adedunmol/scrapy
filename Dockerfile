# Build the application from source
FROM golang:1.24.1-alpine3.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

#RUN #go install github.com/pressly/goose/v3/cmd/goose@latest

RUN CGO_ENABLED=0 GOOS=linux go build -o ./main.exe ./cmd/webserver/main.go

# Development
FROM build-stage AS dev-stage

WORKDIR /app

#RUN go install github.com/air-verse/air@latest

#CMD ["air -c .air.toml"]

ENTRYPOINT ["./main.exe"]